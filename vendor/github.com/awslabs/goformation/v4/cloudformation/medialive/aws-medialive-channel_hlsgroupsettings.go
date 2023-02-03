package medialive

import (
	"github.com/awslabs/goformation/v4/cloudformation/policies"
)

// Channel_HlsGroupSettings AWS CloudFormation Resource (AWS::MediaLive::Channel.HlsGroupSettings)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-hlsgroupsettings.html
type Channel_HlsGroupSettings struct {

	// AdMarkers AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-hlsgroupsettings.html#cfn-medialive-channel-hlsgroupsettings-admarkers
	AdMarkers []string `json:"AdMarkers,omitempty"`

	// BaseUrlContent AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-hlsgroupsettings.html#cfn-medialive-channel-hlsgroupsettings-baseurlcontent
	BaseUrlContent string `json:"BaseUrlContent,omitempty"`

	// BaseUrlContent1 AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-hlsgroupsettings.html#cfn-medialive-channel-hlsgroupsettings-baseurlcontent1
	BaseUrlContent1 string `json:"BaseUrlContent1,omitempty"`

	// BaseUrlManifest AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-hlsgroupsettings.html#cfn-medialive-channel-hlsgroupsettings-baseurlmanifest
	BaseUrlManifest string `json:"BaseUrlManifest,omitempty"`

	// BaseUrlManifest1 AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-hlsgroupsettings.html#cfn-medialive-channel-hlsgroupsettings-baseurlmanifest1
	BaseUrlManifest1 string `json:"BaseUrlManifest1,omitempty"`

	// CaptionLanguageMappings AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-hlsgroupsettings.html#cfn-medialive-channel-hlsgroupsettings-captionlanguagemappings
	CaptionLanguageMappings []Channel_CaptionLanguageMapping `json:"CaptionLanguageMappings,omitempty"`

	// CaptionLanguageSetting AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-hlsgroupsettings.html#cfn-medialive-channel-hlsgroupsettings-captionlanguagesetting
	CaptionLanguageSetting string `json:"CaptionLanguageSetting,omitempty"`

	// ClientCache AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-hlsgroupsettings.html#cfn-medialive-channel-hlsgroupsettings-clientcache
	ClientCache string `json:"ClientCache,omitempty"`

	// CodecSpecification AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-hlsgroupsettings.html#cfn-medialive-channel-hlsgroupsettings-codecspecification
	CodecSpecification string `json:"CodecSpecification,omitempty"`

	// ConstantIv AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-hlsgroupsettings.html#cfn-medialive-channel-hlsgroupsettings-constantiv
	ConstantIv string `json:"ConstantIv,omitempty"`

	// Destination AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-hlsgroupsettings.html#cfn-medialive-channel-hlsgroupsettings-destination
	Destination *Channel_OutputLocationRef `json:"Destination,omitempty"`

	// DirectoryStructure AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-hlsgroupsettings.html#cfn-medialive-channel-hlsgroupsettings-directorystructure
	DirectoryStructure string `json:"DirectoryStructure,omitempty"`

	// DiscontinuityTags AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-hlsgroupsettings.html#cfn-medialive-channel-hlsgroupsettings-discontinuitytags
	DiscontinuityTags string `json:"DiscontinuityTags,omitempty"`

	// EncryptionType AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-hlsgroupsettings.html#cfn-medialive-channel-hlsgroupsettings-encryptiontype
	EncryptionType string `json:"EncryptionType,omitempty"`

	// HlsCdnSettings AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-hlsgroupsettings.html#cfn-medialive-channel-hlsgroupsettings-hlscdnsettings
	HlsCdnSettings *Channel_HlsCdnSettings `json:"HlsCdnSettings,omitempty"`

	// HlsId3SegmentTagging AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-hlsgroupsettings.html#cfn-medialive-channel-hlsgroupsettings-hlsid3segmenttagging
	HlsId3SegmentTagging string `json:"HlsId3SegmentTagging,omitempty"`

	// IFrameOnlyPlaylists AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-hlsgroupsettings.html#cfn-medialive-channel-hlsgroupsettings-iframeonlyplaylists
	IFrameOnlyPlaylists string `json:"IFrameOnlyPlaylists,omitempty"`

	// IncompleteSegmentBehavior AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-hlsgroupsettings.html#cfn-medialive-channel-hlsgroupsettings-incompletesegmentbehavior
	IncompleteSegmentBehavior string `json:"IncompleteSegmentBehavior,omitempty"`

	// IndexNSegments AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-hlsgroupsettings.html#cfn-medialive-channel-hlsgroupsettings-indexnsegments
	IndexNSegments int `json:"IndexNSegments,omitempty"`

	// InputLossAction AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-hlsgroupsettings.html#cfn-medialive-channel-hlsgroupsettings-inputlossaction
	InputLossAction string `json:"InputLossAction,omitempty"`

	// IvInManifest AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-hlsgroupsettings.html#cfn-medialive-channel-hlsgroupsettings-ivinmanifest
	IvInManifest string `json:"IvInManifest,omitempty"`

	// IvSource AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-hlsgroupsettings.html#cfn-medialive-channel-hlsgroupsettings-ivsource
	IvSource string `json:"IvSource,omitempty"`

	// KeepSegments AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-hlsgroupsettings.html#cfn-medialive-channel-hlsgroupsettings-keepsegments
	KeepSegments int `json:"KeepSegments,omitempty"`

	// KeyFormat AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-hlsgroupsettings.html#cfn-medialive-channel-hlsgroupsettings-keyformat
	KeyFormat string `json:"KeyFormat,omitempty"`

	// KeyFormatVersions AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-hlsgroupsettings.html#cfn-medialive-channel-hlsgroupsettings-keyformatversions
	KeyFormatVersions string `json:"KeyFormatVersions,omitempty"`

	// KeyProviderSettings AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-hlsgroupsettings.html#cfn-medialive-channel-hlsgroupsettings-keyprovidersettings
	KeyProviderSettings *Channel_KeyProviderSettings `json:"KeyProviderSettings,omitempty"`

	// ManifestCompression AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-hlsgroupsettings.html#cfn-medialive-channel-hlsgroupsettings-manifestcompression
	ManifestCompression string `json:"ManifestCompression,omitempty"`

	// ManifestDurationFormat AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-hlsgroupsettings.html#cfn-medialive-channel-hlsgroupsettings-manifestdurationformat
	ManifestDurationFormat string `json:"ManifestDurationFormat,omitempty"`

	// MinSegmentLength AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-hlsgroupsettings.html#cfn-medialive-channel-hlsgroupsettings-minsegmentlength
	MinSegmentLength int `json:"MinSegmentLength,omitempty"`

	// Mode AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-hlsgroupsettings.html#cfn-medialive-channel-hlsgroupsettings-mode
	Mode string `json:"Mode,omitempty"`

	// OutputSelection AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-hlsgroupsettings.html#cfn-medialive-channel-hlsgroupsettings-outputselection
	OutputSelection string `json:"OutputSelection,omitempty"`

	// ProgramDateTime AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-hlsgroupsettings.html#cfn-medialive-channel-hlsgroupsettings-programdatetime
	ProgramDateTime string `json:"ProgramDateTime,omitempty"`

	// ProgramDateTimePeriod AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-hlsgroupsettings.html#cfn-medialive-channel-hlsgroupsettings-programdatetimeperiod
	ProgramDateTimePeriod int `json:"ProgramDateTimePeriod,omitempty"`

	// RedundantManifest AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-hlsgroupsettings.html#cfn-medialive-channel-hlsgroupsettings-redundantmanifest
	RedundantManifest string `json:"RedundantManifest,omitempty"`

	// SegmentLength AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-hlsgroupsettings.html#cfn-medialive-channel-hlsgroupsettings-segmentlength
	SegmentLength int `json:"SegmentLength,omitempty"`

	// SegmentationMode AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-hlsgroupsettings.html#cfn-medialive-channel-hlsgroupsettings-segmentationmode
	SegmentationMode string `json:"SegmentationMode,omitempty"`

	// SegmentsPerSubdirectory AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-hlsgroupsettings.html#cfn-medialive-channel-hlsgroupsettings-segmentspersubdirectory
	SegmentsPerSubdirectory int `json:"SegmentsPerSubdirectory,omitempty"`

	// StreamInfResolution AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-hlsgroupsettings.html#cfn-medialive-channel-hlsgroupsettings-streaminfresolution
	StreamInfResolution string `json:"StreamInfResolution,omitempty"`

	// TimedMetadataId3Frame AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-hlsgroupsettings.html#cfn-medialive-channel-hlsgroupsettings-timedmetadataid3frame
	TimedMetadataId3Frame string `json:"TimedMetadataId3Frame,omitempty"`

	// TimedMetadataId3Period AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-hlsgroupsettings.html#cfn-medialive-channel-hlsgroupsettings-timedmetadataid3period
	TimedMetadataId3Period int `json:"TimedMetadataId3Period,omitempty"`

	// TimestampDeltaMilliseconds AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-hlsgroupsettings.html#cfn-medialive-channel-hlsgroupsettings-timestampdeltamilliseconds
	TimestampDeltaMilliseconds int `json:"TimestampDeltaMilliseconds,omitempty"`

	// TsFileMode AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-hlsgroupsettings.html#cfn-medialive-channel-hlsgroupsettings-tsfilemode
	TsFileMode string `json:"TsFileMode,omitempty"`

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
func (r *Channel_HlsGroupSettings) AWSCloudFormationType() string {
	return "AWS::MediaLive::Channel.HlsGroupSettings"
}
