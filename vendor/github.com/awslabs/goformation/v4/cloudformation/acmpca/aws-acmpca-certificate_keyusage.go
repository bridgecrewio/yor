package acmpca

import (
	"github.com/awslabs/goformation/v4/cloudformation/policies"
)

// Certificate_KeyUsage AWS CloudFormation Resource (AWS::ACMPCA::Certificate.KeyUsage)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-acmpca-certificate-keyusage.html
type Certificate_KeyUsage struct {

	// CRLSign AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-acmpca-certificate-keyusage.html#cfn-acmpca-certificate-keyusage-crlsign
	CRLSign bool `json:"CRLSign,omitempty"`

	// DataEncipherment AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-acmpca-certificate-keyusage.html#cfn-acmpca-certificate-keyusage-dataencipherment
	DataEncipherment bool `json:"DataEncipherment,omitempty"`

	// DecipherOnly AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-acmpca-certificate-keyusage.html#cfn-acmpca-certificate-keyusage-decipheronly
	DecipherOnly bool `json:"DecipherOnly,omitempty"`

	// DigitalSignature AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-acmpca-certificate-keyusage.html#cfn-acmpca-certificate-keyusage-digitalsignature
	DigitalSignature bool `json:"DigitalSignature,omitempty"`

	// EncipherOnly AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-acmpca-certificate-keyusage.html#cfn-acmpca-certificate-keyusage-encipheronly
	EncipherOnly bool `json:"EncipherOnly,omitempty"`

	// KeyAgreement AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-acmpca-certificate-keyusage.html#cfn-acmpca-certificate-keyusage-keyagreement
	KeyAgreement bool `json:"KeyAgreement,omitempty"`

	// KeyCertSign AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-acmpca-certificate-keyusage.html#cfn-acmpca-certificate-keyusage-keycertsign
	KeyCertSign bool `json:"KeyCertSign,omitempty"`

	// KeyEncipherment AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-acmpca-certificate-keyusage.html#cfn-acmpca-certificate-keyusage-keyencipherment
	KeyEncipherment bool `json:"KeyEncipherment,omitempty"`

	// NonRepudiation AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-acmpca-certificate-keyusage.html#cfn-acmpca-certificate-keyusage-nonrepudiation
	NonRepudiation bool `json:"NonRepudiation,omitempty"`

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
func (r *Certificate_KeyUsage) AWSCloudFormationType() string {
	return "AWS::ACMPCA::Certificate.KeyUsage"
}
