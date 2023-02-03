package quicksight

import (
	"github.com/awslabs/goformation/v4/cloudformation/policies"
)

// Theme_UIColorPalette AWS CloudFormation Resource (AWS::QuickSight::Theme.UIColorPalette)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-quicksight-theme-uicolorpalette.html
type Theme_UIColorPalette struct {

	// Accent AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-quicksight-theme-uicolorpalette.html#cfn-quicksight-theme-uicolorpalette-accent
	Accent string `json:"Accent,omitempty"`

	// AccentForeground AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-quicksight-theme-uicolorpalette.html#cfn-quicksight-theme-uicolorpalette-accentforeground
	AccentForeground string `json:"AccentForeground,omitempty"`

	// Danger AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-quicksight-theme-uicolorpalette.html#cfn-quicksight-theme-uicolorpalette-danger
	Danger string `json:"Danger,omitempty"`

	// DangerForeground AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-quicksight-theme-uicolorpalette.html#cfn-quicksight-theme-uicolorpalette-dangerforeground
	DangerForeground string `json:"DangerForeground,omitempty"`

	// Dimension AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-quicksight-theme-uicolorpalette.html#cfn-quicksight-theme-uicolorpalette-dimension
	Dimension string `json:"Dimension,omitempty"`

	// DimensionForeground AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-quicksight-theme-uicolorpalette.html#cfn-quicksight-theme-uicolorpalette-dimensionforeground
	DimensionForeground string `json:"DimensionForeground,omitempty"`

	// Measure AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-quicksight-theme-uicolorpalette.html#cfn-quicksight-theme-uicolorpalette-measure
	Measure string `json:"Measure,omitempty"`

	// MeasureForeground AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-quicksight-theme-uicolorpalette.html#cfn-quicksight-theme-uicolorpalette-measureforeground
	MeasureForeground string `json:"MeasureForeground,omitempty"`

	// PrimaryBackground AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-quicksight-theme-uicolorpalette.html#cfn-quicksight-theme-uicolorpalette-primarybackground
	PrimaryBackground string `json:"PrimaryBackground,omitempty"`

	// PrimaryForeground AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-quicksight-theme-uicolorpalette.html#cfn-quicksight-theme-uicolorpalette-primaryforeground
	PrimaryForeground string `json:"PrimaryForeground,omitempty"`

	// SecondaryBackground AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-quicksight-theme-uicolorpalette.html#cfn-quicksight-theme-uicolorpalette-secondarybackground
	SecondaryBackground string `json:"SecondaryBackground,omitempty"`

	// SecondaryForeground AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-quicksight-theme-uicolorpalette.html#cfn-quicksight-theme-uicolorpalette-secondaryforeground
	SecondaryForeground string `json:"SecondaryForeground,omitempty"`

	// Success AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-quicksight-theme-uicolorpalette.html#cfn-quicksight-theme-uicolorpalette-success
	Success string `json:"Success,omitempty"`

	// SuccessForeground AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-quicksight-theme-uicolorpalette.html#cfn-quicksight-theme-uicolorpalette-successforeground
	SuccessForeground string `json:"SuccessForeground,omitempty"`

	// Warning AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-quicksight-theme-uicolorpalette.html#cfn-quicksight-theme-uicolorpalette-warning
	Warning string `json:"Warning,omitempty"`

	// WarningForeground AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-quicksight-theme-uicolorpalette.html#cfn-quicksight-theme-uicolorpalette-warningforeground
	WarningForeground string `json:"WarningForeground,omitempty"`

	// AWSCloudFormationDeletionPolicy represents a CloudFormation DeletionPolicy
	AWSCloudFormationDeletionPolicy policies.DeletionPolicy `json:"-"`

	// AWSCloudFormationUpdateReplacePolicy represents a CloudFormation UpdateReplacePolicy
	AWSCloudFormationUpdateReplacePolicy policies.UpdateReplacePolicy `json:"-"`

	// AWSCloudFormationDependsOn stores the logical ID of the resources to be created before this resource
	AWSCloudFormationDependsOn []string `json:"-"`

	// AWSCloudFormationMetadata stores structured data associated with this resource
	AWSCloudFormationMetadata map[string]interface{} `json:"-"`

	// AWSCloudFormationCondition stores the logical ID of the condition that must be satisfied for this resource to be created
	AWSCloudFormationCondition string `json:"-"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *Theme_UIColorPalette) AWSCloudFormationType() string {
	return "AWS::QuickSight::Theme.UIColorPalette"
}
