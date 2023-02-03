package emr

import (
	"github.com/awslabs/goformation/v4/cloudformation/policies"
)

// Cluster_ComputeLimits AWS CloudFormation Resource (AWS::EMR::Cluster.ComputeLimits)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-elasticmapreduce-cluster-computelimits.html
type Cluster_ComputeLimits struct {

	// MaximumCapacityUnits AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-elasticmapreduce-cluster-computelimits.html#cfn-elasticmapreduce-cluster-computelimits-maximumcapacityunits
	MaximumCapacityUnits int `json:"MaximumCapacityUnits"`

	// MaximumCoreCapacityUnits AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-elasticmapreduce-cluster-computelimits.html#cfn-elasticmapreduce-cluster-computelimits-maximumcorecapacityunits
	MaximumCoreCapacityUnits int `json:"MaximumCoreCapacityUnits,omitempty"`

	// MaximumOnDemandCapacityUnits AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-elasticmapreduce-cluster-computelimits.html#cfn-elasticmapreduce-cluster-computelimits-maximumondemandcapacityunits
	MaximumOnDemandCapacityUnits int `json:"MaximumOnDemandCapacityUnits,omitempty"`

	// MinimumCapacityUnits AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-elasticmapreduce-cluster-computelimits.html#cfn-elasticmapreduce-cluster-computelimits-minimumcapacityunits
	MinimumCapacityUnits int `json:"MinimumCapacityUnits"`

	// UnitType AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-elasticmapreduce-cluster-computelimits.html#cfn-elasticmapreduce-cluster-computelimits-unittype
	UnitType string `json:"UnitType,omitempty"`

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
func (r *Cluster_ComputeLimits) AWSCloudFormationType() string {
	return "AWS::EMR::Cluster.ComputeLimits"
}
