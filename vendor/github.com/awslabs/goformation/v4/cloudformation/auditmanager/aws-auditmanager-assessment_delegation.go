package auditmanager

import (
	"github.com/awslabs/goformation/v4/cloudformation/policies"
)

// Assessment_Delegation AWS CloudFormation Resource (AWS::AuditManager::Assessment.Delegation)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-auditmanager-assessment-delegation.html
type Assessment_Delegation struct {

	// AssessmentId AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-auditmanager-assessment-delegation.html#cfn-auditmanager-assessment-delegation-assessmentid
	AssessmentId string `json:"AssessmentId,omitempty"`

	// AssessmentName AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-auditmanager-assessment-delegation.html#cfn-auditmanager-assessment-delegation-assessmentname
	AssessmentName string `json:"AssessmentName,omitempty"`

	// Comment AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-auditmanager-assessment-delegation.html#cfn-auditmanager-assessment-delegation-comment
	Comment string `json:"Comment,omitempty"`

	// ControlSetId AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-auditmanager-assessment-delegation.html#cfn-auditmanager-assessment-delegation-controlsetid
	ControlSetId string `json:"ControlSetId,omitempty"`

	// CreatedBy AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-auditmanager-assessment-delegation.html#cfn-auditmanager-assessment-delegation-createdby
	CreatedBy string `json:"CreatedBy,omitempty"`

	// CreationTime AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-auditmanager-assessment-delegation.html#cfn-auditmanager-assessment-delegation-creationtime
	CreationTime float64 `json:"CreationTime,omitempty"`

	// Id AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-auditmanager-assessment-delegation.html#cfn-auditmanager-assessment-delegation-id
	Id string `json:"Id,omitempty"`

	// LastUpdated AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-auditmanager-assessment-delegation.html#cfn-auditmanager-assessment-delegation-lastupdated
	LastUpdated float64 `json:"LastUpdated,omitempty"`

	// RoleArn AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-auditmanager-assessment-delegation.html#cfn-auditmanager-assessment-delegation-rolearn
	RoleArn string `json:"RoleArn,omitempty"`

	// RoleType AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-auditmanager-assessment-delegation.html#cfn-auditmanager-assessment-delegation-roletype
	RoleType string `json:"RoleType,omitempty"`

	// Status AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-auditmanager-assessment-delegation.html#cfn-auditmanager-assessment-delegation-status
	Status string `json:"Status,omitempty"`

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
func (r *Assessment_Delegation) AWSCloudFormationType() string {
	return "AWS::AuditManager::Assessment.Delegation"
}
