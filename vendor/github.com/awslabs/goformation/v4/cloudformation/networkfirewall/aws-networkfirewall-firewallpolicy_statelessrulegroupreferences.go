package networkfirewall

import (
	"github.com/awslabs/goformation/v4/cloudformation/policies"
)

// FirewallPolicy_StatelessRuleGroupReferences AWS CloudFormation Resource (AWS::NetworkFirewall::FirewallPolicy.StatelessRuleGroupReferences)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-networkfirewall-firewallpolicy-statelessrulegroupreferences.html
type FirewallPolicy_StatelessRuleGroupReferences struct {

	// StatelessRuleGroupReferences AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-networkfirewall-firewallpolicy-statelessrulegroupreferences.html#cfn-networkfirewall-firewallpolicy-statelessrulegroupreferences-statelessrulegroupreferences
	StatelessRuleGroupReferences []FirewallPolicy_StatelessRuleGroupReference `json:"StatelessRuleGroupReferences,omitempty"`

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
func (r *FirewallPolicy_StatelessRuleGroupReferences) AWSCloudFormationType() string {
	return "AWS::NetworkFirewall::FirewallPolicy.StatelessRuleGroupReferences"
}
