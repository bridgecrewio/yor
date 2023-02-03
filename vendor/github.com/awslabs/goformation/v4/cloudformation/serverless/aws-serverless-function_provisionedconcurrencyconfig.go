package serverless

import (
	"github.com/awslabs/goformation/v4/cloudformation/policies"
)

// Function_ProvisionedConcurrencyConfig AWS CloudFormation Resource (AWS::Serverless::Function.ProvisionedConcurrencyConfig)
// See: https://github.com/awslabs/serverless-application-model/blob/master/versions/2016-10-31.md#provisioned-concurrency-config-object
type Function_ProvisionedConcurrencyConfig struct {

	// ProvisionedConcurrentExecutions AWS CloudFormation Property
	// Required: true
	// See: https://github.com/awslabs/serverless-application-model/blob/master/versions/2016-10-31.md#provisioned-concurrency-config-object
	ProvisionedConcurrentExecutions string `json:"ProvisionedConcurrentExecutions,omitempty"`

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
func (r *Function_ProvisionedConcurrencyConfig) AWSCloudFormationType() string {
	return "AWS::Serverless::Function.ProvisionedConcurrencyConfig"
}
