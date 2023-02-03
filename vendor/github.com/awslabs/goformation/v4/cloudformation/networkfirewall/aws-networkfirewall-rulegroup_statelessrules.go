package networkfirewall

import (
	"github.com/awslabs/goformation/v4/cloudformation/policies"
)

// RuleGroup_StatelessRules AWS CloudFormation Resource (AWS::NetworkFirewall::RuleGroup.StatelessRules)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-networkfirewall-rulegroup-statelessrules.html
type RuleGroup_StatelessRules struct {

	// StatelessRules AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-networkfirewall-rulegroup-statelessrules.html#cfn-networkfirewall-rulegroup-statelessrules-statelessrules
	StatelessRules []RuleGroup_StatelessRule `json:"StatelessRules,omitempty"`

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
func (r *RuleGroup_StatelessRules) AWSCloudFormationType() string {
	return "AWS::NetworkFirewall::RuleGroup.StatelessRules"
}
