package serverless

import (
	"github.com/awslabs/goformation/v5/cloudformation/policies"
)

// HttpApi_HttpApiDomainConfiguration AWS CloudFormation Resource (AWS::Serverless::HttpApi.HttpApiDomainConfiguration)
// See: https://github.com/aws/serverless-application-model/blob/master/versions/2016-10-31.md#domain-configuration-object
type HttpApi_HttpApiDomainConfiguration struct {

	// BasePath AWS CloudFormation Property
	// Required: false
	// See: https://github.com/aws/serverless-application-model/blob/master/versions/2016-10-31.md#domain-configuration-object
	BasePath string `json:"BasePath,omitempty"`

	// CertificateArn AWS CloudFormation Property
	// Required: true
	// See: https://github.com/aws/serverless-application-model/blob/master/versions/2016-10-31.md#domain-configuration-object
	CertificateArn string `json:"CertificateArn,omitempty"`

	// DomainName AWS CloudFormation Property
	// Required: true
	// See: https://github.com/aws/serverless-application-model/blob/master/versions/2016-10-31.md#domain-configuration-object
	DomainName string `json:"DomainName,omitempty"`

	// EndpointConfiguration AWS CloudFormation Property
	// Required: false
	// See: https://github.com/aws/serverless-application-model/blob/master/versions/2016-10-31.md#domain-configuration-object
	EndpointConfiguration string `json:"EndpointConfiguration,omitempty"`

	// MutualTlsAuthentication AWS CloudFormation Property
	// Required: false
	// See: https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/sam-property-httpapi-httpapidomainconfiguration.html#sam-httpapi-httpapidomainconfiguration-mutualtlsauthentication
	MutualTlsAuthentication *HttpApi_MutualTlsAuthentication `json:"MutualTlsAuthentication,omitempty"`

	// Route53 AWS CloudFormation Property
	// Required: false
	// See: https://github.com/aws/serverless-application-model/blob/master/versions/2016-10-31.md#domain-configuration-object
	Route53 *HttpApi_Route53Configuration `json:"Route53,omitempty"`

	// SecurityPolicy AWS CloudFormation Property
	// Required: false
	// See: https://github.com/aws/serverless-application-model/blob/master/versions/2016-10-31.md#domain-configuration-object
	SecurityPolicy string `json:"SecurityPolicy,omitempty"`

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
func (r *HttpApi_HttpApiDomainConfiguration) AWSCloudFormationType() string {
	return "AWS::Serverless::HttpApi.HttpApiDomainConfiguration"
}
