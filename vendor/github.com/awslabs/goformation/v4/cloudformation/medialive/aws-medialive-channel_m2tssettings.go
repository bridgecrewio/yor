package medialive

import (
	"github.com/awslabs/goformation/v4/cloudformation/policies"
)

// Channel_M2tsSettings AWS CloudFormation Resource (AWS::MediaLive::Channel.M2tsSettings)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-m2tssettings.html
type Channel_M2tsSettings struct {

	// AbsentInputAudioBehavior AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-m2tssettings.html#cfn-medialive-channel-m2tssettings-absentinputaudiobehavior
	AbsentInputAudioBehavior string `json:"AbsentInputAudioBehavior,omitempty"`

	// Arib AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-m2tssettings.html#cfn-medialive-channel-m2tssettings-arib
	Arib string `json:"Arib,omitempty"`

	// AribCaptionsPid AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-m2tssettings.html#cfn-medialive-channel-m2tssettings-aribcaptionspid
	AribCaptionsPid string `json:"AribCaptionsPid,omitempty"`

	// AribCaptionsPidControl AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-m2tssettings.html#cfn-medialive-channel-m2tssettings-aribcaptionspidcontrol
	AribCaptionsPidControl string `json:"AribCaptionsPidControl,omitempty"`

	// AudioBufferModel AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-m2tssettings.html#cfn-medialive-channel-m2tssettings-audiobuffermodel
	AudioBufferModel string `json:"AudioBufferModel,omitempty"`

	// AudioFramesPerPes AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-m2tssettings.html#cfn-medialive-channel-m2tssettings-audioframesperpes
	AudioFramesPerPes int `json:"AudioFramesPerPes,omitempty"`

	// AudioPids AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-m2tssettings.html#cfn-medialive-channel-m2tssettings-audiopids
	AudioPids string `json:"AudioPids,omitempty"`

	// AudioStreamType AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-m2tssettings.html#cfn-medialive-channel-m2tssettings-audiostreamtype
	AudioStreamType string `json:"AudioStreamType,omitempty"`

	// Bitrate AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-m2tssettings.html#cfn-medialive-channel-m2tssettings-bitrate
	Bitrate int `json:"Bitrate,omitempty"`

	// BufferModel AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-m2tssettings.html#cfn-medialive-channel-m2tssettings-buffermodel
	BufferModel string `json:"BufferModel,omitempty"`

	// CcDescriptor AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-m2tssettings.html#cfn-medialive-channel-m2tssettings-ccdescriptor
	CcDescriptor string `json:"CcDescriptor,omitempty"`

	// DvbNitSettings AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-m2tssettings.html#cfn-medialive-channel-m2tssettings-dvbnitsettings
	DvbNitSettings *Channel_DvbNitSettings `json:"DvbNitSettings,omitempty"`

	// DvbSdtSettings AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-m2tssettings.html#cfn-medialive-channel-m2tssettings-dvbsdtsettings
	DvbSdtSettings *Channel_DvbSdtSettings `json:"DvbSdtSettings,omitempty"`

	// DvbSubPids AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-m2tssettings.html#cfn-medialive-channel-m2tssettings-dvbsubpids
	DvbSubPids string `json:"DvbSubPids,omitempty"`

	// DvbTdtSettings AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-m2tssettings.html#cfn-medialive-channel-m2tssettings-dvbtdtsettings
	DvbTdtSettings *Channel_DvbTdtSettings `json:"DvbTdtSettings,omitempty"`

	// DvbTeletextPid AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-m2tssettings.html#cfn-medialive-channel-m2tssettings-dvbteletextpid
	DvbTeletextPid string `json:"DvbTeletextPid,omitempty"`

	// Ebif AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-m2tssettings.html#cfn-medialive-channel-m2tssettings-ebif
	Ebif string `json:"Ebif,omitempty"`

	// EbpAudioInterval AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-m2tssettings.html#cfn-medialive-channel-m2tssettings-ebpaudiointerval
	EbpAudioInterval string `json:"EbpAudioInterval,omitempty"`

	// EbpLookaheadMs AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-m2tssettings.html#cfn-medialive-channel-m2tssettings-ebplookaheadms
	EbpLookaheadMs int `json:"EbpLookaheadMs,omitempty"`

	// EbpPlacement AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-m2tssettings.html#cfn-medialive-channel-m2tssettings-ebpplacement
	EbpPlacement string `json:"EbpPlacement,omitempty"`

	// EcmPid AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-m2tssettings.html#cfn-medialive-channel-m2tssettings-ecmpid
	EcmPid string `json:"EcmPid,omitempty"`

	// EsRateInPes AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-m2tssettings.html#cfn-medialive-channel-m2tssettings-esrateinpes
	EsRateInPes string `json:"EsRateInPes,omitempty"`

	// EtvPlatformPid AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-m2tssettings.html#cfn-medialive-channel-m2tssettings-etvplatformpid
	EtvPlatformPid string `json:"EtvPlatformPid,omitempty"`

	// EtvSignalPid AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-m2tssettings.html#cfn-medialive-channel-m2tssettings-etvsignalpid
	EtvSignalPid string `json:"EtvSignalPid,omitempty"`

	// FragmentTime AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-m2tssettings.html#cfn-medialive-channel-m2tssettings-fragmenttime
	FragmentTime float64 `json:"FragmentTime,omitempty"`

	// Klv AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-m2tssettings.html#cfn-medialive-channel-m2tssettings-klv
	Klv string `json:"Klv,omitempty"`

	// KlvDataPids AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-m2tssettings.html#cfn-medialive-channel-m2tssettings-klvdatapids
	KlvDataPids string `json:"KlvDataPids,omitempty"`

	// NielsenId3Behavior AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-m2tssettings.html#cfn-medialive-channel-m2tssettings-nielsenid3behavior
	NielsenId3Behavior string `json:"NielsenId3Behavior,omitempty"`

	// NullPacketBitrate AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-m2tssettings.html#cfn-medialive-channel-m2tssettings-nullpacketbitrate
	NullPacketBitrate float64 `json:"NullPacketBitrate,omitempty"`

	// PatInterval AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-m2tssettings.html#cfn-medialive-channel-m2tssettings-patinterval
	PatInterval int `json:"PatInterval,omitempty"`

	// PcrControl AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-m2tssettings.html#cfn-medialive-channel-m2tssettings-pcrcontrol
	PcrControl string `json:"PcrControl,omitempty"`

	// PcrPeriod AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-m2tssettings.html#cfn-medialive-channel-m2tssettings-pcrperiod
	PcrPeriod int `json:"PcrPeriod,omitempty"`

	// PcrPid AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-m2tssettings.html#cfn-medialive-channel-m2tssettings-pcrpid
	PcrPid string `json:"PcrPid,omitempty"`

	// PmtInterval AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-m2tssettings.html#cfn-medialive-channel-m2tssettings-pmtinterval
	PmtInterval int `json:"PmtInterval,omitempty"`

	// PmtPid AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-m2tssettings.html#cfn-medialive-channel-m2tssettings-pmtpid
	PmtPid string `json:"PmtPid,omitempty"`

	// ProgramNum AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-m2tssettings.html#cfn-medialive-channel-m2tssettings-programnum
	ProgramNum int `json:"ProgramNum,omitempty"`

	// RateMode AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-m2tssettings.html#cfn-medialive-channel-m2tssettings-ratemode
	RateMode string `json:"RateMode,omitempty"`

	// Scte27Pids AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-m2tssettings.html#cfn-medialive-channel-m2tssettings-scte27pids
	Scte27Pids string `json:"Scte27Pids,omitempty"`

	// Scte35Control AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-m2tssettings.html#cfn-medialive-channel-m2tssettings-scte35control
	Scte35Control string `json:"Scte35Control,omitempty"`

	// Scte35Pid AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-m2tssettings.html#cfn-medialive-channel-m2tssettings-scte35pid
	Scte35Pid string `json:"Scte35Pid,omitempty"`

	// SegmentationMarkers AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-m2tssettings.html#cfn-medialive-channel-m2tssettings-segmentationmarkers
	SegmentationMarkers string `json:"SegmentationMarkers,omitempty"`

	// SegmentationStyle AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-m2tssettings.html#cfn-medialive-channel-m2tssettings-segmentationstyle
	SegmentationStyle string `json:"SegmentationStyle,omitempty"`

	// SegmentationTime AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-m2tssettings.html#cfn-medialive-channel-m2tssettings-segmentationtime
	SegmentationTime float64 `json:"SegmentationTime,omitempty"`

	// TimedMetadataBehavior AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-m2tssettings.html#cfn-medialive-channel-m2tssettings-timedmetadatabehavior
	TimedMetadataBehavior string `json:"TimedMetadataBehavior,omitempty"`

	// TimedMetadataPid AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-m2tssettings.html#cfn-medialive-channel-m2tssettings-timedmetadatapid
	TimedMetadataPid string `json:"TimedMetadataPid,omitempty"`

	// TransportStreamId AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-m2tssettings.html#cfn-medialive-channel-m2tssettings-transportstreamid
	TransportStreamId int `json:"TransportStreamId,omitempty"`

	// VideoPid AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-m2tssettings.html#cfn-medialive-channel-m2tssettings-videopid
	VideoPid string `json:"VideoPid,omitempty"`

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
func (r *Channel_M2tsSettings) AWSCloudFormationType() string {
	return "AWS::MediaLive::Channel.M2tsSettings"
}
