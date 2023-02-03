package budgets

import (
	"github.com/awslabs/goformation/v4/cloudformation/policies"
)

// BudgetsAction_IamActionDefinition AWS CloudFormation Resource (AWS::Budgets::BudgetsAction.IamActionDefinition)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-budgets-budgetsaction-iamactiondefinition.html
type BudgetsAction_IamActionDefinition struct {

	// Groups AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-budgets-budgetsaction-iamactiondefinition.html#cfn-budgets-budgetsaction-iamactiondefinition-groups
	Groups []string `json:"Groups,omitempty"`

	// PolicyArn AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-budgets-budgetsaction-iamactiondefinition.html#cfn-budgets-budgetsaction-iamactiondefinition-policyarn
	PolicyArn string `json:"PolicyArn,omitempty"`

	// Roles AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-budgets-budgetsaction-iamactiondefinition.html#cfn-budgets-budgetsaction-iamactiondefinition-roles
	Roles []string `json:"Roles,omitempty"`

	// Users AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-budgets-budgetsaction-iamactiondefinition.html#cfn-budgets-budgetsaction-iamactiondefinition-users
	Users []string `json:"Users,omitempty"`

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
func (r *BudgetsAction_IamActionDefinition) AWSCloudFormationType() string {
	return "AWS::Budgets::BudgetsAction.IamActionDefinition"
}
