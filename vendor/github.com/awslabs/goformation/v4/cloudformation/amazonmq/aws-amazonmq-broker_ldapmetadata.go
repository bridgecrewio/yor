package amazonmq

import (
	"github.com/awslabs/goformation/v4/cloudformation/policies"
)

// Broker_LdapMetadata AWS CloudFormation Resource (AWS::AmazonMQ::Broker.LdapMetadata)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-amazonmq-broker-ldapmetadata.html
type Broker_LdapMetadata struct {

	// InterBrokerCreds AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-amazonmq-broker-ldapmetadata.html#cfn-amazonmq-broker-ldapmetadata-interbrokercreds
	InterBrokerCreds []Broker_InterBrokerCred `json:"InterBrokerCreds,omitempty"`

	// ServerMetadata AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-amazonmq-broker-ldapmetadata.html#cfn-amazonmq-broker-ldapmetadata-servermetadata
	ServerMetadata *Broker_ServerMetadata `json:"ServerMetadata,omitempty"`

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
func (r *Broker_LdapMetadata) AWSCloudFormationType() string {
	return "AWS::AmazonMQ::Broker.LdapMetadata"
}
