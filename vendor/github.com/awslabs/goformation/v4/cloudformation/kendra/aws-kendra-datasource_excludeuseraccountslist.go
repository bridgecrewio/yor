package kendra

import (
	"github.com/awslabs/goformation/v4/cloudformation/policies"
)

// DataSource_ExcludeUserAccountsList AWS CloudFormation Resource (AWS::Kendra::DataSource.ExcludeUserAccountsList)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kendra-datasource-excludeuseraccountslist.html
type DataSource_ExcludeUserAccountsList struct {

	// ExcludeUserAccountsList AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kendra-datasource-excludeuseraccountslist.html#cfn-kendra-datasource-excludeuseraccountslist-excludeuseraccountslist
	ExcludeUserAccountsList []string `json:"ExcludeUserAccountsList,omitempty"`

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
func (r *DataSource_ExcludeUserAccountsList) AWSCloudFormationType() string {
	return "AWS::Kendra::DataSource.ExcludeUserAccountsList"
}
