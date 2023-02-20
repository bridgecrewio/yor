package structure

import (
	stdjson "encoding/json"
	"fmt"
	"io"
	"math"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync"

	goformationTags "github.com/awslabs/goformation/v5/cloudformation/tags"
	"github.com/bridgecrewio/goformation/v5"
	"github.com/bridgecrewio/goformation/v5/cloudformation"
	"github.com/bridgecrewio/goformation/v5/intrinsics"
	"github.com/bridgecrewio/yor/src/common"
	"github.com/bridgecrewio/yor/src/common/json"
	"github.com/bridgecrewio/yor/src/common/logger"
	"github.com/bridgecrewio/yor/src/common/structure"
	"github.com/bridgecrewio/yor/src/common/tagging/tags"
	"github.com/bridgecrewio/yor/src/common/types"
	"github.com/bridgecrewio/yor/src/common/utils"
	"github.com/bridgecrewio/yor/src/common/yaml"
	sanathyaml "github.com/sanathkr/yaml"
)

type CloudformationParser struct {
	*types.YamlParser
	*types.JSONParser
}

const TagsAttributeName = "Tags"
const ResourcesStartToken = "Resources"
const EnvVarsPath = "Resources/*/Properties/Environment/Variables/*"

var goformationLock sync.Mutex

func (p *CloudformationParser) Name() string {
	return "CloudFormation"
}

func (p *CloudformationParser) Init(rootDir string, _ map[string]string) {
	p.YamlParser = &types.YamlParser{
		RootDir: rootDir,
	}
	p.JSONParser = &types.JSONParser{
		RootDir: rootDir,
	}
}

func (p CloudformationParser) Close() {
}

func (p *CloudformationParser) GetSkippedDirs() []string {
	return []string{}
}

func (p *CloudformationParser) GetSupportedFileExtensions() []string {
	return []string{common.YamlFileType.Extension, common.YmlFileType.Extension, common.CFTFileType.Extension, common.JSONFileType.Extension}
}

// ValidFile Validate file has AWSTemplateFormatVersion
func (p *CloudformationParser) ValidFile(filePath string) bool {
	// #nosec G304
	file, err := os.Open(filePath)
	if err != nil {
		logger.Warning(fmt.Sprintf("Error opening file %s, skipping: %v", filePath, err))
		return false
	}
	bytes, err := io.ReadAll(file)
	if err != nil {
		logger.Warning(fmt.Sprintf("Error reading file %s, skipping: %v", filePath, err))
		return false
	}
	if err = file.Close(); err != nil {
		logger.Warning(fmt.Sprintf("Error closing file %s, skipping: %v", filePath, err))
		return false
	}

	if !strings.HasSuffix(filePath, ".json") {
		bytes, err = sanathyaml.YAMLToJSON(bytes)
		if err != nil {
			logger.Warning(fmt.Sprintf("Error converting YAML to JSON for file %s, skipping: %v", filePath, err))
			return false
		}
	}
	var result map[string]interface{}
	err = stdjson.Unmarshal(bytes, &result)
	if err != nil {
		logger.Warning(fmt.Sprintf("Error unmarshalling JSON for file %s, skipping: %v", filePath, err))
		return false
	}
	_, hasHeader := result["AWSTemplateFormatVersion"]
	return hasHeader
}

func goformationParse(file string) (*cloudformation.Template, error) {
	var template *cloudformation.Template
	var err error
	defer func() {
		if e := recover(); e != nil {
			logger.Warning(fmt.Sprintf("Failed to parser cfn file at %v due to: %v", file, e))
			err = fmt.Errorf("failed to parse cfn file %v: %v", file, e)
		}
	}()

	template, err = goformation.OpenWithOptions(file, &intrinsics.ProcessorOptions{
		StringifyPaths: []string{EnvVarsPath},
	})
	return template, err
}

