package common

import (
	"bridgecrewio/yor/src/common/logger"
	"fmt"
	"strings"

	"gopkg.in/validator.v2"
)

var allowedOutputTypes = []string{"cli", "json"}

type TagGroupName string

const (
	SimpleTagGroupName TagGroupName = "simple"
	GitTagGroupName    TagGroupName = "git"
	Code2Cloud         TagGroupName = "code2cloud"
)

var TagGroupNames = []TagGroupName{SimpleTagGroupName, GitTagGroupName, Code2Cloud}

type TagOptions struct {
	Directory      string
	Tag            string
	SkipTags       []string
	CustomTagging  []string
	SkipDirs       []string
	Output         string `validate:"output"`
	OutputJSONFile string
}

type ListTagsOptions struct {
	TagGroups []string `validate:"tagGroupNames"`
}

func (o *TagOptions) Validate() {
	_ = validator.SetValidationFunc("output", validateOutput)
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
	val, ok := v.([]string)
	if ok {
		for _, gn := range val {
			groupName := TagGroupName(gn)
			if !InSlice(TagGroupNames, groupName) {
				return fmt.Errorf("tag group %s is not one of the supported tag groups. supported groups: %v", gn, TagGroupNames)
			}
		}
		return nil
	}
	return fmt.Errorf("unsupported tag group names [%s]. supported types: %s", val, TagGroupNames)
}

func validateOutput(v interface{}, _ string) error {
	val, ok := v.(string)
	if !ok {
		return validator.ErrUnsupported
	}

	if val != "" && !InSlice(allowedOutputTypes, strings.ToLower(val)) {
		return fmt.Errorf("unsupported output type [%s]. allowed types: %s", val, allowedOutputTypes)
	}

	return nil
}
