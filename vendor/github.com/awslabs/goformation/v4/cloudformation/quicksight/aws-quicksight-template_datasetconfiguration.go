package quicksight

import (
	"github.com/awslabs/goformation/v4/cloudformation/policies"
)

// Template_DataSetConfiguration AWS CloudFormation Resource (AWS::QuickSight::Template.DataSetConfiguration)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-quicksight-template-datasetconfiguration.html
type Template_DataSetConfiguration struct {

	// ColumnGroupSchemaList AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-quicksight-template-datasetconfiguration.html#cfn-quicksight-template-datasetconfiguration-columngroupschemalist
	ColumnGroupSchemaList []Template_ColumnGroupSchema `json:"ColumnGroupSchemaList,omitempty"`

	// DataSetSchema AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-quicksight-template-datasetconfiguration.html#cfn-quicksight-template-datasetconfiguration-datasetschema
	DataSetSchema *Template_DataSetSchema `json:"DataSetSchema,omitempty"`

	// Placeholder AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-quicksight-template-datasetconfiguration.html#cfn-quicksight-template-datasetconfiguration-placeholder
	Placeholder string `json:"Placeholder,omitempty"`

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
func (r *Template_DataSetConfiguration) AWSCloudFormationType() string {
	return "AWS::QuickSight::Template.DataSetConfiguration"
}
