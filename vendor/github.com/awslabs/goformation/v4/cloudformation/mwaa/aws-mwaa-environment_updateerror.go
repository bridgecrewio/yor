package mwaa

import (
	"github.com/awslabs/goformation/v4/cloudformation/policies"
)

// Environment_UpdateError AWS CloudFormation Resource (AWS::MWAA::Environment.UpdateError)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-mwaa-environment-updateerror.html
type Environment_UpdateError struct {

	// ErrorCode AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-mwaa-environment-updateerror.html#cfn-mwaa-environment-updateerror-errorcode
	ErrorCode string `json:"ErrorCode,omitempty"`

	// ErrorMessage AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-mwaa-environment-updateerror.html#cfn-mwaa-environment-updateerror-errormessage
	ErrorMessage string `json:"ErrorMessage,omitempty"`

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
func (r *Environment_UpdateError) AWSCloudFormationType() string {
	return "AWS::MWAA::Environment.UpdateError"
}