func (p *CloudformationParser) ParseFile(filePath string) ([]structure.IBlock, error) {
	goformationLock.Lock()
	template, err := goformationParse(filePath)
	goformationLock.Unlock()
	if err != nil || template == nil {
		logger.Warning(fmt.Sprintf("There was an error processing the cloudformation template %v: %s", filePath, err))
		if err == nil {
			err = fmt.Errorf("failed to parse template %v", filePath)
		}
		return nil, err
	}

	if template.Transform != nil {
		logger.Info(fmt.Sprintf("Skipping CFN template %s as SAM templates are not yet supported", filePath))
		return nil, nil
	}

	resourceNames := make([]string, 0)
	if template.Resources != nil && len(template.Resources) > 0 {
		for resourceName := range template.Resources {
			resourceNames = append(resourceNames, resourceName)
		}

		var resourceNamesToLines map[string]*structure.Lines
		switch utils.GetFileFormat(filePath) {
		case common.YmlFileType.FileFormat, common.YamlFileType.FileFormat:
			resourceNamesToLines = yaml.MapResourcesLineYAML(filePath, resourceNames, ResourcesStartToken)
		case common.JSONFileType.FileFormat:
			var fileBracketsMapping map[int]json.BracketPair
			resourceNamesToLines, fileBracketsMapping = json.MapResourcesLineJSON(filePath, resourceNames)
			p.FileToBracketMapping.Store(filePath, fileBracketsMapping)
		default:
			return nil, fmt.Errorf("unsupported file type %s", utils.GetFileFormat(filePath))
		}

		minResourceLine := math.MaxInt8
		maxResourceLine := 0
		parsedBlocks := make([]structure.IBlock, 0)
		for resourceName := range template.Resources {
			resource := template.Resources[resourceName]
			resourceType := resource.AWSCloudFormationType()
			lines := resourceNamesToLines[resourceName]
			isTaggable, tagsValue := utils.StructContainsProperty(resource, TagsAttributeName)
			tagsLines := structure.Lines{Start: -1, End: -1}
			var existingTags []tags.ITag
			if isTaggable {
				tagsLines, existingTags = p.extractTagsAndLines(filePath, lines, tagsValue)
			}
			minResourceLine = int(math.Min(float64(minResourceLine), float64(lines.Start)))
			maxResourceLine = int(math.Max(float64(maxResourceLine), float64(lines.End)))

			cfnBlock := &CloudformationBlock{
				Block: structure.Block{
					FilePath:          filePath,
					ExitingTags:       existingTags,
					RawBlock:          resource,
					IsTaggable:        isTaggable,
					TagsAttributeName: TagsAttributeName,
					Lines:             *lines,
					TagLines:          tagsLines,
					Name:              resourceName,
					Type:              resourceType,
				},
			}
			parsedBlocks = append(parsedBlocks, cfnBlock)
		}

		p.FileToResourcesLines.Store(filePath, structure.Lines{Start: minResourceLine, End: maxResourceLine})

		return parsedBlocks, nil
	}
	return nil, err
}

func (p *CloudformationParser) extractTagsAndLines(filePath string, lines *structure.Lines, tagsValue reflect.Value) (structure.Lines, []tags.ITag) {
	tagsLines := p.getTagsLines(filePath, lines)
	existingTags := p.GetExistingTags(tagsValue)
	return tagsLines, existingTags
}

func (p *CloudformationParser) GetExistingTags(tagsValue reflect.Value) []tags.ITag {
	existingTags := make([]goformationTags.Tag, 0)
	if tagsValue.Kind() == reflect.Slice {
		ok := true
		existingTags, ok = tagsValue.Interface().([]goformationTags.Tag)
		if !ok {
			for i := 0; i < tagsValue.Len(); i++ {
				iTag := tagsValue.Index(i)

				hasKey, tagKey := utils.StructContainsProperty(iTag.Interface(), "Key")
				hasValue, tagValue := utils.StructContainsProperty(iTag.Interface(), "Value")
				if hasKey && hasValue {
					existingTag := goformationTags.Tag{Key: tagKey.String(), Value: tagValue.String()}
					existingTags = append(existingTags, existingTag)
				}
			}
		}
	}

	iTags := make([]tags.ITag, 0)
	for _, goformationTag := range existingTags {
		tag := &tags.Tag{Key: goformationTag.Key, Value: goformationTag.Value}
		iTags = append(iTags, tag)
	}

	return iTags
}

