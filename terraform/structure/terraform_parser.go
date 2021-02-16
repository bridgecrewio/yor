package structure

import (
	"bridgecrewio/yor/common"
	"bridgecrewio/yor/common/logger"
	"bridgecrewio/yor/common/structure"
	"bridgecrewio/yor/common/tagging/tags"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/hashicorp/terraform/command"
	"github.com/minamijoyo/tfschema/tfschema"
	"github.com/mitchellh/cli"
	"github.com/zclconf/go-cty/cty"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
)

const TerraformOutputDir = "/.terraform"

var prefixToTagAttribute = map[string]string{"aws": "tags", "azure": "tags", "gcp": "labels"}

type TerrraformParser struct {
	rootDir                string
	providerToClientMap    map[string]tfschema.Client
	taggableResourcesCache map[string]bool
	tagModules             bool
}

func (p *TerrraformParser) Init(rootDir string, args map[string]string) {
	p.rootDir = rootDir
	p.providerToClientMap = make(map[string]tfschema.Client)
	p.taggableResourcesCache = make(map[string]bool)
	p.tagModules = true
	if argTagModule, ok := args["tag-modules"]; ok {
		p.tagModules, _ = strconv.ParseBool(argTagModule)
	}
}

func (p *TerrraformParser) TerraformInitDirectory(directory string) error {
	terraformOutputPath := directory + TerraformOutputDir
	if _, err := os.Stat(terraformOutputPath); !os.IsNotExist(err) {
		logger.Info("directory already initialized\n")
		return nil
	}
	initCommand := &command.InitCommand{
		Meta: command.Meta{
			Ui:              &cli.MockUi{},
			OverrideDataDir: terraformOutputPath,
		},
	}
	fmt.Printf("Could not locate %s directury under %s, running terraform init\n", TerraformOutputDir, directory)
	args := []string{directory}
	code := initCommand.Run(args)
	if code != 0 {
		return fmt.Errorf("failed to run terraform init on directory %s, please run it manually", directory)
	}
	if _, err := os.Stat(terraformOutputPath); !os.IsNotExist(err) {
		logger.Info("directory initialized successfully")
		return nil
	}

	return fmt.Errorf("failed to initialize directory %s, the folder '%s' was not created", directory, TerraformOutputDir)
}

func (p *TerrraformParser) GetSourceFiles(directory string) ([]string, error) {
	errMsg := "failed to get .tf files because %s"
	var modulesDirectories []string

	err := p.TerraformInitDirectory(directory)
	if err != nil {
		return nil, fmt.Errorf(errMsg, err)
	}

	if p.tagModules {
		modulesDirectories, err = p.getModulesDirectories(directory)
		if err != nil {
			return nil, err
		}
	} else {
		modulesDirectories = []string{directory}
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
	modulesJSONFile, err := os.Open(directory + TerraformOutputDir + "/modules/modules.json")
	var modulesFile ModulesFile
	if err != nil {
		return nil, fmt.Errorf(errMsg, err)
	}

	moduleFileData, _ := ioutil.ReadAll(modulesJSONFile)
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
	// read file bytes
	src, err := ioutil.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file %s because %s", filePath, err)
	}

	// parse the file into hclwrite.File and hclsyntax.File to allow getting existing tags and lines
	hclFile, diagnostics := hclwrite.ParseConfig(src, filePath, hcl.InitialPos)
	if diagnostics != nil && diagnostics.HasErrors() {
		hclErrors := diagnostics.Errs()
		return fmt.Errorf("failed to parse hcl file %s because of errors %s", filePath, hclErrors)
	}

	if hclFile == nil {
		return fmt.Errorf("failed to parse hcl file %s", filePath)
	}

	rawBlocks := hclFile.Body().Blocks()
	for _, rawBlock := range rawBlocks {
		rawBlockLabels := rawBlock.Labels()
		for _, parsedBlock := range blocks {
			if parsedBlock.IsBlockTaggable() {
				parsedBlockLabels := parsedBlock.(*TerraformBlock).HclSyntaxBlock.Labels
				if reflect.DeepEqual(parsedBlockLabels, rawBlockLabels) {
					p.modifyBlockTags(rawBlock, parsedBlock)
				}
				print(parsedBlock)
			}
		}
	}
	f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_TRUNC, os.ModeAppend)
	if err != nil {
		return err
	}
	_, err = hclFile.WriteTo(f)
	if err != nil {
		return fmt.Errorf("failed to write HCL file %s, %s", filePath, err.Error())
	}
	if err = f.Close(); err != nil {
		return err
	}
	return nil
}

