package serverless

import (
	"github.com/awslabs/goformation/v5/cloudformation/policies"
)

// HttpApi_CorsConfigurationObject AWS CloudFormation Resource (AWS::Serverless::HttpApi.CorsConfigurationObject)
// See: https://github.com/aws/serverless-application-model/blob/master/versions/2016-10-31.md#cors-configuration-object
type HttpApi_CorsConfigurationObject struct {

	// AllowCredentials AWS CloudFormation Property
	// Required: false
	// See: https://github.com/aws/serverless-application-model/blob/master/versions/2016-10-31.md#cors-configuration-object
	AllowCredentials bool `json:"AllowCredentials,omitempty"`

	// AllowHeaders AWS CloudFormation Property
	// Required: false
	// See: https://github.com/aws/serverless-application-model/blob/master/versions/2016-10-31.md#cors-configuration-object
	AllowHeaders string `json:"AllowHeaders,omitempty"`

	// AllowMethods AWS CloudFormation Property
	// Required: false
	// See: https://github.com/aws/serverless-application-model/blob/master/versions/2016-10-31.md#cors-configuration-object
	AllowMethods string `json:"AllowMethods,omitempty"`

	// AllowOrigin AWS CloudFormation Property
	// Required: false
	// See: https://github.com/aws/serverless-application-model/blob/master/versions/2016-10-31.md#cors-configuration-object
	AllowOrigin string `json:"AllowOrigin,omitempty"`

	// ExposeHeaders AWS CloudFormation Property
	// Required: false
	// See: https://github.com/aws/serverless-application-model/blob/master/versions/2016-10-31.md#cors-configuration-object
	ExposeHeaders []string `json:"ExposeHeaders,omitempty"`

	// MaxAge AWS CloudFormation Property
	// Required: false
	// See: https://github.com/aws/serverless-application-model/blob/master/versions/2016-10-31.md#cors-configuration-object
	MaxAge string `json:"MaxAge,omitempty"`

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
func (r *HttpApi_CorsConfigurationObject) AWSCloudFormationType() string {
	return "AWS::Serverless::HttpApi.CorsConfigurationObject"
}