func (p *CloudformationParser) WriteFile(readFilePath string, blocks []structure.IBlock, writeFilePath string, addToggle bool) error {
	for _, block := range blocks {
		block := block.(*CloudformationBlock)
		block.UpdateTags()
	}
	tempFile, err := os.CreateTemp(filepath.Dir(readFilePath), "temp.*.template")
	defer func() {
		_ = os.Remove(tempFile.Name())
	}()
	if err != nil {
		return err
	}
	err = p.writeToFile(readFilePath, blocks, tempFile.Name())
	if err != nil {
		return err
	}

	_, err = p.ParseFile(tempFile.Name())
	if err != nil {
		return fmt.Errorf("editing file %v resulted in a malformed template, please open a github issue with the relevant details", readFilePath)
	}
	return p.writeToFile(readFilePath, blocks, writeFilePath)
}

func (p *CloudformationParser) writeToFile(readFilePath string, blocks []structure.IBlock, writeFilePath string) error {
	switch utils.GetFileFormat(readFilePath) {
	case common.YamlFileType.FileFormat, common.YmlFileType.FileFormat:
		return yaml.WriteYAMLFile(readFilePath, blocks, writeFilePath, TagsAttributeName, ResourcesStartToken)
	case common.JSONFileType.FileFormat:
		bracketMapping, _ := p.FileToBracketMapping.Load(readFilePath)
		return json.WriteJSONFile(readFilePath, blocks, writeFilePath, bracketMapping.(map[int]json.BracketPair))
	default:
		return fmt.Errorf("unsupported file type %s", utils.GetFileFormat(readFilePath))
	}
}

func (p *CloudformationParser) getTagsLines(filePath string, resourceLinesRange *structure.Lines) structure.Lines {
	nonFoundLines := structure.Lines{Start: -1, End: -1}
	switch utils.GetFileFormat(filePath) {
	case common.YamlFileType.FileFormat, common.YmlFileType.FileFormat:
		scanner, _ := utils.GetFileScanner(filePath, &nonFoundLines)
		resourceLinesText := make([]string, 0)
		// iterate file line by line
		lineCounter := 0
		for scanner.Scan() {
			if lineCounter > resourceLinesRange.End {
				break
			}
			if lineCounter >= resourceLinesRange.Start && lineCounter <= resourceLinesRange.End {
				resourceLinesText = append(resourceLinesText, scanner.Text())
			}
			lineCounter++
		}
		linesInResource, tagsExist := yaml.FindTagsLinesYAML(resourceLinesText, TagsAttributeName)
		if tagsExist {
			return structure.Lines{Start: linesInResource.Start + resourceLinesRange.Start, End: linesInResource.Start + resourceLinesRange.Start + (linesInResource.End - linesInResource.Start)}
		}
		return structure.Lines{Start: -1, End: -1}
	case common.JSONFileType.FileFormat:
		// #nosec G304
		file, err := os.ReadFile(filePath)
		if err != nil {
			logger.Warning(fmt.Sprintf("failed to read file %s", filePath))
			return structure.Lines{Start: -1, End: -1}
		}
		bracketMapping, _ := p.FileToBracketMapping.Load(filePath)
		tagsBrackets := json.FindScopeInJSON(string(file), TagsAttributeName, bracketMapping.(map[int]json.BracketPair), resourceLinesRange)
		tagsLines := &structure.Lines{Start: tagsBrackets.Open.Line, End: tagsBrackets.Close.Line}
		return *tagsLines
	default:
		return structure.Lines{Start: -1, End: -1}
	}
}
