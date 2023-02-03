package imagebuilder

import (
	"github.com/awslabs/goformation/v4/cloudformation/policies"
)

// ImageRecipe_EbsInstanceBlockDeviceSpecification AWS CloudFormation Resource (AWS::ImageBuilder::ImageRecipe.EbsInstanceBlockDeviceSpecification)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-imagebuilder-imagerecipe-ebsinstanceblockdevicespecification.html
type ImageRecipe_EbsInstanceBlockDeviceSpecification struct {

	// DeleteOnTermination AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-imagebuilder-imagerecipe-ebsinstanceblockdevicespecification.html#cfn-imagebuilder-imagerecipe-ebsinstanceblockdevicespecification-deleteontermination
	DeleteOnTermination bool `json:"DeleteOnTermination,omitempty"`

	// Encrypted AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-imagebuilder-imagerecipe-ebsinstanceblockdevicespecification.html#cfn-imagebuilder-imagerecipe-ebsinstanceblockdevicespecification-encrypted
	Encrypted bool `json:"Encrypted,omitempty"`

	// Iops AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-imagebuilder-imagerecipe-ebsinstanceblockdevicespecification.html#cfn-imagebuilder-imagerecipe-ebsinstanceblockdevicespecification-iops
	Iops int `json:"Iops,omitempty"`

	// KmsKeyId AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-imagebuilder-imagerecipe-ebsinstanceblockdevicespecification.html#cfn-imagebuilder-imagerecipe-ebsinstanceblockdevicespecification-kmskeyid
	KmsKeyId string `json:"KmsKeyId,omitempty"`

	// SnapshotId AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-imagebuilder-imagerecipe-ebsinstanceblockdevicespecification.html#cfn-imagebuilder-imagerecipe-ebsinstanceblockdevicespecification-snapshotid
	SnapshotId string `json:"SnapshotId,omitempty"`

	// VolumeSize AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-imagebuilder-imagerecipe-ebsinstanceblockdevicespecification.html#cfn-imagebuilder-imagerecipe-ebsinstanceblockdevicespecification-volumesize
	VolumeSize int `json:"VolumeSize,omitempty"`

	// VolumeType AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-imagebuilder-imagerecipe-ebsinstanceblockdevicespecification.html#cfn-imagebuilder-imagerecipe-ebsinstanceblockdevicespecification-volumetype
	VolumeType string `json:"VolumeType,omitempty"`

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
func (r *ImageRecipe_EbsInstanceBlockDeviceSpecification) AWSCloudFormationType() string {
	return "AWS::ImageBuilder::ImageRecipe.EbsInstanceBlockDeviceSpecification"
}
