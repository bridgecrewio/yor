package serverless

import (
	"github.com/awslabs/goformation/v5/cloudformation/policies"
)

// Function_EventInvokeConfig AWS CloudFormation Resource (AWS::Serverless::Function.EventInvokeConfig)
// See: https://github.com/aws/serverless-application-model/blob/master/versions/2016-10-31.md#event-invoke-config-object
type Function_EventInvokeConfig struct {

	// DestinationConfig AWS CloudFormation Property
	// Required: false
	// See: https://github.com/aws/serverless-application-model/blob/master/versions/2016-10-31.md#event-invoke-config-object
	DestinationConfig *Function_DestinationConfig `json:"DestinationConfig,omitempty"`

	// MaximumEventAgeInSeconds AWS CloudFormation Property
	// Required: false
	// See: https://github.com/aws/serverless-application-model/blob/master/versions/2016-10-31.md#event-invoke-config-object
	MaximumEventAgeInSeconds int `json:"MaximumEventAgeInSeconds,omitempty"`

	// MaximumRetryAttempts AWS CloudFormation Property
	// Required: false
	// See: https://github.com/aws/serverless-application-model/blob/master/versions/2016-10-31.md#event-invoke-config-object
	MaximumRetryAttempts int `json:"MaximumRetryAttempts,omitempty"`

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
func (r *Function_EventInvokeConfig) AWSCloudFormationType() string {
	return "AWS::Serverless::Function.EventInvokeConfig"
}
