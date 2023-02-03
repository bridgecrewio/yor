package serverless

import (
	"github.com/awslabs/goformation/v5/cloudformation/policies"
)

// Function_ImageConfig AWS CloudFormation Resource (AWS::Serverless::Function.ImageConfig)
// See: https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-lambda-function-imageconfig.html
type Function_ImageConfig struct {

	// Command AWS CloudFormation Property
	// Required: false
	// See: https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-lambda-function-imageconfig.html#cfn-lambda-function-imageconfig-command
	Command []string `json:"Command,omitempty"`

	// EntryPoint AWS CloudFormation Property
	// Required: false
	// See: https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-lambda-function-imageconfig.html#cfn-lambda-function-imageconfig-entrypoint
	EntryPoint []string `json:"EntryPoint,omitempty"`

	// WorkingDirectory AWS CloudFormation Property
	// Required: false
	// See: https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-lambda-function-imageconfig.html#cfn-lambda-function-imageconfig-workingdirectory
	WorkingDirectory string `json:"WorkingDirectory,omitempty"`

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
func (r *Function_ImageConfig) AWSCloudFormationType() string {
	return "AWS::Serverless::Function.ImageConfig"
}
