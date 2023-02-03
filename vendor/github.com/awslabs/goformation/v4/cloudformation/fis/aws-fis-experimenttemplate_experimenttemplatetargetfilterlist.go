package fis

import (
	"github.com/awslabs/goformation/v4/cloudformation/policies"
)

// ExperimentTemplate_ExperimentTemplateTargetFilterList AWS CloudFormation Resource (AWS::FIS::ExperimentTemplate.ExperimentTemplateTargetFilterList)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-fis-experimenttemplate-experimenttemplatetargetfilterlist.html
type ExperimentTemplate_ExperimentTemplateTargetFilterList struct {

	// ExperimentTemplateTargetFilterList AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-fis-experimenttemplate-experimenttemplatetargetfilterlist.html#cfn-fis-experimenttemplate-experimenttemplatetargetfilterlist-experimenttemplatetargetfilterlist
	ExperimentTemplateTargetFilterList []ExperimentTemplate_ExperimentTemplateTargetFilter `json:"ExperimentTemplateTargetFilterList,omitempty"`

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
func (r *ExperimentTemplate_ExperimentTemplateTargetFilterList) AWSCloudFormationType() string {
	return "AWS::FIS::ExperimentTemplate.ExperimentTemplateTargetFilterList"
}
