package acmpca

import (
	"github.com/awslabs/goformation/v4/cloudformation/policies"
)

// Certificate_CertificatePolicyList AWS CloudFormation Resource (AWS::ACMPCA::Certificate.CertificatePolicyList)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-acmpca-certificate-certificatepolicylist.html
type Certificate_CertificatePolicyList struct {

	// CertificatePolicyList AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-acmpca-certificate-certificatepolicylist.html#cfn-acmpca-certificate-certificatepolicylist-certificatepolicylist
	CertificatePolicyList []Certificate_PolicyInformation `json:"CertificatePolicyList,omitempty"`

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
func (r *Certificate_CertificatePolicyList) AWSCloudFormationType() string {
	return "AWS::ACMPCA::Certificate.CertificatePolicyList"
}
