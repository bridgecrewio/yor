package mediapackage

import (
	"github.com/awslabs/goformation/v4/cloudformation/policies"
)

// OriginEndpoint_HlsManifest AWS CloudFormation Resource (AWS::MediaPackage::OriginEndpoint.HlsManifest)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-mediapackage-originendpoint-hlsmanifest.html
type OriginEndpoint_HlsManifest struct {

	// AdMarkers AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-mediapackage-originendpoint-hlsmanifest.html#cfn-mediapackage-originendpoint-hlsmanifest-admarkers
	AdMarkers string `json:"AdMarkers,omitempty"`

	// AdTriggers AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-mediapackage-originendpoint-hlsmanifest.html#cfn-mediapackage-originendpoint-hlsmanifest-adtriggers
	AdTriggers []string `json:"AdTriggers,omitempty"`

	// AdsOnDeliveryRestrictions AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-mediapackage-originendpoint-hlsmanifest.html#cfn-mediapackage-originendpoint-hlsmanifest-adsondeliveryrestrictions
	AdsOnDeliveryRestrictions string `json:"AdsOnDeliveryRestrictions,omitempty"`

	// Id AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-mediapackage-originendpoint-hlsmanifest.html#cfn-mediapackage-originendpoint-hlsmanifest-id
	Id string `json:"Id,omitempty"`

	// IncludeIframeOnlyStream AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-mediapackage-originendpoint-hlsmanifest.html#cfn-mediapackage-originendpoint-hlsmanifest-includeiframeonlystream
	IncludeIframeOnlyStream bool `json:"IncludeIframeOnlyStream,omitempty"`

	// ManifestName AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-mediapackage-originendpoint-hlsmanifest.html#cfn-mediapackage-originendpoint-hlsmanifest-manifestname
	ManifestName string `json:"ManifestName,omitempty"`

	// PlaylistType AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-mediapackage-originendpoint-hlsmanifest.html#cfn-mediapackage-originendpoint-hlsmanifest-playlisttype
	PlaylistType string `json:"PlaylistType,omitempty"`

	// PlaylistWindowSeconds AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-mediapackage-originendpoint-hlsmanifest.html#cfn-mediapackage-originendpoint-hlsmanifest-playlistwindowseconds
	PlaylistWindowSeconds int `json:"PlaylistWindowSeconds,omitempty"`

	// ProgramDateTimeIntervalSeconds AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-mediapackage-originendpoint-hlsmanifest.html#cfn-mediapackage-originendpoint-hlsmanifest-programdatetimeintervalseconds
	ProgramDateTimeIntervalSeconds int `json:"ProgramDateTimeIntervalSeconds,omitempty"`

	// Url AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-mediapackage-originendpoint-hlsmanifest.html#cfn-mediapackage-originendpoint-hlsmanifest-url
	Url string `json:"Url,omitempty"`

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
func (r *OriginEndpoint_HlsManifest) AWSCloudFormationType() string {
	return "AWS::MediaPackage::OriginEndpoint.HlsManifest"
}
