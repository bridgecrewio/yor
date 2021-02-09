package structure

import (
	"bridgecrewio/yor/common"
	"bridgecrewio/yor/common/structure"
	"bridgecrewio/yor/common/tagging/tags"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/hashicorp/terraform/command"
	"github.com/mitchellh/cli"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
)

var prefixToTagAttribute = map[string]string{"aws": "tags", "azure": "tags", "gcp": "labels"}

type TerrraformParser struct {
	generatedPath string
}

func NewTerrraformParser() *TerrraformParser {
	terraformParser := new(TerrraformParser)
	_, currFile, _, _ := runtime.Caller(0)
	currDir := path.Join(path.Dir(currFile))
	terraformParser.generatedPath = path.Join(currDir, "./.generated")

	return terraformParser
}

func (p *TerrraformParser) getGeneratedPathForDir(dir string) string {
	dirName := path.Base(dir)
	return path.Join(p.generatedPath, dirName)
}

func (p *TerrraformParser) TerraformInitDirectory(directory string) error {
	generatedPath := p.getGeneratedPathForDir(directory)
	if _, err := os.Stat(generatedPath); !os.IsNotExist(err) {
		fmt.Printf("directory already initialized\n")
		return nil
	}
	initCommand := &command.InitCommand{
		Meta: command.Meta{
			Ui:              &cli.MockUi{},
			OverrideDataDir: generatedPath,
		},
	}
	args := []string{directory}
	code := initCommand.Run(args)
	if code != 0 {
		return fmt.Errorf("failed to run terraform init on directory %s", directory)
	}

	return nil

}

func (p *TerrraformParser) GetSourceFiles(directory string) ([]string, error) {
	errMsg := "failed to get .tf files because %s"
	err := p.TerraformInitDirectory(directory)
	if err != nil {
		return nil, fmt.Errorf(errMsg, err)
	}

	modulesDirectories, err := p.getModulesDirectories(directory)
	if err != nil {
		return nil, err
	}

	var files []string
	for _, dir := range modulesDirectories {

		err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() && strings.HasSuffix(info.Name(), ".tf") {
				files = append(files, path)
			}
			return nil
		})
		if err != nil {
			return nil, fmt.Errorf(errMsg, err)
		}
	}

	return files, nil
}

func (p *TerrraformParser) getModulesDirectories(directory string) ([]string, error) {
	errMsg := "failed to get all modules directories because %s"
	modulesJsonFile, err := os.Open(p.getGeneratedPathForDir(directory) + "/modules/modules.json")
	var modulesFile ModulesFile
	if err != nil {
		return nil, fmt.Errorf(errMsg, err)
	}

	moduleFileData, _ := ioutil.ReadAll(modulesJsonFile)
	err = json.Unmarshal(moduleFileData, &modulesFile)
	if err != nil {
		return nil, fmt.Errorf(errMsg, err)
	}

	modulesDirectories := make([]string, 0)
	for _, entry := range modulesFile.Modules {
		moduleDir := path.Join(directory, entry.Source)
		if _, err := os.Stat(moduleDir); !os.IsNotExist(err) && !common.InSlice(modulesDirectories, moduleDir) {
			// if directory exists (local module) and modulesDirectories doesn't contain it yet, add it
			modulesDirectories = append(modulesDirectories, moduleDir)
		}
	}

	return modulesDirectories, nil
}

func (p *TerrraformParser) ParseFile(filePath string) ([]structure.IBlock, error) {
	// read file bytes
	src, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s because %s", filePath, err)
	}

	// parse the file into hclwrite.File and hclsyntax.File to allow getting existing tags and lines
	hclFile, diagnostics := hclwrite.ParseConfig(src, filePath, hcl.InitialPos)
	if diagnostics != nil && diagnostics.HasErrors() {
		hclErrors := diagnostics.Errs()
		return nil, fmt.Errorf("failed to parse hcl file %s because of errors %s", filePath, hclErrors)
	}
	hclSyntaxFile, diagnostics := hclsyntax.ParseConfig(src, filePath, hcl.InitialPos)
	if diagnostics != nil && diagnostics.HasErrors() {
		hclErrors := diagnostics.Errs()
		return nil, fmt.Errorf("failed to parse hcl file %s because of errors %s", filePath, hclErrors)
	}

	if hclFile == nil || hclSyntaxFile == nil {
		return nil, fmt.Errorf("failed to parse hcl file %s", filePath)
	}

	syntaxBlocks := hclSyntaxFile.Body.(*hclsyntax.Body).Blocks
	rawBlocks := hclFile.Body().Blocks()
	parsedBlocks := make([]structure.IBlock, 0)
	for i, block := range rawBlocks {
		terraformBlock, err := p.parseBlock(block)
		if err != nil {
			return nil, fmt.Errorf("failed to parse terraform block because %s", err)
		}

		terraformBlock.Init(filePath, block)
		terraformBlock.AddHclSyntaxBlock(syntaxBlocks[i])
		parsedBlocks = append(parsedBlocks, terraformBlock)
	}

	return parsedBlocks, nil
}

