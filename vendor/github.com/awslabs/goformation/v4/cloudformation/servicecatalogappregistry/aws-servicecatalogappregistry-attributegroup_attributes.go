package servicecatalogappregistry

import (
	"github.com/awslabs/goformation/v4/cloudformation/policies"
)

// AttributeGroup_Attributes AWS CloudFormation Resource (AWS::ServiceCatalogAppRegistry::AttributeGroup.Attributes)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-servicecatalogappregistry-attributegroup-attributes.html
type AttributeGroup_Attributes struct {

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
func (r *AttributeGroup_Attributes) AWSCloudFormationType() string {
	return "AWS::ServiceCatalogAppRegistry::AttributeGroup.Attributes"
}
