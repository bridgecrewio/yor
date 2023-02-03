package global

import (
	"github.com/awslabs/goformation/v5/cloudformation/serverless"
)

// Api AWS CloudFormation Resource (Api)
// See: https://github.com/awslabs/serverless-application-model/blob/master/versions/2016-10-31.md#awsserverlessapi
type Api struct {

	// AccessLogSetting AWS CloudFormation Property
	// Required: false
	// See: https://github.com/awslabs/serverless-application-model/blob/master/versions/2016-10-31.md#awsserverlessapi
	AccessLogSetting *serverless.Api_AccessLogSetting `json:"AccessLogSetting,omitempty"`

	// Auth AWS CloudFormation Property
	// Required: false
	// See: https://github.com/awslabs/serverless-application-model/blob/master/versions/2016-10-31.md#awsserverlessapi
	Auth *serverless.Api_Auth `json:"Auth,omitempty"`

	// BinaryMediaTypes AWS CloudFormation Property
	// Required: false
	// See: https://github.com/awslabs/serverless-application-model/blob/master/versions/2016-10-31.md#awsserverlessapi
	BinaryMediaTypes []string `json:"BinaryMediaTypes,omitempty"`

	// CacheClusterEnabled AWS CloudFormation Property
	// Required: false
	// See: https://github.com/awslabs/serverless-application-model/blob/master/versions/2016-10-31.md#awsserverlessapi
	CacheClusterEnabled bool `json:"CacheClusterEnabled,omitempty"`

	// CacheClusterSize AWS CloudFormation Property
	// Required: false
	// See: https://github.com/awslabs/serverless-application-model/blob/master/versions/2016-10-31.md#awsserverlessapi
	CacheClusterSize string `json:"CacheClusterSize,omitempty"`

	// CanarySetting AWS CloudFormation Property
	// Required: false
	// See: https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/sam-resource-api.html#sam-api-canarysetting
	CanarySetting *serverless.Api_CanarySetting `json:"CanarySetting,omitempty"`

	// Cors AWS CloudFormation Property
	// Required: false
	// See: https://github.com/awslabs/serverless-application-model/blob/master/versions/2016-10-31.md#awsserverlessapi
	Cors *serverless.Api_Cors `json:"Cors,omitempty"`

	// DefinitionUri AWS CloudFormation Property
	// Required: false
	// See: https://github.com/awslabs/serverless-application-model/blob/master/versions/2016-10-31.md#awsserverlessapi
	DefinitionUri *serverless.Api_DefinitionUri `json:"DefinitionUri,omitempty"`

	// Description AWS CloudFormation Property
	// Required: false
	// See: https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/sam-resource-api.html#sam-api-description
	Description string `json:"Description,omitempty"`

	// EndpointConfiguration AWS CloudFormation Property
	// Required: false
	// See: https://github.com/awslabs/serverless-application-model/blob/master/versions/2016-10-31.md#awsserverlessapi
	EndpointConfiguration *serverless.Api_EndpointConfiguration `json:"EndpointConfiguration,omitempty"`

	// GatewayResponses AWS CloudFormation Property
	// Required: false
	// See: https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/sam-resource-api.html#sam-api-gatewayresponses
	GatewayResponses map[string]string `json:"GatewayResponses,omitempty"`

	// MethodSettings AWS CloudFormation Property
	// Required: false
	// See: https://github.com/awslabs/serverless-application-model/blob/master/versions/2016-10-31.md#awsserverlessapi
	MethodSettings []interface{} `json:"MethodSettings,omitempty"`

	// MinimumCompressionSize AWS CloudFormation Property
	// Required: false
	// See: https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/sam-resource-api.html#sam-api-minimumcompressionsize
	MinimumCompressionSize int `json:"MinimumCompressionSize,omitempty"`

	// Models AWS CloudFormation Property
	// Required: false
	// See: https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/sam-resource-api.html#sam-api-models
	Models map[string]string `json:"Models,omitempty"`

	// Name AWS CloudFormation Property
	// Required: false
	// See: https://github.com/awslabs/serverless-application-model/blob/master/versions/2016-10-31.md#awsserverlessapi
	Name string `json:"Name,omitempty"`

	// OpenApiVersion AWS CloudFormation Property
	// Required: false
	// See: https://github.com/awslabs/serverless-application-model/blob/master/versions/2016-10-31.md#awsserverlessapi
	OpenApiVersion string `json:"OpenApiVersion,omitempty"`

	// Tags AWS CloudFormation Property
	// Required: false
	// See: https://github.com/aws/serverless-application-model/blob/master/versions/2016-10-31.md#awsserverlesshttpapi
	Tags map[string]string `json:"Tags,omitempty"`

	// TracingEnabled AWS CloudFormation Property
	// Required: false
	// See: https://github.com/awslabs/serverless-application-model/blob/master/versions/2016-10-31.md#awsserverlessapi
	TracingEnabled bool `json:"TracingEnabled,omitempty"`

	// Variables AWS CloudFormation Property
	// Required: false
	// See: https://github.com/awslabs/serverless-application-model/blob/master/versions/2016-10-31.md#awsserverlessapi
	Variables map[string]string `json:"Variables,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *Api) AWSCloudFormationType() string {
	return "Api"
}
