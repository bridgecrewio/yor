package elasticache

import (
	"github.com/awslabs/goformation/v4/cloudformation/policies"
)

// UserGroup_ReplicationGroupIdList AWS CloudFormation Resource (AWS::ElastiCache::UserGroup.ReplicationGroupIdList)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-elasticache-usergroup-replicationgroupidlist.html
type UserGroup_ReplicationGroupIdList struct {

	// ReplicationGroupIdList AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-elasticache-usergroup-replicationgroupidlist.html#cfn-elasticache-usergroup-replicationgroupidlist-replicationgroupidlist
	ReplicationGroupIdList []string `json:"ReplicationGroupIdList,omitempty"`

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
func (r *UserGroup_ReplicationGroupIdList) AWSCloudFormationType() string {
	return "AWS::ElastiCache::UserGroup.ReplicationGroupIdList"
}
