package kendra

import (
	"github.com/awslabs/goformation/v4/cloudformation/policies"
)

// Index_Relevance AWS CloudFormation Resource (AWS::Kendra::Index.Relevance)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kendra-index-relevance.html
type Index_Relevance struct {

	// Duration AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kendra-index-relevance.html#cfn-kendra-index-relevance-duration
	Duration string `json:"Duration,omitempty"`

	// Freshness AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kendra-index-relevance.html#cfn-kendra-index-relevance-freshness
	Freshness bool `json:"Freshness,omitempty"`

	// Importance AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kendra-index-relevance.html#cfn-kendra-index-relevance-importance
	Importance int `json:"Importance,omitempty"`

	// RankOrder AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kendra-index-relevance.html#cfn-kendra-index-relevance-rankorder
	RankOrder string `json:"RankOrder,omitempty"`

	// ValueImportanceItems AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kendra-index-relevance.html#cfn-kendra-index-relevance-valueimportanceitems
	ValueImportanceItems []Index_ValueImportanceItem `json:"ValueImportanceItems,omitempty"`

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
func (r *Index_Relevance) AWSCloudFormationType() string {
	return "AWS::Kendra::Index.Relevance"
}
