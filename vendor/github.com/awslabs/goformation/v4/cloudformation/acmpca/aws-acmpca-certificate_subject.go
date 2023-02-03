package acmpca

import (
	"github.com/awslabs/goformation/v4/cloudformation/policies"
)

// Certificate_Subject AWS CloudFormation Resource (AWS::ACMPCA::Certificate.Subject)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-acmpca-certificate-subject.html
type Certificate_Subject struct {

	// CommonName AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-acmpca-certificate-subject.html#cfn-acmpca-certificate-subject-commonname
	CommonName string `json:"CommonName,omitempty"`

	// Country AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-acmpca-certificate-subject.html#cfn-acmpca-certificate-subject-country
	Country string `json:"Country,omitempty"`

	// DistinguishedNameQualifier AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-acmpca-certificate-subject.html#cfn-acmpca-certificate-subject-distinguishednamequalifier
	DistinguishedNameQualifier string `json:"DistinguishedNameQualifier,omitempty"`

	// GenerationQualifier AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-acmpca-certificate-subject.html#cfn-acmpca-certificate-subject-generationqualifier
	GenerationQualifier string `json:"GenerationQualifier,omitempty"`

	// GivenName AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-acmpca-certificate-subject.html#cfn-acmpca-certificate-subject-givenname
	GivenName string `json:"GivenName,omitempty"`

	// Initials AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-acmpca-certificate-subject.html#cfn-acmpca-certificate-subject-initials
	Initials string `json:"Initials,omitempty"`

	// Locality AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-acmpca-certificate-subject.html#cfn-acmpca-certificate-subject-locality
	Locality string `json:"Locality,omitempty"`

	// Organization AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-acmpca-certificate-subject.html#cfn-acmpca-certificate-subject-organization
	Organization string `json:"Organization,omitempty"`

	// OrganizationalUnit AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-acmpca-certificate-subject.html#cfn-acmpca-certificate-subject-organizationalunit
	OrganizationalUnit string `json:"OrganizationalUnit,omitempty"`

	// Pseudonym AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-acmpca-certificate-subject.html#cfn-acmpca-certificate-subject-pseudonym
	Pseudonym string `json:"Pseudonym,omitempty"`

	// SerialNumber AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-acmpca-certificate-subject.html#cfn-acmpca-certificate-subject-serialnumber
	SerialNumber string `json:"SerialNumber,omitempty"`

	// State AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-acmpca-certificate-subject.html#cfn-acmpca-certificate-subject-state
	State string `json:"State,omitempty"`

	// Surname AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-acmpca-certificate-subject.html#cfn-acmpca-certificate-subject-surname
	Surname string `json:"Surname,omitempty"`

	// Title AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-acmpca-certificate-subject.html#cfn-acmpca-certificate-subject-title
	Title string `json:"Title,omitempty"`

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
func (r *Certificate_Subject) AWSCloudFormationType() string {
	return "AWS::ACMPCA::Certificate.Subject"
}
