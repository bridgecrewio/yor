package serverless

import (
	"github.com/awslabs/goformation/v5/cloudformation/policies"
)

// Function_Auth AWS CloudFormation Resource (AWS::Serverless::Function.Auth)
// See: https://github.com/aws/serverless-application-model/blob/master/versions/2016-10-31.md#function-auth-object
type Function_Auth struct {

	// ApiKeyRequired AWS CloudFormation Property
	// Required: false
	// See: https://github.com/awslabs/serverless-application-model/blob/master/versions/2016-10-31.md#function-auth-object
	ApiKeyRequired bool `json:"ApiKeyRequired,omitempty"`

	// AuthorizationScopes AWS CloudFormation Property
	// Required: false
	// See: https://github.com/awslabs/serverless-application-model/blob/master/versions/2016-10-31.md#function-auth-object
	AuthorizationScopes []string `json:"AuthorizationScopes,omitempty"`

	// Authorizer AWS CloudFormation Property
	// Required: false
	// See: https://github.com/awslabs/serverless-application-model/blob/master/versions/2016-10-31.md#function-auth-object
	Authorizer string `json:"Authorizer,omitempty"`

	// ResourcePolicy AWS CloudFormation Property
	// Required: false
	// See: https://github.com/awslabs/serverless-application-model/blob/master/versions/2016-10-31.md#function-auth-object
	ResourcePolicy *Function_AuthResourcePolicy `json:"ResourcePolicy,omitempty"`

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
func (r *Function_Auth) AWSCloudFormationType() string {
	return "AWS::Serverless::Function.Auth"
}
