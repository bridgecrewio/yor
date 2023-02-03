package medialive

import (
	"github.com/awslabs/goformation/v4/cloudformation/policies"
)

// Channel_BurnInDestinationSettings AWS CloudFormation Resource (AWS::MediaLive::Channel.BurnInDestinationSettings)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-burnindestinationsettings.html
type Channel_BurnInDestinationSettings struct {

	// Alignment AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-burnindestinationsettings.html#cfn-medialive-channel-burnindestinationsettings-alignment
	Alignment string `json:"Alignment,omitempty"`

	// BackgroundColor AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-burnindestinationsettings.html#cfn-medialive-channel-burnindestinationsettings-backgroundcolor
	BackgroundColor string `json:"BackgroundColor,omitempty"`

	// BackgroundOpacity AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-burnindestinationsettings.html#cfn-medialive-channel-burnindestinationsettings-backgroundopacity
	BackgroundOpacity int `json:"BackgroundOpacity,omitempty"`

	// Font AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-burnindestinationsettings.html#cfn-medialive-channel-burnindestinationsettings-font
	Font *Channel_InputLocation `json:"Font,omitempty"`

	// FontColor AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-burnindestinationsettings.html#cfn-medialive-channel-burnindestinationsettings-fontcolor
	FontColor string `json:"FontColor,omitempty"`

	// FontOpacity AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-burnindestinationsettings.html#cfn-medialive-channel-burnindestinationsettings-fontopacity
	FontOpacity int `json:"FontOpacity,omitempty"`

	// FontResolution AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-burnindestinationsettings.html#cfn-medialive-channel-burnindestinationsettings-fontresolution
	FontResolution int `json:"FontResolution,omitempty"`

	// FontSize AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-burnindestinationsettings.html#cfn-medialive-channel-burnindestinationsettings-fontsize
	FontSize string `json:"FontSize,omitempty"`

	// OutlineColor AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-burnindestinationsettings.html#cfn-medialive-channel-burnindestinationsettings-outlinecolor
	OutlineColor string `json:"OutlineColor,omitempty"`

	// OutlineSize AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-burnindestinationsettings.html#cfn-medialive-channel-burnindestinationsettings-outlinesize
	OutlineSize int `json:"OutlineSize,omitempty"`

	// ShadowColor AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-burnindestinationsettings.html#cfn-medialive-channel-burnindestinationsettings-shadowcolor
	ShadowColor string `json:"ShadowColor,omitempty"`

	// ShadowOpacity AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-burnindestinationsettings.html#cfn-medialive-channel-burnindestinationsettings-shadowopacity
	ShadowOpacity int `json:"ShadowOpacity,omitempty"`

	// ShadowXOffset AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-burnindestinationsettings.html#cfn-medialive-channel-burnindestinationsettings-shadowxoffset
	ShadowXOffset int `json:"ShadowXOffset,omitempty"`

	// ShadowYOffset AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-burnindestinationsettings.html#cfn-medialive-channel-burnindestinationsettings-shadowyoffset
	ShadowYOffset int `json:"ShadowYOffset,omitempty"`

	// TeletextGridControl AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-burnindestinationsettings.html#cfn-medialive-channel-burnindestinationsettings-teletextgridcontrol
	TeletextGridControl string `json:"TeletextGridControl,omitempty"`

	// XPosition AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-burnindestinationsettings.html#cfn-medialive-channel-burnindestinationsettings-xposition
	XPosition int `json:"XPosition,omitempty"`

	// YPosition AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-burnindestinationsettings.html#cfn-medialive-channel-burnindestinationsettings-yposition
	YPosition int `json:"YPosition,omitempty"`

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
func (r *Channel_BurnInDestinationSettings) AWSCloudFormationType() string {
	return "AWS::MediaLive::Channel.BurnInDestinationSettings"
}
