package common

import (
	"bridgecrewio/yor/src/common/logger"
	"encoding/json"
	"fmt"
	"strings"

	"gopkg.in/validator.v2"
)

var allowedOutputTypes = []string{"cli", "json"}
var SupportedTaggers = []string{"simple", "code2cloud", "git"}

type TagOptions struct {
	Directory              string
	Tag                    string
	SkipTag                string
	CustomTaggers          []string
	SkipConfigurationPaths []string
	Output                 string `validate:"output"`
	OutputJSONFile         string
	ExtraTags              string `validate:"extraTags"`
}

type DescribeTaggerOptions struct {
	Tagger string `validate:"taggerType"`
}

func (o *TagOptions) Validate() {
	_ = validator.SetValidationFunc("extraTags", validateExtraTags)
	_ = validator.SetValidationFunc("output", validateOutput)
	if err := validator.Validate(o); err != nil {
		logger.Error(err.Error())
	}
}

func (o *DescribeTaggerOptions) Validate() {
	_ = validator.SetValidationFunc("taggerType", validateTaggerType)
}

func validateExtraTags(v interface{}, _ string) error {
	val, ok := v.(string)
	if !ok {
		return validator.ErrUnsupported
	}

	var extraTagMap map[string]string
	if err := json.Unmarshal([]byte(val), &extraTagMap); err != nil {
		return fmt.Errorf("failed to parse extra tags: %s", err)
	}

	return nil
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

func validateTaggerType(v interface{}, _ string) error {
	val, ok := v.(string)
	if !ok {
		return validator.ErrUnsupported
	}

	if val != "" && !InSlice(SupportedTaggers, strings.ToLower(val)) {
		return fmt.Errorf("unsupported output type [%s]. allowed types: %s", val, allowedOutputTypes)
	}

	return nil
}
