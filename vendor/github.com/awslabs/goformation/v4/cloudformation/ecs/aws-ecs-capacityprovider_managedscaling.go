package ecs

import (
	"github.com/awslabs/goformation/v4/cloudformation/policies"
)

// CapacityProvider_ManagedScaling AWS CloudFormation Resource (AWS::ECS::CapacityProvider.ManagedScaling)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ecs-capacityprovider-managedscaling.html
type CapacityProvider_ManagedScaling struct {

	// InstanceWarmupPeriod AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ecs-capacityprovider-managedscaling.html#cfn-ecs-capacityprovider-managedscaling-instancewarmupperiod
	InstanceWarmupPeriod int `json:"InstanceWarmupPeriod,omitempty"`

	// MaximumScalingStepSize AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ecs-capacityprovider-managedscaling.html#cfn-ecs-capacityprovider-managedscaling-maximumscalingstepsize
	MaximumScalingStepSize int `json:"MaximumScalingStepSize,omitempty"`

	// MinimumScalingStepSize AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ecs-capacityprovider-managedscaling.html#cfn-ecs-capacityprovider-managedscaling-minimumscalingstepsize
	MinimumScalingStepSize int `json:"MinimumScalingStepSize,omitempty"`

	// Status AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ecs-capacityprovider-managedscaling.html#cfn-ecs-capacityprovider-managedscaling-status
	Status string `json:"Status,omitempty"`

	// TargetCapacity AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ecs-capacityprovider-managedscaling.html#cfn-ecs-capacityprovider-managedscaling-targetcapacity
	TargetCapacity int `json:"TargetCapacity,omitempty"`

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
func (r *CapacityProvider_ManagedScaling) AWSCloudFormationType() string {
	return "AWS::ECS::CapacityProvider.ManagedScaling"
}
