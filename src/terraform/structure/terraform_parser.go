package structure

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"sync"

	"github.com/bridgecrewio/yor/src/common"
	"github.com/bridgecrewio/yor/src/common/logger"
	"github.com/bridgecrewio/yor/src/common/structure"
	"github.com/bridgecrewio/yor/src/common/tagging/tags"
	"github.com/bridgecrewio/yor/src/common/utils"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/hashicorp/terraform/command"
	"github.com/minamijoyo/tfschema/tfschema"
	"github.com/zclconf/go-cty/cty"
)

var ignoredDirs = []string{".git", ".DS_Store", ".idea", ".terraform"}
var unsupportedTerraformBlocks = []string{
	"aws_autoscaling_group",                  // This resource specifically supports tags with a different structure, see: https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/autoscaling_group#tag-and-tags
	"aws_lb_listener",                        // This resource does not support tags, although docs state otherwise.
	"aws_lb_listener_rule",                   // This resource does not support tags, although docs state otherwise.
	"aws_cloudwatch_log_destination",         // This resource does not support tags, although docs state otherwise.
	"google_monitoring_notification_channel", //This resource uses labels for other purposes.
	"aws_secretsmanager_secret_rotation",         // This resource does not support tags, although tfschema states otherwise.
}

var taggableResourcesLock sync.RWMutex
var hclWriteLock sync.Mutex

type TerraformParser struct {
	rootDir                string
	providerToClientMap    sync.Map
	taggableResourcesCache map[string]bool
	tagModules             bool
	tagLocalModules        bool
	terraformModule        *TerraformModule
	moduleImporter         *command.GetCommand
	moduleInstallDir       string
	downloadedPaths        []string
	tfClientLock           sync.Mutex
}

func (p *TerraformParser) Name() string {
	return "Terraform"
}

func (p *TerraformParser) Init(rootDir string, args map[string]string) {
	p.rootDir = rootDir
	p.taggableResourcesCache = make(map[string]bool)
	p.tagModules = true
	p.tagLocalModules = false
	p.terraformModule = NewTerraformModule(rootDir)
	if argTagModule, ok := args["tag-modules"]; ok {
		p.tagModules, _ = strconv.ParseBool(argTagModule)
	}

	if argTagLocalModule, ok := args["tag-local-modules"]; ok {
		p.tagLocalModules, _ = strconv.ParseBool(argTagLocalModule)
	}

	p.moduleImporter = &command.GetCommand{Meta: command.Meta{Color: false, Ui: customTfLogger{}}}
	pwd, _ := os.Getwd()
	p.moduleInstallDir = filepath.Join(pwd, ".terraform", "modules")
}

func (p *TerraformParser) Close() {
	logger.MuteOutputBlock(func() {
		p.providerToClientMap.Range(func(provider, iClient interface{}) bool {
			client := iClient.(tfschema.Client)
			client.Close()
			return true
		})
	})
}

func (p *TerraformParser) GetSkippedDirs() []string {
	return ignoredDirs
}

func (p *TerraformParser) GetSupportedFileExtensions() []string {
	return []string{common.TfFileType.Extension}
}

