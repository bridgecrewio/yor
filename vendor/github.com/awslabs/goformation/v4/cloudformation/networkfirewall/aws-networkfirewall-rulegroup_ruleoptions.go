package networkfirewall

import (
	"github.com/awslabs/goformation/v4/cloudformation/policies"
)

// RuleGroup_RuleOptions AWS CloudFormation Resource (AWS::NetworkFirewall::RuleGroup.RuleOptions)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-networkfirewall-rulegroup-ruleoptions.html
type RuleGroup_RuleOptions struct {

	// RuleOptions AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-networkfirewall-rulegroup-ruleoptions.html#cfn-networkfirewall-rulegroup-ruleoptions-ruleoptions
	RuleOptions []RuleGroup_RuleOption `json:"RuleOptions,omitempty"`

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
func (r *RuleGroup_RuleOptions) AWSCloudFormationType() string {
	return "AWS::NetworkFirewall::RuleGroup.RuleOptions"
}
