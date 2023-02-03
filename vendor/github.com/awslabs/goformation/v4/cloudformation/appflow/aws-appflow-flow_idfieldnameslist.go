package appflow

import (
	"github.com/awslabs/goformation/v4/cloudformation/policies"
)

// Flow_IdFieldNamesList AWS CloudFormation Resource (AWS::AppFlow::Flow.IdFieldNamesList)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appflow-flow-idfieldnameslist.html
type Flow_IdFieldNamesList struct {

	// IdFieldNamesList AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appflow-flow-idfieldnameslist.html#cfn-appflow-flow-idfieldnameslist-idfieldnameslist
	IdFieldNamesList []string `json:"IdFieldNamesList,omitempty"`

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
func (r *Flow_IdFieldNamesList) AWSCloudFormationType() string {
	return "AWS::AppFlow::Flow.IdFieldNamesList"
}
