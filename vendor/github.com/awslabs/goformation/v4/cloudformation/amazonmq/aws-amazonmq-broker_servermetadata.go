package amazonmq

import (
	"github.com/awslabs/goformation/v4/cloudformation/policies"
)

// Broker_ServerMetadata AWS CloudFormation Resource (AWS::AmazonMQ::Broker.ServerMetadata)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-amazonmq-broker-servermetadata.html
type Broker_ServerMetadata struct {

	// Hosts AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-amazonmq-broker-servermetadata.html#cfn-amazonmq-broker-servermetadata-hosts
	Hosts []string `json:"Hosts,omitempty"`

	// RoleBase AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-amazonmq-broker-servermetadata.html#cfn-amazonmq-broker-servermetadata-rolebase
	RoleBase string `json:"RoleBase,omitempty"`

	// RoleName AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-amazonmq-broker-servermetadata.html#cfn-amazonmq-broker-servermetadata-rolename
	RoleName string `json:"RoleName,omitempty"`

	// RoleSearchMatching AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-amazonmq-broker-servermetadata.html#cfn-amazonmq-broker-servermetadata-rolesearchmatching
	RoleSearchMatching string `json:"RoleSearchMatching,omitempty"`

	// RoleSearchSubtree AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-amazonmq-broker-servermetadata.html#cfn-amazonmq-broker-servermetadata-rolesearchsubtree
	RoleSearchSubtree bool `json:"RoleSearchSubtree,omitempty"`

	// ServiceAccountPassword AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-amazonmq-broker-servermetadata.html#cfn-amazonmq-broker-servermetadata-serviceaccountpassword
	ServiceAccountPassword string `json:"ServiceAccountPassword,omitempty"`

	// ServiceAccountUsername AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-amazonmq-broker-servermetadata.html#cfn-amazonmq-broker-servermetadata-serviceaccountusername
	ServiceAccountUsername string `json:"ServiceAccountUsername,omitempty"`

	// UserBase AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-amazonmq-broker-servermetadata.html#cfn-amazonmq-broker-servermetadata-userbase
	UserBase string `json:"UserBase,omitempty"`

	// UserRoleName AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-amazonmq-broker-servermetadata.html#cfn-amazonmq-broker-servermetadata-userrolename
	UserRoleName string `json:"UserRoleName,omitempty"`

	// UserSearchMatching AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-amazonmq-broker-servermetadata.html#cfn-amazonmq-broker-servermetadata-usersearchmatching
	UserSearchMatching string `json:"UserSearchMatching,omitempty"`

	// UserSearchSubtree AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-amazonmq-broker-servermetadata.html#cfn-amazonmq-broker-servermetadata-usersearchsubtree
	UserSearchSubtree bool `json:"UserSearchSubtree,omitempty"`

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
func (r *Broker_ServerMetadata) AWSCloudFormationType() string {
	return "AWS::AmazonMQ::Broker.ServerMetadata"
}
