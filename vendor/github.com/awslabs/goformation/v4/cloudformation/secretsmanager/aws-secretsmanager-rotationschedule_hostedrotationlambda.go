package secretsmanager

import (
	"github.com/awslabs/goformation/v4/cloudformation/policies"
)

// RotationSchedule_HostedRotationLambda AWS CloudFormation Resource (AWS::SecretsManager::RotationSchedule.HostedRotationLambda)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-secretsmanager-rotationschedule-hostedrotationlambda.html
type RotationSchedule_HostedRotationLambda struct {

	// KmsKeyArn AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-secretsmanager-rotationschedule-hostedrotationlambda.html#cfn-secretsmanager-rotationschedule-hostedrotationlambda-kmskeyarn
	KmsKeyArn string `json:"KmsKeyArn,omitempty"`

	// MasterSecretArn AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-secretsmanager-rotationschedule-hostedrotationlambda.html#cfn-secretsmanager-rotationschedule-hostedrotationlambda-mastersecretarn
	MasterSecretArn string `json:"MasterSecretArn,omitempty"`

	// MasterSecretKmsKeyArn AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-secretsmanager-rotationschedule-hostedrotationlambda.html#cfn-secretsmanager-rotationschedule-hostedrotationlambda-mastersecretkmskeyarn
	MasterSecretKmsKeyArn string `json:"MasterSecretKmsKeyArn,omitempty"`

	// RotationLambdaName AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-secretsmanager-rotationschedule-hostedrotationlambda.html#cfn-secretsmanager-rotationschedule-hostedrotationlambda-rotationlambdaname
	RotationLambdaName string `json:"RotationLambdaName,omitempty"`

	// RotationType AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-secretsmanager-rotationschedule-hostedrotationlambda.html#cfn-secretsmanager-rotationschedule-hostedrotationlambda-rotationtype
	RotationType string `json:"RotationType,omitempty"`

	// VpcSecurityGroupIds AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-secretsmanager-rotationschedule-hostedrotationlambda.html#cfn-secretsmanager-rotationschedule-hostedrotationlambda-vpcsecuritygroupids
	VpcSecurityGroupIds string `json:"VpcSecurityGroupIds,omitempty"`

	// VpcSubnetIds AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-secretsmanager-rotationschedule-hostedrotationlambda.html#cfn-secretsmanager-rotationschedule-hostedrotationlambda-vpcsubnetids
	VpcSubnetIds string `json:"VpcSubnetIds,omitempty"`

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
func (r *RotationSchedule_HostedRotationLambda) AWSCloudFormationType() string {
	return "AWS::SecretsManager::RotationSchedule.HostedRotationLambda"
}
