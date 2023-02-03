package sagemaker

import (
	"github.com/awslabs/goformation/v4/cloudformation/policies"
)

// ModelQualityJobDefinition_ModelQualityAppSpecification AWS CloudFormation Resource (AWS::SageMaker::ModelQualityJobDefinition.ModelQualityAppSpecification)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-sagemaker-modelqualityjobdefinition-modelqualityappspecification.html
type ModelQualityJobDefinition_ModelQualityAppSpecification struct {

	// ContainerArguments AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-sagemaker-modelqualityjobdefinition-modelqualityappspecification.html#cfn-sagemaker-modelqualityjobdefinition-modelqualityappspecification-containerarguments
	ContainerArguments []string `json:"ContainerArguments,omitempty"`

	// ContainerEntrypoint AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-sagemaker-modelqualityjobdefinition-modelqualityappspecification.html#cfn-sagemaker-modelqualityjobdefinition-modelqualityappspecification-containerentrypoint
	ContainerEntrypoint []string `json:"ContainerEntrypoint,omitempty"`

	// Environment AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-sagemaker-modelqualityjobdefinition-modelqualityappspecification.html#cfn-sagemaker-modelqualityjobdefinition-modelqualityappspecification-environment
	Environment *ModelQualityJobDefinition_Environment `json:"Environment,omitempty"`

	// ImageUri AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-sagemaker-modelqualityjobdefinition-modelqualityappspecification.html#cfn-sagemaker-modelqualityjobdefinition-modelqualityappspecification-imageuri
	ImageUri string `json:"ImageUri,omitempty"`

	// PostAnalyticsProcessorSourceUri AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-sagemaker-modelqualityjobdefinition-modelqualityappspecification.html#cfn-sagemaker-modelqualityjobdefinition-modelqualityappspecification-postanalyticsprocessorsourceuri
	PostAnalyticsProcessorSourceUri string `json:"PostAnalyticsProcessorSourceUri,omitempty"`

	// ProblemType AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-sagemaker-modelqualityjobdefinition-modelqualityappspecification.html#cfn-sagemaker-modelqualityjobdefinition-modelqualityappspecification-problemtype
	ProblemType string `json:"ProblemType,omitempty"`

	// RecordPreprocessorSourceUri AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-sagemaker-modelqualityjobdefinition-modelqualityappspecification.html#cfn-sagemaker-modelqualityjobdefinition-modelqualityappspecification-recordpreprocessorsourceuri
	RecordPreprocessorSourceUri string `json:"RecordPreprocessorSourceUri,omitempty"`

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
func (r *ModelQualityJobDefinition_ModelQualityAppSpecification) AWSCloudFormationType() string {
	return "AWS::SageMaker::ModelQualityJobDefinition.ModelQualityAppSpecification"
}
