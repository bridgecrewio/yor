package dms

import (
	"github.com/awslabs/goformation/v4/cloudformation/policies"
)

// Endpoint_KafkaSettings AWS CloudFormation Resource (AWS::DMS::Endpoint.KafkaSettings)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-dms-endpoint-kafkasettings.html
type Endpoint_KafkaSettings struct {

	// Broker AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-dms-endpoint-kafkasettings.html#cfn-dms-endpoint-kafkasettings-broker
	Broker string `json:"Broker,omitempty"`

	// Topic AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-dms-endpoint-kafkasettings.html#cfn-dms-endpoint-kafkasettings-topic
	Topic string `json:"Topic,omitempty"`

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
func (r *Endpoint_KafkaSettings) AWSCloudFormationType() string {
	return "AWS::DMS::Endpoint.KafkaSettings"
}
