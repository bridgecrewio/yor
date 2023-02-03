package mediaconnect

import (
	"github.com/awslabs/goformation/v4/cloudformation/policies"
)

// Flow_Source AWS CloudFormation Resource (AWS::MediaConnect::Flow.Source)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-mediaconnect-flow-source.html
type Flow_Source struct {

	// Decryption AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-mediaconnect-flow-source.html#cfn-mediaconnect-flow-source-decryption
	Decryption *Flow_Encryption `json:"Decryption,omitempty"`

	// Description AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-mediaconnect-flow-source.html#cfn-mediaconnect-flow-source-description
	Description string `json:"Description,omitempty"`

	// EntitlementArn AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-mediaconnect-flow-source.html#cfn-mediaconnect-flow-source-entitlementarn
	EntitlementArn string `json:"EntitlementArn,omitempty"`

	// IngestIp AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-mediaconnect-flow-source.html#cfn-mediaconnect-flow-source-ingestip
	IngestIp string `json:"IngestIp,omitempty"`

	// IngestPort AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-mediaconnect-flow-source.html#cfn-mediaconnect-flow-source-ingestport
	IngestPort int `json:"IngestPort,omitempty"`

	// MaxBitrate AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-mediaconnect-flow-source.html#cfn-mediaconnect-flow-source-maxbitrate
	MaxBitrate int `json:"MaxBitrate,omitempty"`

	// MaxLatency AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-mediaconnect-flow-source.html#cfn-mediaconnect-flow-source-maxlatency
	MaxLatency int `json:"MaxLatency,omitempty"`

	// Name AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-mediaconnect-flow-source.html#cfn-mediaconnect-flow-source-name
	Name string `json:"Name,omitempty"`

	// Protocol AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-mediaconnect-flow-source.html#cfn-mediaconnect-flow-source-protocol
	Protocol string `json:"Protocol,omitempty"`

	// SourceArn AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-mediaconnect-flow-source.html#cfn-mediaconnect-flow-source-sourcearn
	SourceArn string `json:"SourceArn,omitempty"`

	// StreamId AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-mediaconnect-flow-source.html#cfn-mediaconnect-flow-source-streamid
	StreamId string `json:"StreamId,omitempty"`

	// VpcInterfaceName AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-mediaconnect-flow-source.html#cfn-mediaconnect-flow-source-vpcinterfacename
	VpcInterfaceName string `json:"VpcInterfaceName,omitempty"`

	// WhitelistCidr AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-mediaconnect-flow-source.html#cfn-mediaconnect-flow-source-whitelistcidr
	WhitelistCidr string `json:"WhitelistCidr,omitempty"`

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
func (r *Flow_Source) AWSCloudFormationType() string {
	return "AWS::MediaConnect::Flow.Source"
}
