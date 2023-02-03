package mwaa

import (
	"github.com/awslabs/goformation/v4/cloudformation/policies"
)

// Environment_SecurityGroupList AWS CloudFormation Resource (AWS::MWAA::Environment.SecurityGroupList)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-mwaa-environment-securitygrouplist.html
type Environment_SecurityGroupList struct {

	// SecurityGroupList AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-mwaa-environment-securitygrouplist.html#cfn-mwaa-environment-securitygrouplist-securitygrouplist
	SecurityGroupList []string `json:"SecurityGroupList,omitempty"`

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
func (r *Environment_SecurityGroupList) AWSCloudFormationType() string {
	return "AWS::MWAA::Environment.SecurityGroupList"
}