func (p *TerrraformParser) modifyBlockTags(rawBlock *hclwrite.Block, parsedBlock structure.IBlock) {
	tagsAttributeName := parsedBlock.(*TerraformBlock).TagsAttributeName
	tagsAttribute := rawBlock.Body().GetAttribute(tagsAttributeName)
	mergedTags := parsedBlock.MergeTags()
	rawTagsTokens := tagsAttribute.Expr().BuildTokens(hclwrite.Tokens{})
	isMergeOpExists := false

	for _, rawTagsToken := range rawTagsTokens {
		fmt.Println(string(rawTagsToken.Type), string(rawTagsToken.Bytes))
		if string(rawTagsToken.Bytes) == "merge" {
			isMergeOpExists = true
			break
		}
	}
	var tagTypes = tags.TagTypes
	if !isMergeOpExists {
		//Insert the merge token, opening and closing parenthesis tokens
		rawTagsTokens = InsertToken(rawTagsTokens, 0, &hclwrite.Token{
			Type:  hclsyntax.TokenIdent,
			Bytes: []byte("merge"),
		})
		rawTagsTokens = InsertToken(rawTagsTokens, 1, &hclwrite.Token{
			Type:  hclsyntax.TokenOParen,
			Bytes: []byte("("),
		})
		rawTagsTokens = InsertToken(rawTagsTokens, len(rawTagsTokens), &hclwrite.Token{
			Type:  hclsyntax.TokenCParen,
			Bytes: []byte(")"),
		})
	}
	for _, tagType := range tagTypes {
		fmt.Println(tagType)
	}
	//Insert a comma token before the merge closing parenthesis
	rawTagsTokens = InsertToken(rawTagsTokens, len(rawTagsTokens)-1, &hclwrite.Token{
		Type:  hclsyntax.TokenComma,
		Bytes: []byte(","),
	})
	mergedTagsTokens := buildTagsTokens(mergedTags)
	for _, tagToken := range mergedTagsTokens {
		rawTagsTokens = InsertToken(rawTagsTokens, len(rawTagsTokens)-1, tagToken)
	}
	//Set the body's tags to the new built tokens
	rawBlock.Body().SetAttributeRaw(tagsAttributeName, rawTagsTokens)
	print(rawTagsTokens, mergedTags, isMergeOpExists, mergedTagsTokens)
}

func buildTagsTokens(tags []tags.ITag) hclwrite.Tokens {
	tagsMap := make(map[string]cty.Value, len(tags))
	for _, tag := range tags {
		tagsMap[tag.GetKey()] = cty.StringVal(tag.GetValue())
	}
	return hclwrite.TokensForValue(cty.MapVal(tagsMap))
}

// Insert inserts a value at a specific index in a slice
func InsertToken(tokens hclwrite.Tokens, index int, value *hclwrite.Token) hclwrite.Tokens {
	if len(tokens) == index {
		return append(tokens, value)
	}
	tokens = append(tokens[:index+1], tokens[index:]...)
	tokens[index] = value
	return tokens
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
		parsedTags := p.parseTagAttribute(tagsTokens)
		for key := range parsedTags {
			iTag := tags.Init(key, parsedTags[key])
			existingTags = append(existingTags, iTag)
		}
	}

	return existingTags, isTaggable
}

