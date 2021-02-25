package structure

import (
	"bridgecrewio/yor/src/common"
	"bridgecrewio/yor/src/common/logger"
	"bridgecrewio/yor/src/common/structure"
	"bridgecrewio/yor/src/common/tagging/tags"
	"fmt"
	"math"

	"github.com/awslabs/goformation/v4"
	goformation_tags "github.com/awslabs/goformation/v4/cloudformation/tags"

	"reflect"
)

const TagsAttributeName = "Tags"

var TemplateSections = []string{"AWSTemplateFormatVersion", "Transform", "Description", "Metadata", "Parameters", "Mappings", "Conditions", "Outputs", "Resources"}

type CloudformationParser struct {
	rootDir              string
	fileToResourcesLines map[string]common.Lines
}

func (p *CloudformationParser) Init(rootDir string, _ map[string]string) {
	p.rootDir = rootDir
	p.fileToResourcesLines = make(map[string]common.Lines)
}

func (p *CloudformationParser) GetSkippedDirs() []string {
	return []string{}
}

func (p *CloudformationParser) GetAllowedFileTypes() []string {
	return []string{".yaml"}
}

func (p *CloudformationParser) ParseFile(filePath string) ([]structure.IBlock, error) {
	template, err := goformation.Open(filePath)
	if err != nil {
		logger.Warning(fmt.Sprintf("There was an error processing the cloudformation template: %s", err))
	}

	resourceNames := make([]string, 0)
	for resourceName := range template.Resources {
		resourceNames = append(resourceNames, resourceName)
	}

	var resourceNamesToLines map[string]*common.Lines
	switch common.GetFileFormat(filePath) {
	case "yaml":
		resourceNamesToLines = MapResourcesLineYAML(filePath, resourceNames)
	default:
		return nil, fmt.Errorf("unsupported file type %s", common.GetFileFormat(filePath))
	}

	minResourceLine := math.MaxInt8
	maxResourceLine := 0
	parsedBlocks := make([]structure.IBlock, 0)
	for resourceName := range template.Resources {
		resource := template.Resources[resourceName]
		isTaggable, tagsValue := common.StructContainsProperty(resource, TagsAttributeName)
		var existingTags []tags.ITag
		if isTaggable {
			existingTags = p.GetExistingTags(tagsValue)
		}
		lines := resourceNamesToLines[resourceName]
		minResourceLine = int(math.Min(float64(minResourceLine), float64(lines.Start)))
		maxResourceLine = int(math.Max(float64(maxResourceLine), float64(lines.End)))

		cfnBlock := &CloudformationBlock{
			Block: structure.Block{
				FilePath:          filePath,
				ExitingTags:       existingTags,
				RawBlock:          resource,
				IsTaggable:        isTaggable,
				TagsAttributeName: TagsAttributeName,
			},
			lines: *lines,
			name:  resourceName,
		}
		parsedBlocks = append(parsedBlocks, cfnBlock)
	}

	p.fileToResourcesLines[filePath] = common.Lines{Start: minResourceLine, End: maxResourceLine}

	return parsedBlocks, nil
}

func (p *CloudformationParser) GetExistingTags(tagsValue reflect.Value) []tags.ITag {
	existingTags := make([]goformation_tags.Tag, 0)
	if tagsValue.Kind() == reflect.Slice {
		existingTags = tagsValue.Interface().([]goformation_tags.Tag)
	}

	iTags := make([]tags.ITag, 0)
	for _, goformationTag := range existingTags {
		tag := &tags.Tag{Key: goformationTag.Key, Value: goformationTag.Value}
		iTags = append(iTags, tag)
	}

	return iTags
}

func (p *CloudformationParser) WriteFile(readFilePath string, blocks []structure.IBlock, writeFilePath string) error {
	fileFormat := common.GetFileFormat(readFilePath)
	switch fileFormat {
	case "yaml":
		return WriteYAMLFile(readFilePath, blocks, writeFilePath, p.fileToResourcesLines[readFilePath])
	default:
		logger.Warning(fmt.Sprintf("unsupported file type %s", fileFormat))
		return nil
	}
}
