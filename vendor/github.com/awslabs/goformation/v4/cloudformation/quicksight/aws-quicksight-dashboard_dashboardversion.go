package quicksight

import (
	"github.com/awslabs/goformation/v4/cloudformation/policies"
)

// Dashboard_DashboardVersion AWS CloudFormation Resource (AWS::QuickSight::Dashboard.DashboardVersion)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-quicksight-dashboard-dashboardversion.html
type Dashboard_DashboardVersion struct {

	// Arn AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-quicksight-dashboard-dashboardversion.html#cfn-quicksight-dashboard-dashboardversion-arn
	Arn string `json:"Arn,omitempty"`

	// CreatedTime AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-quicksight-dashboard-dashboardversion.html#cfn-quicksight-dashboard-dashboardversion-createdtime
	CreatedTime string `json:"CreatedTime,omitempty"`

	// DataSetArns AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-quicksight-dashboard-dashboardversion.html#cfn-quicksight-dashboard-dashboardversion-datasetarns
	DataSetArns []string `json:"DataSetArns,omitempty"`

	// Description AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-quicksight-dashboard-dashboardversion.html#cfn-quicksight-dashboard-dashboardversion-description
	Description string `json:"Description,omitempty"`

	// Errors AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-quicksight-dashboard-dashboardversion.html#cfn-quicksight-dashboard-dashboardversion-errors
	Errors []Dashboard_DashboardError `json:"Errors,omitempty"`

	// Sheets AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-quicksight-dashboard-dashboardversion.html#cfn-quicksight-dashboard-dashboardversion-sheets
	Sheets []Dashboard_Sheet `json:"Sheets,omitempty"`

	// SourceEntityArn AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-quicksight-dashboard-dashboardversion.html#cfn-quicksight-dashboard-dashboardversion-sourceentityarn
	SourceEntityArn string `json:"SourceEntityArn,omitempty"`

	// Status AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-quicksight-dashboard-dashboardversion.html#cfn-quicksight-dashboard-dashboardversion-status
	Status string `json:"Status,omitempty"`

	// ThemeArn AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-quicksight-dashboard-dashboardversion.html#cfn-quicksight-dashboard-dashboardversion-themearn
	ThemeArn string `json:"ThemeArn,omitempty"`

	// VersionNumber AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-quicksight-dashboard-dashboardversion.html#cfn-quicksight-dashboard-dashboardversion-versionnumber
	VersionNumber float64 `json:"VersionNumber,omitempty"`

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
func (r *Dashboard_DashboardVersion) AWSCloudFormationType() string {
	return "AWS::QuickSight::Dashboard.DashboardVersion"
}
