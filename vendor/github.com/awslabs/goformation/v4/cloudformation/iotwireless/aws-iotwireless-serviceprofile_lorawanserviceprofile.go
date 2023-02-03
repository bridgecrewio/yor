package iotwireless

import (
	"github.com/awslabs/goformation/v4/cloudformation/policies"
)

// ServiceProfile_LoRaWANServiceProfile AWS CloudFormation Resource (AWS::IoTWireless::ServiceProfile.LoRaWANServiceProfile)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotwireless-serviceprofile-lorawanserviceprofile.html
type ServiceProfile_LoRaWANServiceProfile struct {

	// AddGwMetadata AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotwireless-serviceprofile-lorawanserviceprofile.html#cfn-iotwireless-serviceprofile-lorawanserviceprofile-addgwmetadata
	AddGwMetadata bool `json:"AddGwMetadata,omitempty"`

	// ChannelMask AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotwireless-serviceprofile-lorawanserviceprofile.html#cfn-iotwireless-serviceprofile-lorawanserviceprofile-channelmask
	ChannelMask string `json:"ChannelMask,omitempty"`

	// DevStatusReqFreq AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotwireless-serviceprofile-lorawanserviceprofile.html#cfn-iotwireless-serviceprofile-lorawanserviceprofile-devstatusreqfreq
	DevStatusReqFreq int `json:"DevStatusReqFreq,omitempty"`

	// DlBucketSize AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotwireless-serviceprofile-lorawanserviceprofile.html#cfn-iotwireless-serviceprofile-lorawanserviceprofile-dlbucketsize
	DlBucketSize int `json:"DlBucketSize,omitempty"`

	// DlRate AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotwireless-serviceprofile-lorawanserviceprofile.html#cfn-iotwireless-serviceprofile-lorawanserviceprofile-dlrate
	DlRate int `json:"DlRate,omitempty"`

	// DlRatePolicy AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotwireless-serviceprofile-lorawanserviceprofile.html#cfn-iotwireless-serviceprofile-lorawanserviceprofile-dlratepolicy
	DlRatePolicy string `json:"DlRatePolicy,omitempty"`

	// DrMax AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotwireless-serviceprofile-lorawanserviceprofile.html#cfn-iotwireless-serviceprofile-lorawanserviceprofile-drmax
	DrMax int `json:"DrMax,omitempty"`

	// DrMin AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotwireless-serviceprofile-lorawanserviceprofile.html#cfn-iotwireless-serviceprofile-lorawanserviceprofile-drmin
	DrMin int `json:"DrMin,omitempty"`

	// HrAllowed AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotwireless-serviceprofile-lorawanserviceprofile.html#cfn-iotwireless-serviceprofile-lorawanserviceprofile-hrallowed
	HrAllowed bool `json:"HrAllowed,omitempty"`

	// MinGwDiversity AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotwireless-serviceprofile-lorawanserviceprofile.html#cfn-iotwireless-serviceprofile-lorawanserviceprofile-mingwdiversity
	MinGwDiversity int `json:"MinGwDiversity,omitempty"`

	// NwkGeoLoc AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotwireless-serviceprofile-lorawanserviceprofile.html#cfn-iotwireless-serviceprofile-lorawanserviceprofile-nwkgeoloc
	NwkGeoLoc bool `json:"NwkGeoLoc,omitempty"`

	// PrAllowed AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotwireless-serviceprofile-lorawanserviceprofile.html#cfn-iotwireless-serviceprofile-lorawanserviceprofile-prallowed
	PrAllowed bool `json:"PrAllowed,omitempty"`

	// RaAllowed AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotwireless-serviceprofile-lorawanserviceprofile.html#cfn-iotwireless-serviceprofile-lorawanserviceprofile-raallowed
	RaAllowed bool `json:"RaAllowed,omitempty"`

	// ReportDevStatusBattery AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotwireless-serviceprofile-lorawanserviceprofile.html#cfn-iotwireless-serviceprofile-lorawanserviceprofile-reportdevstatusbattery
	ReportDevStatusBattery bool `json:"ReportDevStatusBattery,omitempty"`

	// ReportDevStatusMargin AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotwireless-serviceprofile-lorawanserviceprofile.html#cfn-iotwireless-serviceprofile-lorawanserviceprofile-reportdevstatusmargin
	ReportDevStatusMargin bool `json:"ReportDevStatusMargin,omitempty"`

	// TargetPer AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotwireless-serviceprofile-lorawanserviceprofile.html#cfn-iotwireless-serviceprofile-lorawanserviceprofile-targetper
	TargetPer int `json:"TargetPer,omitempty"`

	// UlBucketSize AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotwireless-serviceprofile-lorawanserviceprofile.html#cfn-iotwireless-serviceprofile-lorawanserviceprofile-ulbucketsize
	UlBucketSize int `json:"UlBucketSize,omitempty"`

	// UlRate AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotwireless-serviceprofile-lorawanserviceprofile.html#cfn-iotwireless-serviceprofile-lorawanserviceprofile-ulrate
	UlRate int `json:"UlRate,omitempty"`

	// UlRatePolicy AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotwireless-serviceprofile-lorawanserviceprofile.html#cfn-iotwireless-serviceprofile-lorawanserviceprofile-ulratepolicy
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
func (r *ServiceProfile_LoRaWANServiceProfile) AWSCloudFormationType() string {
	return "AWS::IoTWireless::ServiceProfile.LoRaWANServiceProfile"
}
