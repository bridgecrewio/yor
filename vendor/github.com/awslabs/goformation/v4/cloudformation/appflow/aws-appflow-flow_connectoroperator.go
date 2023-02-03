package appflow

import (
	"github.com/awslabs/goformation/v4/cloudformation/policies"
)

// Flow_ConnectorOperator AWS CloudFormation Resource (AWS::AppFlow::Flow.ConnectorOperator)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appflow-flow-connectoroperator.html
type Flow_ConnectorOperator struct {

	// Amplitude AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appflow-flow-connectoroperator.html#cfn-appflow-flow-connectoroperator-amplitude
	Amplitude string `json:"Amplitude,omitempty"`

	// Datadog AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appflow-flow-connectoroperator.html#cfn-appflow-flow-connectoroperator-datadog
	Datadog string `json:"Datadog,omitempty"`

	// Dynatrace AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appflow-flow-connectoroperator.html#cfn-appflow-flow-connectoroperator-dynatrace
	Dynatrace string `json:"Dynatrace,omitempty"`

	// GoogleAnalytics AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appflow-flow-connectoroperator.html#cfn-appflow-flow-connectoroperator-googleanalytics
	GoogleAnalytics string `json:"GoogleAnalytics,omitempty"`

	// InforNexus AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appflow-flow-connectoroperator.html#cfn-appflow-flow-connectoroperator-infornexus
	InforNexus string `json:"InforNexus,omitempty"`

	// Marketo AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appflow-flow-connectoroperator.html#cfn-appflow-flow-connectoroperator-marketo
	Marketo string `json:"Marketo,omitempty"`

	// S3 AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appflow-flow-connectoroperator.html#cfn-appflow-flow-connectoroperator-s3
	S3 string `json:"S3,omitempty"`

	// Salesforce AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appflow-flow-connectoroperator.html#cfn-appflow-flow-connectoroperator-salesforce
	Salesforce string `json:"Salesforce,omitempty"`

	// ServiceNow AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appflow-flow-connectoroperator.html#cfn-appflow-flow-connectoroperator-servicenow
	ServiceNow string `json:"ServiceNow,omitempty"`

	// Singular AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appflow-flow-connectoroperator.html#cfn-appflow-flow-connectoroperator-singular
	Singular string `json:"Singular,omitempty"`

	// Slack AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appflow-flow-connectoroperator.html#cfn-appflow-flow-connectoroperator-slack
	Slack string `json:"Slack,omitempty"`

	// Trendmicro AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appflow-flow-connectoroperator.html#cfn-appflow-flow-connectoroperator-trendmicro
	Trendmicro string `json:"Trendmicro,omitempty"`

	// Veeva AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appflow-flow-connectoroperator.html#cfn-appflow-flow-connectoroperator-veeva
	Veeva string `json:"Veeva,omitempty"`

	// Zendesk AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appflow-flow-connectoroperator.html#cfn-appflow-flow-connectoroperator-zendesk
	Zendesk string `json:"Zendesk,omitempty"`

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
func (r *Flow_ConnectorOperator) AWSCloudFormationType() string {
	return "AWS::AppFlow::Flow.ConnectorOperator"
}
