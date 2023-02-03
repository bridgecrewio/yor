package global

import (
	"github.com/awslabs/goformation/v5/cloudformation/serverless"
)

// HttpApi AWS CloudFormation Resource (HttpApi)
// See: https://github.com/aws/serverless-application-model/blob/master/versions/2016-10-31.md#awsserverlesshttpapi
type HttpApi struct {

	// AccessLogSetting AWS CloudFormation Property
	// Required: false
	// See: https://github.com/aws/serverless-application-model/blob/master/versions/2016-10-31.md#awsserverlesshttpapi
	AccessLogSetting *serverless.HttpApi_AccessLogSetting `json:"AccessLogSetting,omitempty"`

	// Auth AWS CloudFormation Property
	// Required: false
	// See: https://github.com/aws/serverless-application-model/blob/master/versions/2016-10-31.md#awsserverlesshttpapi
	Auth *serverless.HttpApi_HttpApiAuth `json:"Auth,omitempty"`

	// CorsConfiguration AWS CloudFormation Property
	// Required: false
	// See: https://github.com/aws/serverless-application-model/blob/master/versions/2016-10-31.md#awsserverlesshttpapi
	CorsConfiguration *serverless.HttpApi_CorsConfiguration `json:"CorsConfiguration,omitempty"`

	// DefaultRouteSettings AWS CloudFormation Property
	// Required: false
	// See: https://github.com/aws/serverless-application-model/blob/master/versions/2016-10-31.md#awsserverlesshttpapi
	DefaultRouteSettings *serverless.HttpApi_RouteSettings `json:"DefaultRouteSettings,omitempty"`

	// Description AWS CloudFormation Property
	// Required: false
	// See: https://github.com/aws/serverless-application-model/blob/master/versions/2016-10-31.md#awsserverlesshttpapi
	Description string `json:"Description,omitempty"`

	// DisableExecuteApiEndpoint AWS CloudFormation Property
	// Required: false
	// See: https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/sam-resource-httpapi.html#sam-httpapi-disableexecuteapiendpoint
	DisableExecuteApiEndpoint bool `json:"DisableExecuteApiEndpoint,omitempty"`

	// Domain AWS CloudFormation Property
	// Required: false
	// See: https://github.com/aws/serverless-application-model/blob/master/versions/2016-10-31.md#awsserverlesshttpapi
	Domain *serverless.HttpApi_HttpApiDomainConfiguration `json:"Domain,omitempty"`

	// FailOnWarnings AWS CloudFormation Property
	// Required: false
	// See: https://github.com/aws/serverless-application-model/blob/master/versions/2016-10-31.md#awsserverlesshttpapi
	FailOnWarnings bool `json:"FailOnWarnings,omitempty"`

	// RouteSettings AWS CloudFormation Property
	// Required: false
	// See: https://github.com/aws/serverless-application-model/blob/master/versions/2016-10-31.md#awsserverlesshttpapi
	RouteSettings *serverless.HttpApi_RouteSettings `json:"RouteSettings,omitempty"`

	// StageVariables AWS CloudFormation Property
	// Required: false
	// See: https://github.com/aws/serverless-application-model/blob/master/versions/2016-10-31.md#awsserverlesshttpapi
	StageVariables map[string]string `json:"StageVariables,omitempty"`

	// Tags AWS CloudFormation Property
	// Required: false
	// See: https://github.com/aws/serverless-application-model/blob/master/versions/2016-10-31.md#awsserverlesshttpapi
	Tags map[string]string `json:"Tags,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *HttpApi) AWSCloudFormationType() string {
	return "HttpApi"
}
