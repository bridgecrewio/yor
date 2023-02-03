package global

import (
	"github.com/awslabs/goformation/v5/cloudformation/serverless"
)

// Function AWS CloudFormation Resource (Function)
// See: https://github.com/awslabs/serverless-application-model/blob/master/versions/2016-10-31.md#awsserverlessfunction
type Function struct {

	// AssumeRolePolicyDocument AWS CloudFormation Property
	// Required: false
	// See: https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/sam-resource-function.html#sam-function-assumerolepolicydocument
	AssumeRolePolicyDocument interface{} `json:"AssumeRolePolicyDocument,omitempty"`

	// AutoPublishAlias AWS CloudFormation Property
	// Required: false
	// See: https://github.com/awslabs/serverless-application-model/blob/master/versions/2016-10-31.md#awsserverlessfunction
	AutoPublishAlias string `json:"AutoPublishAlias,omitempty"`

	// AutoPublishCodeSha256 AWS CloudFormation Property
	// Required: false
	// See: https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/sam-resource-function.html#sam-function-autopublishcodesha256
	AutoPublishCodeSha256 string `json:"AutoPublishCodeSha256,omitempty"`

	// CodeSigningConfigArn AWS CloudFormation Property
	// Required: false
	// See: https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/sam-resource-function.html#sam-function-codesigningconfigarn
	CodeSigningConfigArn string `json:"CodeSigningConfigArn,omitempty"`

	// CodeUri AWS CloudFormation Property
	// Required: false
	// See: https://github.com/awslabs/serverless-application-model/blob/master/versions/2016-10-31.md#awsserverlessfunction
	CodeUri *serverless.Function_CodeUri `json:"CodeUri,omitempty"`

	// DeadLetterQueue AWS CloudFormation Property
	// Required: false
	// See: https://github.com/awslabs/serverless-application-model/blob/master/versions/2016-10-31.md#awsserverlessfunction
	DeadLetterQueue *serverless.Function_DeadLetterQueue `json:"DeadLetterQueue,omitempty"`

	// DeploymentPreference AWS CloudFormation Property
	// Required: false
	// See: https://github.com/awslabs/serverless-application-model/blob/master/versions/2016-10-31.md#awsserverlessfunction
	DeploymentPreference *serverless.Function_DeploymentPreference `json:"DeploymentPreference,omitempty"`

	// Description AWS CloudFormation Property
	// Required: false
	// See: https://github.com/awslabs/serverless-application-model/blob/master/versions/2016-10-31.md#awsserverlessfunction
	Description string `json:"Description,omitempty"`

	// Environment AWS CloudFormation Property
	// Required: false
	// See: https://github.com/awslabs/serverless-application-model/blob/master/versions/2016-10-31.md#awsserverlessfunction
	Environment *serverless.Function_FunctionEnvironment `json:"Environment,omitempty"`

	// EventInvokeConfig AWS CloudFormation Property
	// Required: false
	// See: https://github.com/awslabs/serverless-application-model/blob/master/versions/2016-10-31.md#awsserverlessfunction
	EventInvokeConfig *serverless.Function_EventInvokeConfig `json:"EventInvokeConfig,omitempty"`

	// Handler AWS CloudFormation Property
	// Required: false
	// See: https://github.com/awslabs/serverless-application-model/blob/master/versions/2016-10-31.md#awsserverlessfunction
	Handler string `json:"Handler,omitempty"`

	// ImageConfig AWS CloudFormation Property
	// Required: false
	// See: https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/sam-resource-function.html#sam-function-imageconfig
	ImageConfig *serverless.Function_ImageConfig `json:"ImageConfig,omitempty"`

	// ImageUri AWS CloudFormation Property
	// Required: false
	// See: https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/sam-resource-function.html#sam-function-imageuri
	ImageUri string `json:"ImageUri,omitempty"`

	// InlineCode AWS CloudFormation Property
	// Required: false
	// See: https://github.com/awslabs/serverless-application-model/blob/master/versions/2016-10-31.md#awsserverlessfunction
	InlineCode string `json:"InlineCode,omitempty"`

	// KmsKeyArn AWS CloudFormation Property
	// Required: false
	// See: https://github.com/awslabs/serverless-application-model/blob/master/versions/2016-10-31.md#awsserverlessfunction
	KmsKeyArn string `json:"KmsKeyArn,omitempty"`

	// Layers AWS CloudFormation Property
	// Required: false
	// See: https://github.com/awslabs/serverless-application-model/blob/master/versions/2016-10-31.md#awsserverlessfunction
	Layers []string `json:"Layers,omitempty"`

	// MemorySize AWS CloudFormation Property
	// Required: false
	// See: https://github.com/awslabs/serverless-application-model/blob/master/versions/2016-10-31.md#awsserverlessfunction
	MemorySize int `json:"MemorySize,omitempty"`

	// PackageType AWS CloudFormation Property
	// Required: false
	// See: https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/sam-resource-function.html#sam-function-packagetype
	PackageType string `json:"PackageType,omitempty"`

	// PermissionsBoundary AWS CloudFormation Property
	// Required: false
	// See: https://github.com/awslabs/serverless-application-model/blob/master/versions/2016-10-31.md#awsserverlessfunction
	PermissionsBoundary string `json:"PermissionsBoundary,omitempty"`

	// ReservedConcurrentExecutions AWS CloudFormation Property
	// Required: false
	// See: https://github.com/awslabs/serverless-application-model/blob/master/versions/2016-10-31.md#awsserverlessfunction
	ReservedConcurrentExecutions int `json:"ReservedConcurrentExecutions,omitempty"`

	// Runtime AWS CloudFormation Property
	// Required: false
	// See: https://github.com/awslabs/serverless-application-model/blob/master/versions/2016-10-31.md#awsserverlessfunction
	Runtime string `json:"Runtime,omitempty"`

	// Tags AWS CloudFormation Property
	// Required: false
	// See: https://github.com/awslabs/serverless-application-model/blob/master/versions/2016-10-31.md#awsserverlessfunction
	Tags map[string]string `json:"Tags,omitempty"`

	// Timeout AWS CloudFormation Property
	// Required: false
	// See: https://github.com/awslabs/serverless-application-model/blob/master/versions/2016-10-31.md#awsserverlessfunction
	Timeout int `json:"Timeout,omitempty"`

	// Tracing AWS CloudFormation Property
	// Required: false
	// See: https://github.com/awslabs/serverless-application-model/blob/master/versions/2016-10-31.md#awsserverlessfunction
	Tracing string `json:"Tracing,omitempty"`

	// VersionDescription AWS CloudFormation Property
	// Required: false
	// See: https://github.com/awslabs/serverless-application-model/blob/master/versions/2016-10-31.md#awsserverlessfunction
	VersionDescription string `json:"VersionDescription,omitempty"`

	// VpcConfig AWS CloudFormation Property
	// Required: false
	// See: https://github.com/awslabs/serverless-application-model/blob/master/versions/2016-10-31.md#awsserverlessfunction
	VpcConfig *serverless.Function_VpcConfig `json:"VpcConfig,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *Function) AWSCloudFormationType() string {
	return "Function"
}
