package serverless

import (
	"github.com/awslabs/goformation/v5/cloudformation/policies"
)

// HttpApi_HttpApiAuth AWS CloudFormation Resource (AWS::Serverless::HttpApi.HttpApiAuth)
// See: https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/sam-property-httpapi-httpapiauth.html
type HttpApi_HttpApiAuth struct {

	// Authorizers AWS CloudFormation Property
	// Required: false
	// See: https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/sam-property-httpapi-httpapiauth.html#sam-httpapi-httpapiauth-defaultauthorizer
	Authorizers interface{} `json:"Authorizers,omitempty"`

	// DefaultAuthorizer AWS CloudFormation Property
	// Required: false
	// See: https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/sam-property-httpapi-httpapiauth.html#sam-httpapi-httpapiauth-authorizers
	DefaultAuthorizer string `json:"DefaultAuthorizer,omitempty"`

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
func (r *HttpApi_HttpApiAuth) AWSCloudFormationType() string {
	return "AWS::Serverless::HttpApi.HttpApiAuth"
}