func (p *TerrraformParser) isBlockTaggable(hclBlock *hclwrite.Block) (bool, error) {
	resourceType := hclBlock.Labels()[0]
	if val, ok := p.taggableResourcesCache[resourceType]; ok {
		return val, nil
	}
	tagAtt, err := getTagAttributeByResourceType(resourceType)
	if err != nil {
		return false, err
	}

	providerName := getProviderFromResourceType(resourceType)

	client := p.getClient(providerName)
	taggable := false
	if client != nil {
		typeSchema, err := client.GetResourceTypeSchema(resourceType)
		if err != nil {
			if strings.Contains(err.Error(), "Failed to find resource type") {
				// Resource Type doesn't have schema yet in the provider
				return false, nil
			}
			return false, err
		}

		if _, ok := typeSchema.Attributes[tagAtt]; ok {
			taggable = true
		}
	}
	p.taggableResourcesCache[resourceType] = taggable
	return taggable, nil
}

func (p *TerrraformParser) getHclMapsContents(tokens hclwrite.Tokens) []hclwrite.Tokens {
	// The function gets tokens and returns an array of tokens that are found between curly brackets '{...}'
	// example: tokens: "merge({a=1, b=2}, {c=3})", return: ["a=1, b=2", "c=3"]
	hclMaps := make([]hclwrite.Tokens, 0)
	bracketOpenIndex := -1

	for i, token := range tokens {
		if token.Type == hclsyntax.TokenOBrace {
			bracketOpenIndex = i
		}
		if token.Type == hclsyntax.TokenCBrace {
			hclMaps = append(hclMaps, tokens[bracketOpenIndex+1:i])
		}
	}

	return hclMaps
}

func (p *TerrraformParser) extractTagPairs(tokens hclwrite.Tokens) []hclwrite.Tokens {
	// The function gets tokens and returns an array of tokens that represent key and value
	// example: tokens: "a=1\n b=2, c=3", returns: ["a=1", "b=2", "c=3"]
	separatorTokens := []hclsyntax.TokenType{hclsyntax.TokenComma, hclsyntax.TokenNewline}
	tagPairs := make([]hclwrite.Tokens, 0)
	startIndex := 0
	hasEq := false
	for i, token := range tokens {
		if common.InSlice(separatorTokens, token.Type) {
			if hasEq {
				tagPairs = append(tagPairs, tokens[startIndex:i])
			}
			startIndex = i + 1
			hasEq = false
		}
		if token.Type == hclsyntax.TokenEqual {
			hasEq = true
		}
	}
	if hasEq {
		tagPairs = append(tagPairs, tokens[startIndex:])
	}

	return tagPairs
}

func (p *TerrraformParser) parseTagAttribute(tokens hclwrite.Tokens) map[string]string {
	hclMaps := p.getHclMapsContents(tokens)
	tagPairs := make([]hclwrite.Tokens, 0)
	for _, hclMap := range hclMaps {
		tagPairs = append(tagPairs, p.extractTagPairs(hclMap)...)
	}

	// for each tag pair, find the key and value
	parsedTags := make(map[string]string)
	for _, entry := range tagPairs {
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

func (p *TerrraformParser) getClient(providerName string) tfschema.Client {
	hclLogger := hclog.New(&hclog.LoggerOptions{
		Name:   "plugin",
		Level:  hclog.Trace,
		Output: hclog.DefaultOutput,
	})
	client, exists := p.providerToClientMap[providerName]
	if exists {
		return client
	}
	newClient, err := tfschema.NewClient(providerName, tfschema.Option{
		RootDir: p.rootDir,
		Logger:  hclLogger,
	})

	if err != nil {
		if strings.Contains(err.Error(), "Failed to find plugin") {
			logger.Warning(fmt.Sprintf("Could not load provider %v, resources from this provider will not be tagged", providerName))
		}
		return nil
	}

	p.providerToClientMap[providerName] = newClient
	return newClient
}

type ModulesFile struct {
	Modules []ModuleEntry `json:"Modules"`
}

type ModuleEntry struct {
	Key    string `json:"Key"`
	Source string `json:"Source"`
	Dir    string `json:"Dir"`
}
