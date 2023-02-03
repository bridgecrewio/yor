package medialive

import (
	"github.com/awslabs/goformation/v4/cloudformation/policies"
)

// Channel_Eac3Settings AWS CloudFormation Resource (AWS::MediaLive::Channel.Eac3Settings)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-eac3settings.html
type Channel_Eac3Settings struct {

	// AttenuationControl AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-eac3settings.html#cfn-medialive-channel-eac3settings-attenuationcontrol
	AttenuationControl string `json:"AttenuationControl,omitempty"`

	// Bitrate AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-eac3settings.html#cfn-medialive-channel-eac3settings-bitrate
	Bitrate float64 `json:"Bitrate,omitempty"`

	// BitstreamMode AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-eac3settings.html#cfn-medialive-channel-eac3settings-bitstreammode
	BitstreamMode string `json:"BitstreamMode,omitempty"`

	// CodingMode AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-eac3settings.html#cfn-medialive-channel-eac3settings-codingmode
	CodingMode string `json:"CodingMode,omitempty"`

	// DcFilter AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-eac3settings.html#cfn-medialive-channel-eac3settings-dcfilter
	DcFilter string `json:"DcFilter,omitempty"`

	// Dialnorm AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-eac3settings.html#cfn-medialive-channel-eac3settings-dialnorm
	Dialnorm int `json:"Dialnorm,omitempty"`

	// DrcLine AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-eac3settings.html#cfn-medialive-channel-eac3settings-drcline
	DrcLine string `json:"DrcLine,omitempty"`

	// DrcRf AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-eac3settings.html#cfn-medialive-channel-eac3settings-drcrf
	DrcRf string `json:"DrcRf,omitempty"`

	// LfeControl AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-eac3settings.html#cfn-medialive-channel-eac3settings-lfecontrol
	LfeControl string `json:"LfeControl,omitempty"`

	// LfeFilter AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-eac3settings.html#cfn-medialive-channel-eac3settings-lfefilter
	LfeFilter string `json:"LfeFilter,omitempty"`

	// LoRoCenterMixLevel AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-eac3settings.html#cfn-medialive-channel-eac3settings-lorocentermixlevel
	LoRoCenterMixLevel float64 `json:"LoRoCenterMixLevel,omitempty"`

	// LoRoSurroundMixLevel AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-eac3settings.html#cfn-medialive-channel-eac3settings-lorosurroundmixlevel
	LoRoSurroundMixLevel float64 `json:"LoRoSurroundMixLevel,omitempty"`

	// LtRtCenterMixLevel AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-eac3settings.html#cfn-medialive-channel-eac3settings-ltrtcentermixlevel
	LtRtCenterMixLevel float64 `json:"LtRtCenterMixLevel,omitempty"`

	// LtRtSurroundMixLevel AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-eac3settings.html#cfn-medialive-channel-eac3settings-ltrtsurroundmixlevel
	LtRtSurroundMixLevel float64 `json:"LtRtSurroundMixLevel,omitempty"`

	// MetadataControl AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-eac3settings.html#cfn-medialive-channel-eac3settings-metadatacontrol
	MetadataControl string `json:"MetadataControl,omitempty"`

	// PassthroughControl AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-eac3settings.html#cfn-medialive-channel-eac3settings-passthroughcontrol
	PassthroughControl string `json:"PassthroughControl,omitempty"`

	// PhaseControl AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-eac3settings.html#cfn-medialive-channel-eac3settings-phasecontrol
	PhaseControl string `json:"PhaseControl,omitempty"`

	// StereoDownmix AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-eac3settings.html#cfn-medialive-channel-eac3settings-stereodownmix
	StereoDownmix string `json:"StereoDownmix,omitempty"`

	// SurroundExMode AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-eac3settings.html#cfn-medialive-channel-eac3settings-surroundexmode
	SurroundExMode string `json:"SurroundExMode,omitempty"`

	// SurroundMode AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-eac3settings.html#cfn-medialive-channel-eac3settings-surroundmode
	SurroundMode string `json:"SurroundMode,omitempty"`

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
func (r *Channel_Eac3Settings) AWSCloudFormationType() string {
	return "AWS::MediaLive::Channel.Eac3Settings"
}
