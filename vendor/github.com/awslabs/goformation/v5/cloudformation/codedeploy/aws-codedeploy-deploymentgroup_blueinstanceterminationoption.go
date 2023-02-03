package codedeploy

import (
	"github.com/awslabs/goformation/v5/cloudformation/policies"
)

// DeploymentGroup_BlueInstanceTerminationOption AWS CloudFormation Resource (AWS::CodeDeploy::DeploymentGroup.BlueInstanceTerminationOption)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-codedeploy-deploymentgroup-blueinstanceterminationoption.html
type DeploymentGroup_BlueInstanceTerminationOption struct {

	// Action AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-codedeploy-deploymentgroup-blueinstanceterminationoption.html#cfn-codedeploy-deploymentgroup-bluegreendeploymentconfiguration-blueinstanceterminationoption-action
	Action string `json:"Action,omitempty"`

	// TerminationWaitTimeInMinutes AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-codedeploy-deploymentgroup-blueinstanceterminationoption.html#cfn-codedeploy-deploymentgroup-bluegreendeploymentconfiguration-blueinstanceterminationoption-terminationwaittimeinminutes
	TerminationWaitTimeInMinutes int `json:"TerminationWaitTimeInMinutes,omitempty"`

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
func (r *DeploymentGroup_BlueInstanceTerminationOption) AWSCloudFormationType() string {
	return "AWS::CodeDeploy::DeploymentGroup.BlueInstanceTerminationOption"
}
