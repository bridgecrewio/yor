package quicksight

import (
	"github.com/awslabs/goformation/v4/cloudformation/policies"
)

// Template_DataSetSchema AWS CloudFormation Resource (AWS::QuickSight::Template.DataSetSchema)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-quicksight-template-datasetschema.html
type Template_DataSetSchema struct {

	// ColumnSchemaList AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-quicksight-template-datasetschema.html#cfn-quicksight-template-datasetschema-columnschemalist
	ColumnSchemaList []Template_ColumnSchema `json:"ColumnSchemaList,omitempty"`

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
func (r *Template_DataSetSchema) AWSCloudFormationType() string {
	return "AWS::QuickSight::Template.DataSetSchema"
}
