package serverless

import (
	"github.com/awslabs/goformation/v5/cloudformation/policies"
)

// HttpApi_Route53Configuration AWS CloudFormation Resource (AWS::Serverless::HttpApi.Route53Configuration)
// See: https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/sam-property-httpapi-route53configuration.html
type HttpApi_Route53Configuration struct {

	// DistributedDomainName AWS CloudFormation Property
	// Required: false
	// See: https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/sam-property-httpapi-route53configuration.html#sam-httpapi-route53configuration-distributiondomainname
	DistributedDomainName string `json:"DistributedDomainName,omitempty"`

	// EvaluateTargetHealth AWS CloudFormation Property
	// Required: false
	// See: https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/sam-property-httpapi-route53configuration.html#sam-httpapi-route53configuration-evaluatetargethealth
	EvaluateTargetHealth bool `json:"EvaluateTargetHealth,omitempty"`

	// HostedZoneId AWS CloudFormation Property
	// Required: false
	// See: https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/sam-property-httpapi-route53configuration.html#sam-httpapi-route53configuration-hostedzoneid
	HostedZoneId string `json:"HostedZoneId,omitempty"`

	// HostedZoneName AWS CloudFormation Property
	// Required: false
	// See: https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/sam-property-httpapi-route53configuration.html#sam-httpapi-route53configuration-hostedzonename
	HostedZoneName string `json:"HostedZoneName,omitempty"`

	// IpV6 AWS CloudFormation Property
	// Required: false
	// See: https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/sam-property-httpapi-route53configuration.html#sam-httpapi-route53configuration-ipv6
	IpV6 bool `json:"IpV6,omitempty"`

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
func (r *HttpApi_Route53Configuration) AWSCloudFormationType() string {
	return "AWS::Serverless::HttpApi.Route53Configuration"
}
