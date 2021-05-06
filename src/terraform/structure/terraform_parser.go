package structure

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/bridgecrewio/yor/src/common"
	"github.com/bridgecrewio/yor/src/common/logger"
	"github.com/bridgecrewio/yor/src/common/structure"
	"github.com/bridgecrewio/yor/src/common/tagging/tags"
	"github.com/bridgecrewio/yor/src/common/utils"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/minamijoyo/tfschema/tfschema"
	"github.com/zclconf/go-cty/cty"
)

var ProviderToTagAttribute = map[string]string{"aws": "tags", "azurerm": "tags", "google": "labels"}
var ignoredDirs = []string{".git", ".DS_Store", ".idea", ".terraform"}

type TerrraformParser struct {
	rootDir                string
	providerToClientMap    map[string]tfschema.Client
	taggableResourcesCache map[string]bool
	tagModules             bool
	terraformModule        *TerraformModule
}

func (p *TerrraformParser) Init(rootDir string, args map[string]string) {
	p.rootDir = rootDir
	p.providerToClientMap = make(map[string]tfschema.Client)
	p.taggableResourcesCache = make(map[string]bool)
	p.tagModules = true
	p.terraformModule = NewTerraformModule(rootDir)
	if argTagModule, ok := args["tag-modules"]; ok {
		p.tagModules, _ = strconv.ParseBool(argTagModule)
	}
}

func (p *TerrraformParser) GetSkippedDirs() []string {
	return ignoredDirs
}

func (p *TerrraformParser) GetSupportedFileExtensions() []string {
	return []string{common.TfFileType.Extension}
}

func (p *TerrraformParser) GetSourceFiles(directory string) ([]string, error) {
	errMsg := "failed to get .tf files because %s"
	var modulesDirectories []string

	if p.tagModules {
		modulesDirectories = p.terraformModule.GetModulesDirectories()
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

func (p *TerrraformParser) ParseFile(filePath string) ([]structure.IBlock, error) {
	// read file bytes
	// #nosec G304
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
		if block.Type() != "resource" {
			continue
		}
		blockID := strings.Join(block.Labels(), ".")
		terraformBlock, err := p.parseBlock(block)
		if err != nil {
			if strings.HasPrefix(err.Error(), "resource belongs to skipped provider") {
				logger.Info(fmt.Sprintf("skipping block %s because the provider %s does not support tags", blockID, strings.Split(blockID, "_")[0]))
			} else {
				logger.Warning(fmt.Sprintf("failed to parse terraform block because %s", err.Error()))
			}
			continue
		}
		if terraformBlock == nil {
			logger.Warning(fmt.Sprintf("Found a malformed block according to block scheme %v", blockID))
			continue
		}
		terraformBlock.Init(filePath, block)
		terraformBlock.AddHclSyntaxBlock(syntaxBlocks[i])
		parsedBlocks = append(parsedBlocks, terraformBlock)
	}

	return parsedBlocks, nil
}

func (p *TerrraformParser) WriteFile(readFilePath string, blocks []structure.IBlock, writeFilePath string) error {
	// read file bytes
	// #nosec G304
	src, err := ioutil.ReadFile(readFilePath)
	if err != nil {
		return fmt.Errorf("failed to read file %s because %s", readFilePath, err)
	}

	// parse the file into hclwrite.File and hclsyntax.File to allow getting existing tags and lines
	hclFile, diagnostics := hclwrite.ParseConfig(src, readFilePath, hcl.InitialPos)
	if diagnostics != nil && diagnostics.HasErrors() {
		hclErrors := diagnostics.Errs()
		return fmt.Errorf("failed to parse hcl file %s because of errors %s", readFilePath, hclErrors)
	}

	if hclFile == nil {
		return fmt.Errorf("failed to parse hcl file %s", readFilePath)
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
			}
		}
	}
	// #nosec G304
	f, err := os.OpenFile(writeFilePath, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	_, err = hclFile.WriteTo(f)
	if err != nil {
		return fmt.Errorf("failed to write HCL file %s, %s", readFilePath, err.Error())
	}
	if err = f.Close(); err != nil {
		return err
	}
	return nil
}

