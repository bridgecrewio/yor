package structure

import (
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"path/filepath"

	goformationTags "github.com/awslabs/goformation/v4/cloudformation/tags"
	"github.com/bridgecrewio/goformation/v4"
	"github.com/bridgecrewio/goformation/v4/intrinsics"
	"github.com/bridgecrewio/yor/src/common"
	"github.com/bridgecrewio/yor/src/common/json"
	"github.com/bridgecrewio/yor/src/common/logger"
	"github.com/bridgecrewio/yor/src/common/structure"
	"github.com/bridgecrewio/yor/src/common/tagging/tags"
	"github.com/bridgecrewio/yor/src/common/types"
	"github.com/bridgecrewio/yor/src/common/utils"
	"github.com/bridgecrewio/yor/src/common/yaml"

	"reflect"
)

type CloudformationParser struct {
	*types.YamlParser
	*types.JSONParser
}

const TagsAttributeName = "Tags"
const ResourcesStartToken = "Resources"
const EnvVarsPath = "Resources/*/Properties/Environment/Variables/*"

func (p *CloudformationParser) Name() string {
	return "CloudFormation"
}

func (p *CloudformationParser) Init(rootDir string, _ map[string]string) {
	p.YamlParser = &types.YamlParser{
		RootDir:              rootDir,
		FileToResourcesLines: make(map[string]structure.Lines),
	}
	p.JSONParser = &types.JSONParser{
		RootDir:              rootDir,
		FileToBracketMapping: make(map[string]map[int]json.BracketPair),
	}
}

func (p *CloudformationParser) GetSkippedDirs() []string {
	return []string{}
}

func (p *CloudformationParser) GetSupportedFileExtensions() []string {
	return []string{common.YamlFileType.Extension, common.YmlFileType.Extension, common.CFTFileType.Extension, common.JSONFileType.Extension}
}

func (p *CloudformationParser) ParseFile(filePath string) ([]structure.IBlock, error) {
	template, err := goformation.OpenWithOptions(filePath, &intrinsics.ProcessorOptions{
		StringifyPaths: []string{EnvVarsPath},
	})
	if err != nil || template == nil {
		logger.Warning(fmt.Sprintf("There was an error processing the cloudformation template %v: %s", filePath, err))
		if err == nil {
			err = fmt.Errorf("failed to parse template %v", filePath)
		}
		return nil, err
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
			p.FileToBracketMapping[filePath] = fileBracketsMapping
		default:
			return nil, fmt.Errorf("unsupported file type %s", utils.GetFileFormat(filePath))
		}

		minResourceLine := math.MaxInt8
		maxResourceLine := 0
		parsedBlocks := make([]structure.IBlock, 0)
		for resourceName := range template.Resources {
			resource := template.Resources[resourceName]
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
				},
			}
			parsedBlocks = append(parsedBlocks, cfnBlock)
		}

		p.FileToResourcesLines[filePath] = structure.Lines{Start: minResourceLine, End: maxResourceLine}

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
		existingTags = tagsValue.Interface().([]goformationTags.Tag)
	}

	iTags := make([]tags.ITag, 0)
	for _, goformationTag := range existingTags {
		tag := &tags.Tag{Key: goformationTag.Key, Value: goformationTag.Value}
		iTags = append(iTags, tag)
	}

	return iTags
}

func (p *CloudformationParser) WriteFile(readFilePath string, blocks []structure.IBlock, writeFilePath string) error {
	for _, block := range blocks {
		block := block.(*CloudformationBlock)
		block.UpdateTags()
	}
	tempFile, err := ioutil.TempFile(filepath.Dir(readFilePath), "temp.*.template")
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
		return json.WriteJSONFile(readFilePath, blocks, writeFilePath, p.FileToBracketMapping[readFilePath])
	default:
		return fmt.Errorf("unsupported file type %s", utils.GetFileFormat(readFilePath))
	}
}

func (p *CloudformationParser) getTagsLines(filePath string, resourceLinesRange *structure.Lines) structure.Lines {
	nonFoundLines := structure.Lines{Start: -1, End: -1}
	switch utils.GetFileFormat(filePath) {
	case common.YamlFileType.FileFormat, common.YmlFileType.FileFormat:
		file, scanner, _ := utils.GetFileScanner(filePath, &nonFoundLines)
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
		defer func() {
			_ = file.Close()
		}()
		linesInResource, tagsExist := yaml.FindTagsLinesYAML(resourceLinesText, TagsAttributeName)
		if tagsExist {
			return structure.Lines{Start: linesInResource.Start + resourceLinesRange.Start, End: linesInResource.Start + resourceLinesRange.Start + (linesInResource.End - linesInResource.Start)}
		}
		return structure.Lines{Start: -1, End: -1}
	case common.JSONFileType.FileFormat:
		// #nosec G304
		file, err := ioutil.ReadFile(filePath)
		if err != nil {
			logger.Warning(fmt.Sprintf("failed to read file %s", filePath))
			return structure.Lines{Start: -1, End: -1}
		}
		tagsBrackets := json.FindScopeInJSON(string(file), TagsAttributeName, p.FileToBracketMapping[filePath], resourceLinesRange)
		tagsLines := &structure.Lines{Start: tagsBrackets.Open.Line, End: tagsBrackets.Close.Line}
		return *tagsLines
	default:
		return structure.Lines{Start: -1, End: -1}
	}
}