func (p *TerrraformParser) WriteFile(filePath string, blocks []structure.IBlock) error {
	for _, block := range blocks {
		tfBlock := block.(*TerraformBlock)
		tfBlock.MergeTags()
	}
	return nil
}

func (p *TerrraformParser) parseBlock(hclBlock *hclwrite.Block) (*TerraformBlock, error) {
	var existingTags []tags.ITag
	var tagsAttributeName string
	var err error
	isTaggable := false

	if hclBlock.Type() == "resource" {
		tagsAttributeName, err = p.getTagsAttributeName(hclBlock)
		if err != nil {
			return nil, err
		}
		existingTags, isTaggable = p.getExistingTags(hclBlock, tagsAttributeName)

		if !isTaggable {
			isTaggable, err = p.isBlockTaggable(hclBlock)
			if err != nil {
				return nil, err
			}
		}
	}

	terraformBlock := TerraformBlock{
		Block: structure.Block{
			ExitingTags:       existingTags,
			IsTaggable:        isTaggable,
			TagsAttributeName: tagsAttributeName,
		},
	}

	return &terraformBlock, nil
}

func (p *TerrraformParser) getTagsAttributeName(hclBlock *hclwrite.Block) (string, error) {
	resourceType := hclBlock.Labels()[0]
	tagsAttributeName, err := getTagAttributeByResourceType(resourceType)
	if err != nil {
		return "", err
	}

	return tagsAttributeName, nil
}

func getProviderFromResourceType(resourceType string) string {
	provider := strings.Split(resourceType, "_")[0]
	return provider
}

func getTagAttributeByResourceType(resourceType string) (string, error) {
	prefix := prefixToTagAttribute[getProviderFromResourceType(resourceType)]
	if prefix == "" {
		return "", fmt.Errorf("failed to find tags attribute name for resource type %s", resourceType)
	}

	return prefix, nil
}

func (p *TerrraformParser) getExistingTags(hclBlock *hclwrite.Block, tagsAttributeName string) ([]tags.ITag, bool) {
	isTaggable := false
	existingTags := make([]tags.ITag, 0)

	tagsAttribute := hclBlock.Body().GetAttribute(tagsAttributeName)
	if tagsAttribute != nil {
		// if tags exists in resource
		isTaggable = true
		tagsTokens := tagsAttribute.Expr().BuildTokens(hclwrite.Tokens{})
		parsedTags := p.parseTagLines(tagsTokens)
		for key := range parsedTags {
			iTag := tags.Init(key, parsedTags[key])
			existingTags = append(existingTags, iTag)
		}
	}

	return existingTags, isTaggable
}

func (p *TerrraformParser) isBlockTaggable(hclBlock *hclwrite.Block) (bool, error) {
	// TODO - implement like IsTaggable in https://github.com/env0/terratag/blob/master/tfschema/tfschema.go

	resourceType := hclBlock.Labels()[0]
	return common.InSlice(TaggableResourceTypes, resourceType), nil
}

func (p *TerrraformParser) parseTagLines(tokens hclwrite.Tokens) map[string]string {
	parsedTags := make(map[string]string)
	entries := make([]hclwrite.Tokens, 0)
	startIndex := 0
	hasEq := false
	for i, token := range tokens {
		if token.Type == hclsyntax.TokenNewline {
			if hasEq {
				entries = append(entries, tokens[startIndex:i])
			}
			startIndex = i + 1
			hasEq = false
		}
		if token.Type == hclsyntax.TokenEqual {
			hasEq = true
		}
	}
	if hasEq {
		entries = append(entries, tokens[startIndex:])
	}

	for _, entry := range entries {
		eqIndex := -1
		var key string
		for j, token := range entry {
			if token.Type == hclsyntax.TokenEqual {
				eqIndex = j + 1
				key = strings.TrimSpace(string(entry[:j].Bytes()))
			}
		}
		value := string(entry[eqIndex:].Bytes())
		value = strings.TrimPrefix(strings.TrimSuffix(value, " "), " ")
		parsedTags[key] = value
	}

	return parsedTags
}

type ModulesFile struct {
	Modules []ModuleEntry `json:"Modules"`
}

type ModuleEntry struct {
	Key    string `json:"Key"`
	Source string `json:"Source"`
	Dir    string `json:"Dir"`
}
