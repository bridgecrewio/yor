package greengrassv2

import (
	"github.com/awslabs/goformation/v4/cloudformation/policies"
)

// ComponentVersion_LambdaExecutionParameters AWS CloudFormation Resource (AWS::GreengrassV2::ComponentVersion.LambdaExecutionParameters)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-greengrassv2-componentversion-lambdaexecutionparameters.html
type ComponentVersion_LambdaExecutionParameters struct {

	// EnvironmentVariables AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-greengrassv2-componentversion-lambdaexecutionparameters.html#cfn-greengrassv2-componentversion-lambdaexecutionparameters-environmentvariables
	EnvironmentVariables map[string]string `json:"EnvironmentVariables,omitempty"`

	// EventSources AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-greengrassv2-componentversion-lambdaexecutionparameters.html#cfn-greengrassv2-componentversion-lambdaexecutionparameters-eventsources
	EventSources []ComponentVersion_LambdaEventSource `json:"EventSources,omitempty"`

	// ExecArgs AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-greengrassv2-componentversion-lambdaexecutionparameters.html#cfn-greengrassv2-componentversion-lambdaexecutionparameters-execargs
	ExecArgs []string `json:"ExecArgs,omitempty"`

	// InputPayloadEncodingType AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-greengrassv2-componentversion-lambdaexecutionparameters.html#cfn-greengrassv2-componentversion-lambdaexecutionparameters-inputpayloadencodingtype
	InputPayloadEncodingType string `json:"InputPayloadEncodingType,omitempty"`

	// LinuxProcessParams AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-greengrassv2-componentversion-lambdaexecutionparameters.html#cfn-greengrassv2-componentversion-lambdaexecutionparameters-linuxprocessparams
	LinuxProcessParams *ComponentVersion_LambdaLinuxProcessParams `json:"LinuxProcessParams,omitempty"`

	// MaxIdleTimeInSeconds AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-greengrassv2-componentversion-lambdaexecutionparameters.html#cfn-greengrassv2-componentversion-lambdaexecutionparameters-maxidletimeinseconds
	MaxIdleTimeInSeconds int `json:"MaxIdleTimeInSeconds,omitempty"`

	// MaxInstancesCount AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-greengrassv2-componentversion-lambdaexecutionparameters.html#cfn-greengrassv2-componentversion-lambdaexecutionparameters-maxinstancescount
	MaxInstancesCount int `json:"MaxInstancesCount,omitempty"`

	// MaxQueueSize AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-greengrassv2-componentversion-lambdaexecutionparameters.html#cfn-greengrassv2-componentversion-lambdaexecutionparameters-maxqueuesize
	MaxQueueSize int `json:"MaxQueueSize,omitempty"`

	// Pinned AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-greengrassv2-componentversion-lambdaexecutionparameters.html#cfn-greengrassv2-componentversion-lambdaexecutionparameters-pinned
	Pinned bool `json:"Pinned,omitempty"`

	// StatusTimeoutInSeconds AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-greengrassv2-componentversion-lambdaexecutionparameters.html#cfn-greengrassv2-componentversion-lambdaexecutionparameters-statustimeoutinseconds
	StatusTimeoutInSeconds int `json:"StatusTimeoutInSeconds,omitempty"`

	// TimeoutInSeconds AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-greengrassv2-componentversion-lambdaexecutionparameters.html#cfn-greengrassv2-componentversion-lambdaexecutionparameters-timeoutinseconds
	TimeoutInSeconds int `json:"TimeoutInSeconds,omitempty"`

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
func (r *ComponentVersion_LambdaExecutionParameters) AWSCloudFormationType() string {
	return "AWS::GreengrassV2::ComponentVersion.LambdaExecutionParameters"
}
