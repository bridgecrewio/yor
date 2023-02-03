package acmpca

import (
	"github.com/awslabs/goformation/v4/cloudformation/policies"
)

// CertificateAuthority_KeyUsage AWS CloudFormation Resource (AWS::ACMPCA::CertificateAuthority.KeyUsage)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-acmpca-certificateauthority-keyusage.html
type CertificateAuthority_KeyUsage struct {

	// CRLSign AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-acmpca-certificateauthority-keyusage.html#cfn-acmpca-certificateauthority-keyusage-crlsign
	CRLSign bool `json:"CRLSign,omitempty"`

	// DataEncipherment AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-acmpca-certificateauthority-keyusage.html#cfn-acmpca-certificateauthority-keyusage-dataencipherment
	DataEncipherment bool `json:"DataEncipherment,omitempty"`

	// DecipherOnly AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-acmpca-certificateauthority-keyusage.html#cfn-acmpca-certificateauthority-keyusage-decipheronly
	DecipherOnly bool `json:"DecipherOnly,omitempty"`

	// DigitalSignature AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-acmpca-certificateauthority-keyusage.html#cfn-acmpca-certificateauthority-keyusage-digitalsignature
	DigitalSignature bool `json:"DigitalSignature,omitempty"`

	// EncipherOnly AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-acmpca-certificateauthority-keyusage.html#cfn-acmpca-certificateauthority-keyusage-encipheronly
	EncipherOnly bool `json:"EncipherOnly,omitempty"`

	// KeyAgreement AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-acmpca-certificateauthority-keyusage.html#cfn-acmpca-certificateauthority-keyusage-keyagreement
	KeyAgreement bool `json:"KeyAgreement,omitempty"`

	// KeyCertSign AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-acmpca-certificateauthority-keyusage.html#cfn-acmpca-certificateauthority-keyusage-keycertsign
	KeyCertSign bool `json:"KeyCertSign,omitempty"`

	// KeyEncipherment AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-acmpca-certificateauthority-keyusage.html#cfn-acmpca-certificateauthority-keyusage-keyencipherment
	KeyEncipherment bool `json:"KeyEncipherment,omitempty"`

	// NonRepudiation AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-acmpca-certificateauthority-keyusage.html#cfn-acmpca-certificateauthority-keyusage-nonrepudiation
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
func (r *CertificateAuthority_KeyUsage) AWSCloudFormationType() string {
	return "AWS::ACMPCA::CertificateAuthority.KeyUsage"
}