func (p *TerrraformParser) modifyBlockTags(rawBlock *hclwrite.Block, parsedBlock structure.IBlock) {
	mergedTags := parsedBlock.MergeTags()
	tagsAttributeName := parsedBlock.(*TerraformBlock).TagsAttributeName
	tagsAttribute := rawBlock.Body().GetAttribute(tagsAttributeName)
	if tagsAttribute == nil {
		mergedTagsTokens := buildTagsTokens(mergedTags)
		if mergedTagsTokens != nil {
			rawBlock.Body().SetAttributeRaw(tagsAttributeName, mergedTagsTokens)
		}
	} else {
		rawTagsTokens := tagsAttribute.Expr().BuildTokens(hclwrite.Tokens{})
		isMergeOpExists := false
		isRenderedAttribute := false
		existingParsedTags := p.parseTagAttribute(rawTagsTokens)
		for _, rawTagsToken := range rawTagsTokens {
			tokenStr := string(rawTagsToken.Bytes)
			if tokenStr == "merge" {
				isMergeOpExists = true
				break
			}
			if utils.InSlice([]string{"var", "local"}, tokenStr) {
				isRenderedAttribute = true
				break
			}
		}

		var replacedTags []tags.ITag
		k := 0
		for _, tag := range mergedTags {
			tagReplaced := false
			strippedTagKey := strings.ReplaceAll(tag.GetKey(), `"`, "")
			for _, rawTagsToken := range rawTagsTokens {
				if string(rawTagsToken.Bytes) == tag.GetKey() || string(rawTagsToken.Bytes) == strippedTagKey {
					replacedTags = append(replacedTags, tag)
					tagReplaced = true
					break
				}
			}
			if !tagReplaced {
				// Keep only new tags (non-appearing) in mergedTags
				mergedTags[k] = tag
				k++
			}
		}

		for _, replacedTag := range replacedTags {
			tagKey := replacedTag.GetKey()
			var existingTagValue string
			if existingTagValue = strings.ReplaceAll(existingParsedTags[tagKey], `"`, ""); existingTagValue == "" {
				quotedTagKey := fmt.Sprintf(`"%s"`, tagKey)
				existingTagValue = strings.ReplaceAll(existingParsedTags[quotedTagKey], `"`, "")
			}
			replacedValue := strings.ReplaceAll(replacedTag.GetValue(), `"`, "")
			foundKey := false
			for _, rawToken := range rawTagsTokens {
				if string(rawToken.Bytes) == tagKey {
					foundKey = true
				}
				if string(rawToken.Bytes) == existingTagValue && foundKey {
					rawToken.Bytes = []byte(replacedValue)
				}
			}
		}

		if !isMergeOpExists && !isRenderedAttribute {
			newTagsTokens := buildTagsTokens(mergedTags)
			rawTagsTokens = InsertTokens(rawTagsTokens, newTagsTokens[2:len(newTagsTokens)-2])
			rawBlock.Body().SetAttributeRaw(tagsAttributeName, rawTagsTokens)
			return
		}

		// These lines execute if there is either a `merge` operator at the start of the tags,
		// or if it is rendered via a variable / local.
		mergedTags = mergedTags[:k]
		mergedTagsTokens := buildTagsTokens(mergedTags)
		if !isMergeOpExists && mergedTagsTokens != nil {
			// Insert the merge token, opening and closing parenthesis tokens
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
		// Insert a comma token before the merge closing parenthesis and add as a separate dict
		if mergedTagsTokens != nil {
			rawTagsTokens = InsertToken(rawTagsTokens, len(rawTagsTokens)-1, &hclwrite.Token{
				Type:  hclsyntax.TokenComma,
				Bytes: []byte(","),
			})
			for _, tagToken := range mergedTagsTokens {
				rawTagsTokens = InsertToken(rawTagsTokens, len(rawTagsTokens)-1, tagToken)
			}
		}
		// Set the body's tags to the new built tokens
		rawBlock.Body().SetAttributeRaw(tagsAttributeName, rawTagsTokens)
	}
}

func buildTagsTokens(tags []tags.ITag) hclwrite.Tokens {
	tagsMap := make(map[string]cty.Value, len(tags))
	for _, tag := range tags {
		tagsMap[tag.GetKey()] = cty.StringVal(tag.GetValue())
	}
	if len(tagsMap) > 0 {
		return hclwrite.TokensForValue(cty.MapVal(tagsMap))
	}
	return nil
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

// Inserts a list of tags at end of list
func InsertTokens(tokens hclwrite.Tokens, values []*hclwrite.Token) hclwrite.Tokens {
	suffix := tokens[len(tokens)-2:]
	result := append(tokens[:len(tokens)-2], &hclwrite.Token{Type: hclsyntax.TokenNewline, Bytes: []byte("\n")})
	result = append(result, values...)
	return append(result, suffix...)
}

func (p *TerrraformParser) parseBlock(hclBlock *hclwrite.Block) (*TerraformBlock, error) {
	var existingTags []tags.ITag
	isTaggable := false
	var tagsAttributeName string
	if hclBlock.Type() == "resource" {
		resourceType := hclBlock.Labels()[0]
		providerName := getProviderFromResourceType(resourceType)
		if utils.InSlice(SkippedProviders, providerName) {
			return nil, fmt.Errorf("resource belongs to skipped provider %s", providerName)
		}
		client := p.getClient(providerName)
		if client == nil {
			return nil, fmt.Errorf("could not find client of %s", providerName)
		}
		logger.MuteLogging()
		resourceScheme, err := client.GetResourceTypeSchema(resourceType)
		logger.UnmuteLogging()
		if err != nil {
			return nil, err
		}
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
		if isSchemeViolated := p.isSchemeViolated(hclBlock, tagsAttributeName, resourceScheme); isSchemeViolated {
			return nil, nil
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

func (p *TerrraformParser) isSchemeViolated(hclBlock *hclwrite.Block, tagsAttributeName string, resourceScheme *tfschema.Block) bool {
	bodyTokens := hclBlock.Body().BuildTokens(hclwrite.Tokens{})
	foundTagToken := false
	tagTokenRegex := regexp.MustCompile(`^tag[\d]?$`)
	foundTagsToken := false
	tagsTokensRegex := regexp.MustCompile(fmt.Sprintf(`^%s$`, tagsAttributeName))
	for _, token := range bodyTokens {
		if matched := tagTokenRegex.Match(token.Bytes); matched {
			foundTagToken = true
		}
		if matched := tagsTokensRegex.Match(token.Bytes); matched {
			foundTagsToken = true
		}
	}
	if foundTagToken && foundTagsToken {
		if _, okTags := resourceScheme.Attributes[tagsAttributeName]; okTags {
			return true
		}
	}
	return false
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
	prefix := ProviderToTagAttribute[getProviderFromResourceType(resourceType)]
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
		logger.MuteLogging()
		typeSchema, err := client.GetResourceTypeSchema(resourceType)
		logger.UnmuteLogging()
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
		if utils.InSlice(separatorTokens, token.Type) {
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
		_ = json.Unmarshal([]byte(key), &key)
		_ = json.Unmarshal([]byte(value), &value)
		parsedTags[key] = value
	}

	return parsedTags
}

func (p *TerrraformParser) getClient(providerName string) tfschema.Client {
	if utils.InSlice(SkippedProviders, providerName) {
		return nil
	}

	hclLogger := hclog.New(&hclog.LoggerOptions{
		Name:   "plugin",
		Level:  hclog.Error,
		Output: hclog.DefaultOutput,
	})
	client, exists := p.providerToClientMap[providerName]
	if exists {
		return client
	}
	logger.MuteLogging()
	newClient, err := tfschema.NewClient(providerName, tfschema.Option{
		RootDir: p.terraformModule.ProvidersInstallDir,
		Logger:  hclLogger,
	})
	logger.UnmuteLogging()
	if err != nil {
		if strings.Contains(err.Error(), "Failed to find plugin") {
			logger.Warning(fmt.Sprintf("Could not load provider %v, resources from this provider will not be tagged", providerName))
			logger.Warning(fmt.Sprintf("Try to run `terraform init` in the given root dir: [%s] and try again.", p.rootDir))
		}
		return nil
	}

	p.providerToClientMap[providerName] = newClient
	return newClient
}
