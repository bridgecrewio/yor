package global

import (
	"github.com/awslabs/goformation/v5/cloudformation/serverless"
)

// SimpleTable AWS CloudFormation Resource (SimpleTable)
// See: https://github.com/awslabs/serverless-application-model/blob/master/versions/2016-10-31.md#awsserverlesssimpletable
type SimpleTable struct {

	// SSESpecification AWS CloudFormation Property
	// Required: false
	// See: https://github.com/awslabs/serverless-application-model/blob/master/versions/2016-10-31.md#awsserverlesssimpletable
	SSESpecification *serverless.SimpleTable_SSESpecification `json:"SSESpecification,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *SimpleTable) AWSCloudFormationType() string {
	return "SimpleTable"
}
