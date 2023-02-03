package auditmanager

import (
	"github.com/awslabs/goformation/v4/cloudformation/policies"
)

// Assessment_AWSAccounts AWS CloudFormation Resource (AWS::AuditManager::Assessment.AWSAccounts)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-auditmanager-assessment-awsaccounts.html
type Assessment_AWSAccounts struct {

	// AWSAccounts AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-auditmanager-assessment-awsaccounts.html#cfn-auditmanager-assessment-awsaccounts-awsaccounts
	AWSAccounts []Assessment_AWSAccount `json:"AWSAccounts,omitempty"`

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
func (r *Assessment_AWSAccounts) AWSCloudFormationType() string {
	return "AWS::AuditManager::Assessment.AWSAccounts"
}
