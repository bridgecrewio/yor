package mwaa

import (
	"github.com/awslabs/goformation/v4/cloudformation/policies"
)

// Environment_LastUpdate AWS CloudFormation Resource (AWS::MWAA::Environment.LastUpdate)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-mwaa-environment-lastupdate.html
type Environment_LastUpdate struct {

	// CreatedAt AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-mwaa-environment-lastupdate.html#cfn-mwaa-environment-lastupdate-createdat
	CreatedAt string `json:"CreatedAt,omitempty"`

	// Error AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-mwaa-environment-lastupdate.html#cfn-mwaa-environment-lastupdate-error
	Error *Environment_UpdateError `json:"Error,omitempty"`

	// Status AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-mwaa-environment-lastupdate.html#cfn-mwaa-environment-lastupdate-status
	Status string `json:"Status,omitempty"`

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
func (r *Environment_LastUpdate) AWSCloudFormationType() string {
	return "AWS::MWAA::Environment.LastUpdate"
}
