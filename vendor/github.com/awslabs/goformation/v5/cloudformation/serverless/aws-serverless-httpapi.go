package serverless

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/awslabs/goformation/v5/cloudformation/policies"
)

// HttpApi AWS CloudFormation Resource (AWS::Serverless::HttpApi)
// See: https://github.com/aws/serverless-application-model/blob/master/versions/2016-10-31.md#awsserverlesshttpapi
type HttpApi struct {

	// AccessLogSetting AWS CloudFormation Property
	// Required: false
	// See: https://github.com/aws/serverless-application-model/blob/master/versions/2016-10-31.md#awsserverlesshttpapi
	AccessLogSetting *HttpApi_AccessLogSetting `json:"AccessLogSetting,omitempty"`

	// Auth AWS CloudFormation Property
	// Required: false
	// See: https://github.com/aws/serverless-application-model/blob/master/versions/2016-10-31.md#awsserverlesshttpapi
	Auth *HttpApi_HttpApiAuth `json:"Auth,omitempty"`

	// CorsConfiguration AWS CloudFormation Property
	// Required: false
	// See: https://github.com/aws/serverless-application-model/blob/master/versions/2016-10-31.md#awsserverlesshttpapi
	CorsConfiguration *HttpApi_CorsConfiguration `json:"CorsConfiguration,omitempty"`

	// DefaultRouteSettings AWS CloudFormation Property
	// Required: false
	// See: https://github.com/aws/serverless-application-model/blob/master/versions/2016-10-31.md#awsserverlesshttpapi
	DefaultRouteSettings *HttpApi_RouteSettings `json:"DefaultRouteSettings,omitempty"`

	// DefinitionBody AWS CloudFormation Property
	// Required: false
	// See: https://github.com/aws/serverless-application-model/blob/master/versions/2016-10-31.md#awsserverlesshttpapi
	DefinitionBody interface{} `json:"DefinitionBody,omitempty"`

	// DefinitionUri AWS CloudFormation Property
	// Required: false
	// See: https://github.com/aws/serverless-application-model/blob/master/versions/2016-10-31.md#awsserverlesshttpapi
	DefinitionUri *HttpApi_DefinitionUri `json:"DefinitionUri,omitempty"`

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
	Domain *HttpApi_HttpApiDomainConfiguration `json:"Domain,omitempty"`

	// FailOnWarnings AWS CloudFormation Property
	// Required: false
	// See: https://github.com/aws/serverless-application-model/blob/master/versions/2016-10-31.md#awsserverlesshttpapi
	FailOnWarnings bool `json:"FailOnWarnings,omitempty"`

	// RouteSettings AWS CloudFormation Property
	// Required: false
	// See: https://github.com/aws/serverless-application-model/blob/master/versions/2016-10-31.md#awsserverlesshttpapi
	RouteSettings *HttpApi_RouteSettings `json:"RouteSettings,omitempty"`

	// StageName AWS CloudFormation Property
	// Required: false
	// See: https://github.com/aws/serverless-application-model/blob/master/versions/2016-10-31.md#awsserverlesshttpapi
	StageName string `json:"StageName,omitempty"`

	// StageVariables AWS CloudFormation Property
	// Required: false
	// See: https://github.com/aws/serverless-application-model/blob/master/versions/2016-10-31.md#awsserverlesshttpapi
	StageVariables map[string]string `json:"StageVariables,omitempty"`

	// Tags AWS CloudFormation Property
	// Required: false
	// See: https://github.com/aws/serverless-application-model/blob/master/versions/2016-10-31.md#awsserverlesshttpapi
	Tags map[string]string `json:"Tags,omitempty"`

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
func (r *HttpApi) AWSCloudFormationType() string {
	return "AWS::Serverless::HttpApi"
}

// MarshalJSON is a custom JSON marshalling hook that embeds this object into
// an AWS CloudFormation JSON resource's 'Properties' field and adds a 'Type'.
func (r HttpApi) MarshalJSON() ([]byte, error) {
	type Properties HttpApi
	return json.Marshal(&struct {
		Type                string
		Properties          Properties
		DependsOn           []string                     `json:"DependsOn,omitempty"`
		Metadata            map[string]interface{}       `json:"Metadata,omitempty"`
		DeletionPolicy      policies.DeletionPolicy      `json:"DeletionPolicy,omitempty"`
		UpdateReplacePolicy policies.UpdateReplacePolicy `json:"UpdateReplacePolicy,omitempty"`
		Condition           string                       `json:"Condition,omitempty"`
	}{
		Type:                r.AWSCloudFormationType(),
		Properties:          (Properties)(r),
		DependsOn:           r.AWSCloudFormationDependsOn,
		Metadata:            r.AWSCloudFormationMetadata,
		DeletionPolicy:      r.AWSCloudFormationDeletionPolicy,
		UpdateReplacePolicy: r.AWSCloudFormationUpdateReplacePolicy,
		Condition:           r.AWSCloudFormationCondition,
	})
}

// UnmarshalJSON is a custom JSON unmarshalling hook that strips the outer
// AWS CloudFormation resource object, and just keeps the 'Properties' field.
func (r *HttpApi) UnmarshalJSON(b []byte) error {
	type Properties HttpApi
	res := &struct {
		Type                string
		Properties          *Properties
		DependsOn           []string
		Metadata            map[string]interface{}
		DeletionPolicy      string
		UpdateReplacePolicy string
		Condition           string
	}{}

	dec := json.NewDecoder(bytes.NewReader(b))
	dec.DisallowUnknownFields() // Force error if unknown field is found

	if err := dec.Decode(&res); err != nil {
		fmt.Printf("ERROR: %s\n", err)
		return err
	}

	// If the resource has no Properties set, it could be nil
	if res.Properties != nil {
		*r = HttpApi(*res.Properties)
	}
	if res.DependsOn != nil {
		r.AWSCloudFormationDependsOn = res.DependsOn
	}
	if res.Metadata != nil {
		r.AWSCloudFormationMetadata = res.Metadata
	}
	if res.DeletionPolicy != "" {
		r.AWSCloudFormationDeletionPolicy = policies.DeletionPolicy(res.DeletionPolicy)
	}
	if res.UpdateReplacePolicy != "" {
		r.AWSCloudFormationUpdateReplacePolicy = policies.UpdateReplacePolicy(res.UpdateReplacePolicy)
	}
	if res.Condition != "" {
		r.AWSCloudFormationCondition = res.Condition
	}
	return nil
}
