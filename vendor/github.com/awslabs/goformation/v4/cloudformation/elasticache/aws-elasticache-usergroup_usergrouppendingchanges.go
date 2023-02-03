package elasticache

import (
	"github.com/awslabs/goformation/v4/cloudformation/policies"
)

// UserGroup_UserGroupPendingChanges AWS CloudFormation Resource (AWS::ElastiCache::UserGroup.UserGroupPendingChanges)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-elasticache-usergroup-usergrouppendingchanges.html
type UserGroup_UserGroupPendingChanges struct {

	// UserIdsToAdd AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-elasticache-usergroup-usergrouppendingchanges.html#cfn-elasticache-usergroup-usergrouppendingchanges-useridstoadd
	UserIdsToAdd []string `json:"UserIdsToAdd,omitempty"`

	// UserIdsToRemove AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-elasticache-usergroup-usergrouppendingchanges.html#cfn-elasticache-usergroup-usergrouppendingchanges-useridstoremove
	UserIdsToRemove []string `json:"UserIdsToRemove,omitempty"`

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
func (r *UserGroup_UserGroupPendingChanges) AWSCloudFormationType() string {
	return "AWS::ElastiCache::UserGroup.UserGroupPendingChanges"
}
