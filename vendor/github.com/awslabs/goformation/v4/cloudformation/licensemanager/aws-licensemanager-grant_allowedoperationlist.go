package licensemanager

import (
	"github.com/awslabs/goformation/v4/cloudformation/policies"
)

// Grant_AllowedOperationList AWS CloudFormation Resource (AWS::LicenseManager::Grant.AllowedOperationList)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-licensemanager-grant-allowedoperationlist.html
type Grant_AllowedOperationList struct {

	// AllowedOperationList AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-licensemanager-grant-allowedoperationlist.html#cfn-licensemanager-grant-allowedoperationlist-allowedoperationlist
	AllowedOperationList []string `json:"AllowedOperationList,omitempty"`

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
func (r *Grant_AllowedOperationList) AWSCloudFormationType() string {
	return "AWS::LicenseManager::Grant.AllowedOperationList"
}
