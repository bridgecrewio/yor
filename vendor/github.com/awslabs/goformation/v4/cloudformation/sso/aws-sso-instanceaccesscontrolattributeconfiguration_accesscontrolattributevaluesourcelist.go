package sso

import (
	"github.com/awslabs/goformation/v4/cloudformation/policies"
)

// InstanceAccessControlAttributeConfiguration_AccessControlAttributeValueSourceList AWS CloudFormation Resource (AWS::SSO::InstanceAccessControlAttributeConfiguration.AccessControlAttributeValueSourceList)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-sso-instanceaccesscontrolattributeconfiguration-accesscontrolattributevaluesourcelist.html
type InstanceAccessControlAttributeConfiguration_AccessControlAttributeValueSourceList struct {

	// AccessControlAttributeValueSourceList AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-sso-instanceaccesscontrolattributeconfiguration-accesscontrolattributevaluesourcelist.html#cfn-sso-instanceaccesscontrolattributeconfiguration-accesscontrolattributevaluesourcelist-accesscontrolattributevaluesourcelist
	AccessControlAttributeValueSourceList []string `json:"AccessControlAttributeValueSourceList,omitempty"`

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
func (r *InstanceAccessControlAttributeConfiguration_AccessControlAttributeValueSourceList) AWSCloudFormationType() string {
	return "AWS::SSO::InstanceAccessControlAttributeConfiguration.AccessControlAttributeValueSourceList"
}
