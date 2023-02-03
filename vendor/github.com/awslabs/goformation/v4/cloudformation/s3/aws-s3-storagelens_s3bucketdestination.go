package s3

import (
	"github.com/awslabs/goformation/v4/cloudformation/policies"
)

// StorageLens_S3BucketDestination AWS CloudFormation Resource (AWS::S3::StorageLens.S3BucketDestination)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-s3-storagelens-s3bucketdestination.html
type StorageLens_S3BucketDestination struct {

	// AccountId AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-s3-storagelens-s3bucketdestination.html#cfn-s3-storagelens-s3bucketdestination-accountid
	AccountId string `json:"AccountId,omitempty"`

	// Arn AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-s3-storagelens-s3bucketdestination.html#cfn-s3-storagelens-s3bucketdestination-arn
	Arn string `json:"Arn,omitempty"`

	// Encryption AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-s3-storagelens-s3bucketdestination.html#cfn-s3-storagelens-s3bucketdestination-encryption
	Encryption *StorageLens_Encryption `json:"Encryption,omitempty"`

	// Format AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-s3-storagelens-s3bucketdestination.html#cfn-s3-storagelens-s3bucketdestination-format
	Format string `json:"Format,omitempty"`

	// OutputSchemaVersion AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-s3-storagelens-s3bucketdestination.html#cfn-s3-storagelens-s3bucketdestination-outputschemaversion
	OutputSchemaVersion string `json:"OutputSchemaVersion,omitempty"`

	// Prefix AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-s3-storagelens-s3bucketdestination.html#cfn-s3-storagelens-s3bucketdestination-prefix
	Prefix string `json:"Prefix,omitempty"`

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
func (r *StorageLens_S3BucketDestination) AWSCloudFormationType() string {
	return "AWS::S3::StorageLens.S3BucketDestination"
}
