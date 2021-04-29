package cli

import (
	"bridgecrewio/yor/src/common/logger"
	taggingUtils "bridgecrewio/yor/src/common/tagging/utils"
	"bridgecrewio/yor/src/common/utils"
	"fmt"
	"strings"

	"gopkg.in/validator.v2"
)

var allowedOutputTypes = []string{"cli", "json"}

type TagOptions struct {
	Directory      string
	Tag            string
	SkipTags       []string
	CustomTagging  []string
	SkipDirs       []string
	Output         string `validate:"output"`
	OutputJSONFile string
	TagGroups      []string `validate:"tagGroupNames"`
}

type ListTagsOptions struct {
	TagGroups []string `validate:"tagGroupNames"`
}

func (o *TagOptions) Validate() {
	_ = validator.SetValidationFunc("output", validateOutput)
	_ = validator.SetValidationFunc("tagGroupNames", validateTagGroupNames)
	if err := validator.Validate(o); err != nil {
		logger.Error(err.Error())
	}
}

func (l *ListTagsOptions) Validate() {
	_ = validator.SetValidationFunc("tagGroupNames", validateTagGroupNames)
	if err := validator.Validate(l); err != nil {
		logger.Error(err.Error())
	}
}

func validateTagGroupNames(v interface{}, _ string) error {
	tagGroupsNames := taggingUtils.GetAllTagGroupsNames()
	val, ok := v.([]string)
	if ok {
		for _, gn := range val {
			if !utils.InSlice(tagGroupsNames, gn) {
				return fmt.Errorf("tag group %s is not one of the supported tag groups. supported groups: %v", gn, tagGroupsNames)
			}
		}
		return nil
	}
	return fmt.Errorf("unsupported tag group names [%s]. supported types: %s", val, tagGroupsNames)
}

func validateOutput(v interface{}, _ string) error {
	val, ok := v.(string)
	if !ok {
		return validator.ErrUnsupported
	}

	if val != "" && !utils.InSlice(allowedOutputTypes, strings.ToLower(val)) {
		return fmt.Errorf("unsupported output type [%s]. allowed types: %s", val, allowedOutputTypes)
	}

	return nil
}
