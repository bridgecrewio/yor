package serverless

import (
	"github.com/awslabs/goformation/v5/cloudformation/policies"
)

// Api_EndpointConfiguration AWS CloudFormation Resource (AWS::Serverless::Api.EndpointConfiguration)
// See: https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/sam-property-api-endpointconfiguration.html
type Api_EndpointConfiguration struct {

	// Type AWS CloudFormation Property
	// Required: false
	// See: https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/sam-property-api-endpointconfiguration.html#sam-api-endpointconfiguration-type
	Type string `json:"Type,omitempty"`

	// VpcEndpointIds AWS CloudFormation Property
	// Required: false
	// See: https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/sam-property-api-endpointconfiguration.html#sam-api-endpointconfiguration-vpcendpointids
	VpcEndpointIds []string `json:"VpcEndpointIds,omitempty"`

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
func (r *Api_EndpointConfiguration) AWSCloudFormationType() string {
	return "AWS::Serverless::Api.EndpointConfiguration"
}