func (p *TerraformParser) GetSourceFiles(directory string) ([]string, error) {
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

func (p *TerraformParser) ValidFile(_ string) bool {
	return true
}

func (p *TerraformParser) ParseFile(filePath string) ([]structure.IBlock, error) {
	// #nosec G304
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
		if !utils.InSlice(SupportedBlockTypes, block.Type()) {
			continue
		}
		blockID := strings.Join(block.Labels(), ".")
		terraformBlock, err := p.parseBlock(block, filePath)
		if err != nil {
			if strings.HasPrefix(err.Error(), "resource belongs to skipped") || strings.HasPrefix(err.Error(), "could not find client") {
				logger.Info(fmt.Sprintf("skipping block %s because the provider %s does not exist locally or does not support tags",
					blockID, strings.Split(blockID, "_")[0]))
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

func (p *TerraformParser) WriteFile(readFilePath string, blocks []structure.IBlock, writeFilePath string) error {
	// #nosec G304
	// read file bytes
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

	tempFile, err := ioutil.TempFile(filepath.Dir(readFilePath), "temp.*.tf")
	if err != nil {
		return err
	}
	fd, err := os.OpenFile(tempFile.Name(), os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	hclWriteLock.Lock()
	defer hclWriteLock.Unlock()
	_, err = hclFile.WriteTo(fd)
	if err != nil {
		return err
	}

	_, err = p.ParseFile(tempFile.Name())
	if err != nil {
		return fmt.Errorf("editing file %v resulted in malformed terraform, please open a github issue with the relevant details", readFilePath)
	}
	err = tempFile.Close()
	if err != nil {
		return err
	}

	//cant delete files on windows if you dont close them
	if err = fd.Close(); err != nil {
		return err
	}
	err = os.Remove(tempFile.Name())
	if err != nil {
		return err
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

func (p *TerraformParser) modifyBlockTags(rawBlock *hclwrite.Block, parsedBlock structure.IBlock) {
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
		for i, rawTagsToken := range rawTagsTokens {
			tokenStr := string(rawTagsToken.Bytes)
			if tokenStr == "merge" {
				isMergeOpExists = true
				break
			}
			if i == 0 && utils.InSlice([]string{VarBlockType, LocalBlockType, ModuleBlockType, DataBlockType}, tokenStr) {
				isRenderedAttribute = true
				break
			}
		}

		var replacedTags []tags.ITag
		var newTags []tags.ITag
		possibleTagKeys := p.extractTagKeysFromRawTokens(rawTagsTokens)
		for _, tag := range mergedTags {
			tagReplaced := false
			strippedTagKey := strings.ReplaceAll(tag.GetKey(), `"`, "")
			for _, t := range possibleTagKeys {
				if t == tag.GetKey() || t == strippedTagKey || strings.Contains(t, strippedTagKey) {
					replacedTags = append(replacedTags, tag)
					tagReplaced = true
					break
				}
			}
			if !tagReplaced {
				newTags = append(newTags, tag)
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

		if len(newTags) == 0 {
			logger.Debug(fmt.Sprintf("Nothing to update for block %v (%v)", parsedBlock.GetResourceID(), parsedBlock.GetFilePath()))
			return
		}

		if !isMergeOpExists && !isRenderedAttribute {
			newTagsTokens := buildTagsTokens(newTags)
			if len(rawTagsTokens) == 1 {
				// The line is:
				//    tags = null
				// => we should replace it!
				rawTagsTokens = newTagsTokens
			} else {
				rawTagsTokens = InsertTokens(rawTagsTokens, newTagsTokens[2:len(newTagsTokens)-2])
			}
			rawBlock.Body().SetAttributeRaw(tagsAttributeName, rawTagsTokens)
			return
		}

		// These lines execute if there is either a `merge` operator at the start of the tags,
		// or if it is rendered via a variable / local.
		newTagsTokens := buildTagsTokens(newTags)
		if !isMergeOpExists && newTagsTokens != nil {
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
		if newTagsTokens != nil {
			if rawTagsTokens[len(rawTagsTokens)-3].Type != hclsyntax.TokenComma &&
				rawTagsTokens[len(rawTagsTokens)-2].Type != hclsyntax.TokenComma {
				// Insert a comma token before the merge closing parenthesis and add as a separate dict
				rawTagsTokens = InsertToken(rawTagsTokens, len(rawTagsTokens)-1, &hclwrite.Token{
					Type:  hclsyntax.TokenComma,
					Bytes: []byte(","),
				})
			}

			for _, tagToken := range newTagsTokens {
				rawTagsTokens = InsertToken(rawTagsTokens, len(rawTagsTokens)-1, tagToken)
			}
		}
		// Set the body's tags to the new built tokens
		rawBlock.Body().SetAttributeRaw(tagsAttributeName, rawTagsTokens)
	}
}

func (p *TerraformParser) extractTagKeysFromRawTokens(rawTagsTokens hclwrite.Tokens) []string {
	var tokens []string
	for _, t := range rawTagsTokens {
		tokens = append(tokens, string(t.Bytes))
	}
	var possibleTagKeys []string
	var currentToken string
	var inInterpolation bool
	for _, t := range tokens {
		switch t {
		case "{":
			continue
		case "${":
			currentToken = fmt.Sprintf("%v%v", currentToken, t)
			inInterpolation = true
		case "}":
			if inInterpolation {
				currentToken = fmt.Sprintf("%v%v", currentToken, t)
				inInterpolation = false
			} else {
				continue
			}
		case "=", "\n", ",":
			possibleTagKeys = append(possibleTagKeys, currentToken)
			currentToken = ""
		default:
			currentToken = fmt.Sprintf("%v%v", currentToken, t)
		}
	}
	// Cleanup unnecessary quotes around tag name
	for i, t := range possibleTagKeys {
		if strings.HasPrefix(t, "\"") && strings.HasSuffix(t, "\"") && len(t) > 2 {
			possibleTagKeys[i] = t[1 : len(t)-1]
		}
	}
	return possibleTagKeys
}

func buildTagsTokens(tags []tags.ITag) hclwrite.Tokens {
	tagsMap := make(map[string]cty.Value, len(tags))
	for _, tag := range tags {
		tagsMap[tag.GetKey()] = cty.StringVal(tag.GetValue())
	}
	if len(tagsMap) > 0 {
		hclWriteLock.Lock()
		defer hclWriteLock.Unlock()
		return hclwrite.TokensForValue(cty.MapVal(tagsMap))
	}
	return nil
}

// InsertToken Insert inserts a value at a specific index in a slice
func InsertToken(tokens hclwrite.Tokens, index int, value *hclwrite.Token) hclwrite.Tokens {
	if len(tokens) == index {
		return append(tokens, value)
	}
	tokens = append(tokens[:index+1], tokens[index:]...)
	tokens[index] = value
	return tokens
}

// InsertTokens Inserts a list of tags at end of list
func InsertTokens(tokens hclwrite.Tokens, values []*hclwrite.Token) hclwrite.Tokens {
	suffixLength := 1 // Only the closing parenthesis
	if tokens[len(tokens)-2].Type == hclsyntax.TokenNewline {
		suffixLength = 2
	}
	var result hclwrite.Tokens
	result = append(result, tokens[:len(tokens)-suffixLength]...)
	result = append(result, &hclwrite.Token{Type: hclsyntax.TokenNewline, Bytes: []byte("\n")})
	result = append(result, values...)
	if suffixLength == 1 {
		result = append(result, &hclwrite.Token{Type: hclsyntax.TokenNewline, Bytes: []byte("\n")})
	}
	return append(result, tokens[len(tokens)-suffixLength:]...)
}

func (p *TerraformParser) parseBlock(hclBlock *hclwrite.Block, filePath string) (*TerraformBlock, error) {
	var existingTags []tags.ITag
	isTaggable := false
	var tagsAttributeName string
	var resourceType string
	var err error

	switch hclBlock.Type() {
	case ResourceBlockType:
		resourceType = hclBlock.Labels()[0]
		providerName := getProviderFromResourceType(resourceType)
		if utils.InSlice(SkippedProviders, providerName) {
			return nil, fmt.Errorf("resource belongs to skipped provider %s", providerName)
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
	case ModuleBlockType:
		resourceType = "module"
		defer func() {
			if e := recover(); e != nil {
				logger.Warning(fmt.Sprintf("Failed to parse module module.%v (%v)", strings.Join(hclBlock.Labels(), "."), filePath))
				err = fmt.Errorf("failed to parse module.%v", strings.Join(hclBlock.Labels(), "."))
			}
		}()
		isTaggable, existingTags, tagsAttributeName = p.extractTagsFromModule(hclBlock, filePath, isTaggable, existingTags, tagsAttributeName)
	}

	terraformBlock := TerraformBlock{
		Block: structure.Block{
			ExitingTags:       existingTags,
			IsTaggable:        isTaggable,
			TagsAttributeName: tagsAttributeName,
			Type:              resourceType,
		},
	}

	return &terraformBlock, err
}

func (p *TerraformParser) extractTagsFromModule(hclBlock *hclwrite.Block, filePath string, isTaggable bool, existingTags []tags.ITag, tagsAttributeName string) (bool, []tags.ITag, string) {
	moduleSource := string(hclBlock.Body().GetAttribute("source").Expr().BuildTokens(hclwrite.Tokens{}).Bytes())
	// source is always wrapped in " front and back
	moduleSource = strings.Trim(moduleSource, "\" ")

	if !isRemoteModule(moduleSource) && !isTerraformRegistryModule(moduleSource) && !p.tagLocalModules {
		// Don't use the tags label on local modules - the underlying resources will be tagged by themselves
		isTaggable = false
	} else {
		// This is a remote module - if it has tags attribute, tag it!
		moduleProvider := ExtractProviderFromModuleSrc(moduleSource)
		possibleTagAttributeNames := []string{"extra_tags", "tags", "common_tags", "labels"}
		if val, ok := ProviderToTagAttribute[moduleProvider]; ok {
			possibleTagAttributeNames = append(possibleTagAttributeNames, val)
		}
		for _, tan := range possibleTagAttributeNames {
			existingTags, isTaggable = p.getModuleTags(hclBlock, tan)

			if isTaggable {
				tagsAttributeName = tan
				break
			}
		}
		if !isTaggable {
			isTaggable, tagsAttributeName = p.isModuleTaggable(filePath, strings.Join(hclBlock.Labels(), "."), possibleTagAttributeNames)
		}
	}
	return isTaggable, existingTags, tagsAttributeName
}

func ExtractProviderFromModuleSrc(source string) string {
	if strings.HasPrefix(source, "app.terraform.io") {
		// Terraform modules in private registry follow this structure: <HOSTNAME>/<ORGANIZATION>/<MODULE NAME>/<PROVIDER>
		// https://www.terraform.io/docs/cloud/registry/using.html
		return strings.Split(source, "/")[3]
	}
	if isTerraformRegistryModule(source) {
		matches := utils.FindSubMatchByGroup(RegistryModuleRegex, source)
		val, _ := matches["PROVIDER"]
		return val
	}
	withoutRef := strings.Split(source, "//")[0]
	parts := strings.Split(strings.TrimRight(withoutRef, ".git"), "/")
	for _, part := range parts {
		if strings.HasPrefix(part, "terraform-") {
			return strings.Split(part, "-")[1]
		}
	}
	return ""
}

func (p *TerraformParser) isModuleTaggable(fp string, moduleName string, tagAtts []string) (bool, string) {
	logger.Info(fmt.Sprintf("Searching module %v for %v", moduleName, tagAtts))
	actualPath, _ := filepath.Rel(p.rootDir, filepath.Dir(fp))
	absRootPath, _ := filepath.Abs(p.rootDir)
	actualPath, _ = filepath.Abs(filepath.Join(absRootPath, actualPath))
	if !utils.InSlice(p.downloadedPaths, fp) && os.Getenv("YOR_DISABLE_TF_MODULE_DOWNLOAD") != "TRUE" {
		logger.MuteOutputBlock(func() {
			logger.Info(fmt.Sprintf("Downloading modules for dir %v\n", actualPath))
			_ = p.moduleImporter.Run([]string{actualPath})
			p.downloadedPaths = append(p.downloadedPaths, fp)
		})
	}
	expectedModuleDir := filepath.Join(p.moduleInstallDir, moduleName)
	if _, err := os.Stat(expectedModuleDir); os.IsNotExist(err) {
		return false, ""
	}

	files, _ := ioutil.ReadDir(expectedModuleDir)
	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".tf") {
			blocks, _ := p.ParseFile(filepath.Join(expectedModuleDir, f.Name()))
			for _, b := range blocks {
				if b.(*TerraformBlock).HclSyntaxBlock.Type == VariableBlockType {
					for _, tagAtt := range tagAtts {
						if b.GetResourceID() == tagAtt {
							return true, tagAtt
						}
					}
				}
			}
		}
	}

	return false, ""
}

func (p *TerraformParser) getTagsAttributeName(hclBlock *hclwrite.Block) (string, error) {
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

func (p *TerraformParser) getExistingTags(hclBlock *hclwrite.Block, tagsAttributeName string) ([]tags.ITag, bool) {
	isTaggable := false
	existingTags := make([]tags.ITag, 0)

	tagsAttribute := hclBlock.Body().GetAttribute(tagsAttributeName)
	if tagsAttribute != nil {
		// if tags exists in resource
		isTaggable, _ = p.isBlockTaggable(hclBlock)
		tagsTokens := tagsAttribute.Expr().BuildTokens(hclwrite.Tokens{})
		parsedTags := p.parseTagAttribute(tagsTokens)
		for key := range parsedTags {
			iTag := tags.Init(key, parsedTags[key])
			existingTags = append(existingTags, iTag)
		}
	}

	return existingTags, isTaggable
}

func (p *TerraformParser) isBlockTaggable(hclBlock *hclwrite.Block) (bool, error) {
	resourceType := hclBlock.Labels()[0]
	if utils.InSlice(unsupportedTerraformBlocks, resourceType) {
		return false, nil
	}
	if utils.InSlice(TfTaggableResourceTypes, resourceType) {
		return true, nil
	}
	taggableResourcesLock.RLock()
	val, ok := p.taggableResourcesCache[resourceType]
	taggableResourcesLock.RUnlock()
	if ok {
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
		var typeSchema *tfschema.Block
		logger.MuteOutputBlock(func() {
			typeSchema, err = client.GetResourceTypeSchema(resourceType)
		})
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
	taggableResourcesLock.Lock()
	p.taggableResourcesCache[resourceType] = taggable
	taggableResourcesLock.Unlock()
	return taggable, nil
}

func (p *TerraformParser) getHclMapsContents(tokens hclwrite.Tokens) []hclwrite.Tokens {
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

func (p *TerraformParser) extractTagPairs(tokens hclwrite.Tokens) []hclwrite.Tokens {
	// The function gets tokens and returns an array of tokens that represent key and value
	// example: tokens: "a=1\n b=2, c=3", returns: ["a=1", "b=2", "c=3"]
	separatorTokens := []hclsyntax.TokenType{hclsyntax.TokenComma, hclsyntax.TokenNewline}

	bracketsCounters := map[hclsyntax.TokenType]int{
		hclsyntax.TokenOParen: 0,
		hclsyntax.TokenOBrack: 0,
	}

	openingBrackets := []hclsyntax.TokenType{hclsyntax.TokenOParen, hclsyntax.TokenOBrack}
	closingBrackets := []hclsyntax.TokenType{hclsyntax.TokenCParen, hclsyntax.TokenCBrack}

	bracketsPairs := map[hclsyntax.TokenType]hclsyntax.TokenType{
		hclsyntax.TokenCParen: hclsyntax.TokenOParen,
		hclsyntax.TokenCBrack: hclsyntax.TokenOBrack,
	}

	tagPairs := make([]hclwrite.Tokens, 0)
	startIndex := 0
	hasEq := false
	for i, token := range tokens {
		if utils.InSlice(separatorTokens, token.Type) && getUncloseBracketsCount(bracketsCounters) == 0 {
			if hasEq {
				tagPairs = append(tagPairs, tokens[startIndex:i])
			}
			startIndex = i + 1
			hasEq = false
		}
		if token.Type == hclsyntax.TokenEqual {
			hasEq = true
		}
		if utils.InSlice(openingBrackets, token.Type) {
			bracketsCounters[token.Type]++
		}
		if utils.InSlice(closingBrackets, token.Type) {
			matchingOpen := bracketsPairs[token.Type]
			bracketsCounters[matchingOpen]--
		}
	}
	if hasEq {
		tagPairs = append(tagPairs, tokens[startIndex:])
	}

	return tagPairs
}

func getUncloseBracketsCount(bracketsCounters map[hclsyntax.TokenType]int) int {
	sum := 0
	for b := range bracketsCounters {
		sum += bracketsCounters[b]
	}

	return sum
}

func (p *TerraformParser) parseTagAttribute(tokens hclwrite.Tokens) map[string]string {
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

func (p *TerraformParser) getClient(providerName string) tfschema.Client {
	if utils.InSlice(SkippedProviders, providerName) {
		return nil
	}

	p.tfClientLock.Lock()
	defer p.tfClientLock.Unlock()

	client, exists := p.providerToClientMap.Load(providerName)
	if exists {
		return client.(tfschema.Client)
	}

	hclLogger := hclog.New(&hclog.LoggerOptions{
		Name:   "plugin",
		Level:  hclog.Error,
		Output: hclog.DefaultOutput,
	})
	var err error
	var newClient tfschema.Client
	if p.terraformModule == nil {
		logger.Warning(fmt.Sprintf("Failed to initialize terraform module, it might be due to a malformed file in the given root dir: [%s]", p.rootDir))
		return nil
	}
	logger.MuteOutputBlock(func() {
		newClient, err = tfschema.NewClient(providerName, tfschema.Option{
			RootDir: p.terraformModule.ProvidersInstallDir,
			Logger:  hclLogger,
		})
	})
	if err != nil {
		if strings.Contains(err.Error(), "Failed to find plugin") {
			logger.Warning(fmt.Sprintf("Could not load provider %v, resources from this provider will not be tagged", providerName))
			logger.Warning(fmt.Sprintf("Try to run `terraform init` in the given root dir: [%s] and try again.", p.rootDir))
		}
		return nil
	}

	p.providerToClientMap.Store(providerName, newClient)
	return newClient
}

func (p *TerraformParser) getModuleTags(hclBlock *hclwrite.Block, tagsAttributeName string) ([]tags.ITag, bool) {
	isTaggable := false
	existingTags := make([]tags.ITag, 0)

	tagsAttribute := hclBlock.Body().GetAttribute(tagsAttributeName)
	if tagsAttribute != nil {
		// if tags exists in module
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
