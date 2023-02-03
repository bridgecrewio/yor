package kendra

import (
	"github.com/awslabs/goformation/v4/cloudformation/policies"
)

// DataSource_ConfluenceBlogFieldMappingsList AWS CloudFormation Resource (AWS::Kendra::DataSource.ConfluenceBlogFieldMappingsList)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kendra-datasource-confluenceblogfieldmappingslist.html
type DataSource_ConfluenceBlogFieldMappingsList struct {

	// ConfluenceBlogFieldMappingsList AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kendra-datasource-confluenceblogfieldmappingslist.html#cfn-kendra-datasource-confluenceblogfieldmappingslist-confluenceblogfieldmappingslist
	ConfluenceBlogFieldMappingsList []DataSource_ConfluenceBlogToIndexFieldMapping `json:"ConfluenceBlogFieldMappingsList,omitempty"`

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
func (r *DataSource_ConfluenceBlogFieldMappingsList) AWSCloudFormationType() string {
	return "AWS::Kendra::DataSource.ConfluenceBlogFieldMappingsList"
}
