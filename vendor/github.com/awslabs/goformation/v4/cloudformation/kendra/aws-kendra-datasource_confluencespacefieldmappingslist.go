package kendra

import (
	"github.com/awslabs/goformation/v4/cloudformation/policies"
)

// DataSource_ConfluenceSpaceFieldMappingsList AWS CloudFormation Resource (AWS::Kendra::DataSource.ConfluenceSpaceFieldMappingsList)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kendra-datasource-confluencespacefieldmappingslist.html
type DataSource_ConfluenceSpaceFieldMappingsList struct {

	// ConfluenceSpaceFieldMappingsList AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kendra-datasource-confluencespacefieldmappingslist.html#cfn-kendra-datasource-confluencespacefieldmappingslist-confluencespacefieldmappingslist
	ConfluenceSpaceFieldMappingsList []DataSource_ConfluenceSpaceToIndexFieldMapping `json:"ConfluenceSpaceFieldMappingsList,omitempty"`

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
func (r *DataSource_ConfluenceSpaceFieldMappingsList) AWSCloudFormationType() string {
	return "AWS::Kendra::DataSource.ConfluenceSpaceFieldMappingsList"
}
