package kendra

import (
	"github.com/awslabs/goformation/v4/cloudformation/policies"
)

// DataSource_ConfluencePageFieldMappingsList AWS CloudFormation Resource (AWS::Kendra::DataSource.ConfluencePageFieldMappingsList)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kendra-datasource-confluencepagefieldmappingslist.html
type DataSource_ConfluencePageFieldMappingsList struct {

	// ConfluencePageFieldMappingsList AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kendra-datasource-confluencepagefieldmappingslist.html#cfn-kendra-datasource-confluencepagefieldmappingslist-confluencepagefieldmappingslist
	ConfluencePageFieldMappingsList []DataSource_ConfluencePageToIndexFieldMapping `json:"ConfluencePageFieldMappingsList,omitempty"`

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
func (r *DataSource_ConfluencePageFieldMappingsList) AWSCloudFormationType() string {
	return "AWS::Kendra::DataSource.ConfluencePageFieldMappingsList"
}
