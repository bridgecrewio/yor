package serverless

import (
	"github.com/awslabs/goformation/v5/cloudformation/policies"
)

// Function_AuthResourcePolicy AWS CloudFormation Resource (AWS::Serverless::Function.AuthResourcePolicy)
// See: https://github.com/aws/serverless-application-model/blob/master/versions/2016-10-31.md#function-auth-object
type Function_AuthResourcePolicy struct {

	// AwsAccountBlacklist AWS CloudFormation Property
	// Required: false
	// See: https://github.com/awslabs/serverless-application-model/blob/master/versions/2016-10-31.md#function-auth-object
	AwsAccountBlacklist []string `json:"AwsAccountBlacklist,omitempty"`

	// AwsAccountWhitelist AWS CloudFormation Property
	// Required: false
	// See: https://github.com/awslabs/serverless-application-model/blob/master/versions/2016-10-31.md#function-auth-object
	AwsAccountWhitelist []string `json:"AwsAccountWhitelist,omitempty"`

	// CustomStatements AWS CloudFormation Property
	// Required: false
	// See: https://github.com/awslabs/serverless-application-model/blob/master/versions/2016-10-31.md#function-auth-object
	CustomStatements []interface{} `json:"CustomStatements,omitempty"`

	// IntrinsicVpcBlacklist AWS CloudFormation Property
	// Required: false
	// See: https://github.com/awslabs/serverless-application-model/blob/master/versions/2016-10-31.md#function-auth-object
	IntrinsicVpcBlacklist []string `json:"IntrinsicVpcBlacklist,omitempty"`

	// IntrinsicVpcWhitelist AWS CloudFormation Property
	// Required: false
	// See: https://github.com/awslabs/serverless-application-model/blob/master/versions/2016-10-31.md#function-auth-object
	IntrinsicVpcWhitelist []string `json:"IntrinsicVpcWhitelist,omitempty"`

	// IntrinsicVpceBlacklist AWS CloudFormation Property
	// Required: false
	// See: https://github.com/awslabs/serverless-application-model/blob/master/versions/2016-10-31.md#function-auth-object
	IntrinsicVpceBlacklist []string `json:"IntrinsicVpceBlacklist,omitempty"`

	// IntrinsicVpceWhitelist AWS CloudFormation Property
	// Required: false
	// See: https://github.com/awslabs/serverless-application-model/blob/master/versions/2016-10-31.md#function-auth-object
	IntrinsicVpceWhitelist []string `json:"IntrinsicVpceWhitelist,omitempty"`

	// IpRangeBlacklist AWS CloudFormation Property
	// Required: false
	// See: https://github.com/awslabs/serverless-application-model/blob/master/versions/2016-10-31.md#function-auth-object
	IpRangeBlacklist []string `json:"IpRangeBlacklist,omitempty"`

	// IpRangeWhitelist AWS CloudFormation Property
	// Required: false
	// See: https://github.com/awslabs/serverless-application-model/blob/master/versions/2016-10-31.md#function-auth-object
	IpRangeWhitelist []string `json:"IpRangeWhitelist,omitempty"`

	// SourceVpcBlacklist AWS CloudFormation Property
	// Required: false
	// See: https://github.com/awslabs/serverless-application-model/blob/master/versions/2016-10-31.md#function-auth-object
	SourceVpcBlacklist []string `json:"SourceVpcBlacklist,omitempty"`

	// SourceVpcWhitelist AWS CloudFormation Property
	// Required: false
	// See: https://github.com/awslabs/serverless-application-model/blob/master/versions/2016-10-31.md#function-auth-object
	SourceVpcWhitelist []string `json:"SourceVpcWhitelist,omitempty"`

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
func (r *Function_AuthResourcePolicy) AWSCloudFormationType() string {
	return "AWS::Serverless::Function.AuthResourcePolicy"
}
