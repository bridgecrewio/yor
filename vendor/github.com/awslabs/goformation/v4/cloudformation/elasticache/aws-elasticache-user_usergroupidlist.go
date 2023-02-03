package elasticache

import (
	"github.com/awslabs/goformation/v4/cloudformation/policies"
)

// User_UserGroupIdList AWS CloudFormation Resource (AWS::ElastiCache::User.UserGroupIdList)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-elasticache-user-usergroupidlist.html
type User_UserGroupIdList struct {

	// UserGroupIdList AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-elasticache-user-usergroupidlist.html#cfn-elasticache-user-usergroupidlist-usergroupidlist
	UserGroupIdList []string `json:"UserGroupIdList,omitempty"`

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
func (r *User_UserGroupIdList) AWSCloudFormationType() string {
	return "AWS::ElastiCache::User.UserGroupIdList"
}
