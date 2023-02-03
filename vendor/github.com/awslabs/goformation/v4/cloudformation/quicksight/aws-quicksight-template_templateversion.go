package quicksight

import (
	"github.com/awslabs/goformation/v4/cloudformation/policies"
)

// Template_TemplateVersion AWS CloudFormation Resource (AWS::QuickSight::Template.TemplateVersion)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-quicksight-template-templateversion.html
type Template_TemplateVersion struct {

	// CreatedTime AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-quicksight-template-templateversion.html#cfn-quicksight-template-templateversion-createdtime
	CreatedTime string `json:"CreatedTime,omitempty"`

	// DataSetConfigurations AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-quicksight-template-templateversion.html#cfn-quicksight-template-templateversion-datasetconfigurations
	DataSetConfigurations []Template_DataSetConfiguration `json:"DataSetConfigurations,omitempty"`

	// Description AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-quicksight-template-templateversion.html#cfn-quicksight-template-templateversion-description
	Description string `json:"Description,omitempty"`

	// Errors AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-quicksight-template-templateversion.html#cfn-quicksight-template-templateversion-errors
	Errors []Template_TemplateError `json:"Errors,omitempty"`

	// Sheets AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-quicksight-template-templateversion.html#cfn-quicksight-template-templateversion-sheets
	Sheets []Template_Sheet `json:"Sheets,omitempty"`

	// SourceEntityArn AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-quicksight-template-templateversion.html#cfn-quicksight-template-templateversion-sourceentityarn
	SourceEntityArn string `json:"SourceEntityArn,omitempty"`

	// Status AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-quicksight-template-templateversion.html#cfn-quicksight-template-templateversion-status
	Status string `json:"Status,omitempty"`

	// ThemeArn AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-quicksight-template-templateversion.html#cfn-quicksight-template-templateversion-themearn
	ThemeArn string `json:"ThemeArn,omitempty"`

	// VersionNumber AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-quicksight-template-templateversion.html#cfn-quicksight-template-templateversion-versionnumber
	VersionNumber float64 `json:"VersionNumber,omitempty"`

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
func (r *Template_TemplateVersion) AWSCloudFormationType() string {
	return "AWS::QuickSight::Template.TemplateVersion"
}
