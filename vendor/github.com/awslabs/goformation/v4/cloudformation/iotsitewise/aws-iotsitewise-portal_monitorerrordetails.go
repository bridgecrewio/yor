package iotsitewise

import (
	"github.com/awslabs/goformation/v4/cloudformation/policies"
)

// Portal_MonitorErrorDetails AWS CloudFormation Resource (AWS::IoTSiteWise::Portal.MonitorErrorDetails)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotsitewise-portal-monitorerrordetails.html
type Portal_MonitorErrorDetails struct {

	// code AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotsitewise-portal-monitorerrordetails.html#cfn-iotsitewise-portal-monitorerrordetails-code
	code string `json:"code,omitempty"`

	// message AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotsitewise-portal-monitorerrordetails.html#cfn-iotsitewise-portal-monitorerrordetails-message
	message string `json:"message,omitempty"`

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
func (r *Portal_MonitorErrorDetails) AWSCloudFormationType() string {
	return "AWS::IoTSiteWise::Portal.MonitorErrorDetails"
}
