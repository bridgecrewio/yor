package medialive

import (
	"github.com/awslabs/goformation/v4/cloudformation/policies"
)

// Channel_EbuTtDDestinationSettings AWS CloudFormation Resource (AWS::MediaLive::Channel.EbuTtDDestinationSettings)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-ebuttddestinationsettings.html
type Channel_EbuTtDDestinationSettings struct {

	// CopyrightHolder AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-ebuttddestinationsettings.html#cfn-medialive-channel-ebuttddestinationsettings-copyrightholder
	CopyrightHolder string `json:"CopyrightHolder,omitempty"`

	// FillLineGap AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-ebuttddestinationsettings.html#cfn-medialive-channel-ebuttddestinationsettings-filllinegap
	FillLineGap string `json:"FillLineGap,omitempty"`

	// FontFamily AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-ebuttddestinationsettings.html#cfn-medialive-channel-ebuttddestinationsettings-fontfamily
	FontFamily string `json:"FontFamily,omitempty"`

	// StyleControl AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-ebuttddestinationsettings.html#cfn-medialive-channel-ebuttddestinationsettings-stylecontrol
	StyleControl string `json:"StyleControl,omitempty"`

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
func (r *Channel_EbuTtDDestinationSettings) AWSCloudFormationType() string {
	return "AWS::MediaLive::Channel.EbuTtDDestinationSettings"
}
