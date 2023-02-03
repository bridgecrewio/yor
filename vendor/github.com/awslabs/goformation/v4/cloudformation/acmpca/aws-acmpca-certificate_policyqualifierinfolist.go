package acmpca

import (
	"github.com/awslabs/goformation/v4/cloudformation/policies"
)

// Certificate_PolicyQualifierInfoList AWS CloudFormation Resource (AWS::ACMPCA::Certificate.PolicyQualifierInfoList)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-acmpca-certificate-policyqualifierinfolist.html
type Certificate_PolicyQualifierInfoList struct {

	// PolicyQualifierInfoList AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-acmpca-certificate-policyqualifierinfolist.html#cfn-acmpca-certificate-policyqualifierinfolist-policyqualifierinfolist
	PolicyQualifierInfoList []Certificate_PolicyQualifierInfo `json:"PolicyQualifierInfoList,omitempty"`

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
func (r *Certificate_PolicyQualifierInfoList) AWSCloudFormationType() string {
	return "AWS::ACMPCA::Certificate.PolicyQualifierInfoList"
}
