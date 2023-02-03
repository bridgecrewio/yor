package iotwireless

import (
	"github.com/awslabs/goformation/v4/cloudformation/policies"
)

// ServiceProfile_LoRaWANGetServiceProfileInfo AWS CloudFormation Resource (AWS::IoTWireless::ServiceProfile.LoRaWANGetServiceProfileInfo)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotwireless-serviceprofile-lorawangetserviceprofileinfo.html
type ServiceProfile_LoRaWANGetServiceProfileInfo struct {

	// AddGwMetadata AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotwireless-serviceprofile-lorawangetserviceprofileinfo.html#cfn-iotwireless-serviceprofile-lorawangetserviceprofileinfo-addgwmetadata
	AddGwMetadata bool `json:"AddGwMetadata,omitempty"`

	// ChannelMask AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotwireless-serviceprofile-lorawangetserviceprofileinfo.html#cfn-iotwireless-serviceprofile-lorawangetserviceprofileinfo-channelmask
	ChannelMask string `json:"ChannelMask,omitempty"`

	// DevStatusReqFreq AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotwireless-serviceprofile-lorawangetserviceprofileinfo.html#cfn-iotwireless-serviceprofile-lorawangetserviceprofileinfo-devstatusreqfreq
	DevStatusReqFreq int `json:"DevStatusReqFreq,omitempty"`

	// DlBucketSize AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotwireless-serviceprofile-lorawangetserviceprofileinfo.html#cfn-iotwireless-serviceprofile-lorawangetserviceprofileinfo-dlbucketsize
	DlBucketSize int `json:"DlBucketSize,omitempty"`

	// DlRate AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotwireless-serviceprofile-lorawangetserviceprofileinfo.html#cfn-iotwireless-serviceprofile-lorawangetserviceprofileinfo-dlrate
	DlRate int `json:"DlRate,omitempty"`

	// DlRatePolicy AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotwireless-serviceprofile-lorawangetserviceprofileinfo.html#cfn-iotwireless-serviceprofile-lorawangetserviceprofileinfo-dlratepolicy
	DlRatePolicy string `json:"DlRatePolicy,omitempty"`

	// DrMax AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotwireless-serviceprofile-lorawangetserviceprofileinfo.html#cfn-iotwireless-serviceprofile-lorawangetserviceprofileinfo-drmax
	DrMax int `json:"DrMax,omitempty"`

	// DrMin AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotwireless-serviceprofile-lorawangetserviceprofileinfo.html#cfn-iotwireless-serviceprofile-lorawangetserviceprofileinfo-drmin
	DrMin int `json:"DrMin,omitempty"`

	// HrAllowed AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotwireless-serviceprofile-lorawangetserviceprofileinfo.html#cfn-iotwireless-serviceprofile-lorawangetserviceprofileinfo-hrallowed
	HrAllowed bool `json:"HrAllowed,omitempty"`

	// MinGwDiversity AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotwireless-serviceprofile-lorawangetserviceprofileinfo.html#cfn-iotwireless-serviceprofile-lorawangetserviceprofileinfo-mingwdiversity
	MinGwDiversity int `json:"MinGwDiversity,omitempty"`

	// NwkGeoLoc AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotwireless-serviceprofile-lorawangetserviceprofileinfo.html#cfn-iotwireless-serviceprofile-lorawangetserviceprofileinfo-nwkgeoloc
	NwkGeoLoc bool `json:"NwkGeoLoc,omitempty"`

	// PrAllowed AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotwireless-serviceprofile-lorawangetserviceprofileinfo.html#cfn-iotwireless-serviceprofile-lorawangetserviceprofileinfo-prallowed
	PrAllowed bool `json:"PrAllowed,omitempty"`

	// RaAllowed AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotwireless-serviceprofile-lorawangetserviceprofileinfo.html#cfn-iotwireless-serviceprofile-lorawangetserviceprofileinfo-raallowed
	RaAllowed bool `json:"RaAllowed,omitempty"`

	// ReportDevStatusBattery AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotwireless-serviceprofile-lorawangetserviceprofileinfo.html#cfn-iotwireless-serviceprofile-lorawangetserviceprofileinfo-reportdevstatusbattery
	ReportDevStatusBattery bool `json:"ReportDevStatusBattery,omitempty"`

	// ReportDevStatusMargin AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotwireless-serviceprofile-lorawangetserviceprofileinfo.html#cfn-iotwireless-serviceprofile-lorawangetserviceprofileinfo-reportdevstatusmargin
	ReportDevStatusMargin bool `json:"ReportDevStatusMargin,omitempty"`

	// TargetPer AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotwireless-serviceprofile-lorawangetserviceprofileinfo.html#cfn-iotwireless-serviceprofile-lorawangetserviceprofileinfo-targetper
	TargetPer int `json:"TargetPer,omitempty"`

	// UlBucketSize AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotwireless-serviceprofile-lorawangetserviceprofileinfo.html#cfn-iotwireless-serviceprofile-lorawangetserviceprofileinfo-ulbucketsize
	UlBucketSize int `json:"UlBucketSize,omitempty"`

	// UlRate AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotwireless-serviceprofile-lorawangetserviceprofileinfo.html#cfn-iotwireless-serviceprofile-lorawangetserviceprofileinfo-ulrate
	UlRate int `json:"UlRate,omitempty"`

	// UlRatePolicy AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotwireless-serviceprofile-lorawangetserviceprofileinfo.html#cfn-iotwireless-serviceprofile-lorawangetserviceprofileinfo-ulratepolicy
	UlRatePolicy string `json:"UlRatePolicy,omitempty"`

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
func (r *ServiceProfile_LoRaWANGetServiceProfileInfo) AWSCloudFormationType() string {
	return "AWS::IoTWireless::ServiceProfile.LoRaWANGetServiceProfileInfo"
}
