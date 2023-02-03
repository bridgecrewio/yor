package ecs

import (
	"github.com/awslabs/goformation/v4/cloudformation/policies"
)

// Cluster_CapacityProviderStrategyItem AWS CloudFormation Resource (AWS::ECS::Cluster.CapacityProviderStrategyItem)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ecs-cluster-capacityproviderstrategyitem.html
type Cluster_CapacityProviderStrategyItem struct {

	// Base AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ecs-cluster-capacityproviderstrategyitem.html#cfn-ecs-cluster-capacityproviderstrategyitem-base
	Base int `json:"Base,omitempty"`

	// CapacityProvider AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ecs-cluster-capacityproviderstrategyitem.html#cfn-ecs-cluster-capacityproviderstrategyitem-capacityprovider
	CapacityProvider string `json:"CapacityProvider,omitempty"`

	// Weight AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ecs-cluster-capacityproviderstrategyitem.html#cfn-ecs-cluster-capacityproviderstrategyitem-weight
	Weight int `json:"Weight,omitempty"`

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
func (r *Cluster_CapacityProviderStrategyItem) AWSCloudFormationType() string {
	return "AWS::ECS::Cluster.CapacityProviderStrategyItem"
}
