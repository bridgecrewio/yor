package mediapackage

import (
	"github.com/awslabs/goformation/v4/cloudformation/policies"
)

// OriginEndpoint_HlsPackage AWS CloudFormation Resource (AWS::MediaPackage::OriginEndpoint.HlsPackage)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-mediapackage-originendpoint-hlspackage.html
type OriginEndpoint_HlsPackage struct {

	// AdMarkers AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-mediapackage-originendpoint-hlspackage.html#cfn-mediapackage-originendpoint-hlspackage-admarkers
	AdMarkers string `json:"AdMarkers,omitempty"`

	// AdTriggers AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-mediapackage-originendpoint-hlspackage.html#cfn-mediapackage-originendpoint-hlspackage-adtriggers
	AdTriggers []string `json:"AdTriggers,omitempty"`

	// AdsOnDeliveryRestrictions AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-mediapackage-originendpoint-hlspackage.html#cfn-mediapackage-originendpoint-hlspackage-adsondeliveryrestrictions
	AdsOnDeliveryRestrictions string `json:"AdsOnDeliveryRestrictions,omitempty"`

	// Encryption AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-mediapackage-originendpoint-hlspackage.html#cfn-mediapackage-originendpoint-hlspackage-encryption
	Encryption *OriginEndpoint_HlsEncryption `json:"Encryption,omitempty"`

	// IncludeIframeOnlyStream AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-mediapackage-originendpoint-hlspackage.html#cfn-mediapackage-originendpoint-hlspackage-includeiframeonlystream
	IncludeIframeOnlyStream bool `json:"IncludeIframeOnlyStream,omitempty"`

	// PlaylistType AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-mediapackage-originendpoint-hlspackage.html#cfn-mediapackage-originendpoint-hlspackage-playlisttype
	PlaylistType string `json:"PlaylistType,omitempty"`

	// PlaylistWindowSeconds AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-mediapackage-originendpoint-hlspackage.html#cfn-mediapackage-originendpoint-hlspackage-playlistwindowseconds
	PlaylistWindowSeconds int `json:"PlaylistWindowSeconds,omitempty"`

	// ProgramDateTimeIntervalSeconds AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-mediapackage-originendpoint-hlspackage.html#cfn-mediapackage-originendpoint-hlspackage-programdatetimeintervalseconds
	ProgramDateTimeIntervalSeconds int `json:"ProgramDateTimeIntervalSeconds,omitempty"`

	// SegmentDurationSeconds AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-mediapackage-originendpoint-hlspackage.html#cfn-mediapackage-originendpoint-hlspackage-segmentdurationseconds
	SegmentDurationSeconds int `json:"SegmentDurationSeconds,omitempty"`

	// StreamSelection AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-mediapackage-originendpoint-hlspackage.html#cfn-mediapackage-originendpoint-hlspackage-streamselection
	StreamSelection *OriginEndpoint_StreamSelection `json:"StreamSelection,omitempty"`

	// UseAudioRenditionGroup AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-mediapackage-originendpoint-hlspackage.html#cfn-mediapackage-originendpoint-hlspackage-useaudiorenditiongroup
	UseAudioRenditionGroup bool `json:"UseAudioRenditionGroup,omitempty"`

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
func (r *OriginEndpoint_HlsPackage) AWSCloudFormationType() string {
	return "AWS::MediaPackage::OriginEndpoint.HlsPackage"
}
