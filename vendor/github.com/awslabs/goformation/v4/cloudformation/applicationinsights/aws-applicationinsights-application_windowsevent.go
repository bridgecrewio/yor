package applicationinsights

import (
	"github.com/awslabs/goformation/v4/cloudformation/policies"
)

// Application_WindowsEvent AWS CloudFormation Resource (AWS::ApplicationInsights::Application.WindowsEvent)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-applicationinsights-application-windowsevent.html
type Application_WindowsEvent struct {

	// EventLevels AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-applicationinsights-application-windowsevent.html#cfn-applicationinsights-application-windowsevent-eventlevels
	EventLevels []string `json:"EventLevels,omitempty"`

	// EventName AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-applicationinsights-application-windowsevent.html#cfn-applicationinsights-application-windowsevent-eventname
	EventName string `json:"EventName,omitempty"`

	// LogGroupName AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-applicationinsights-application-windowsevent.html#cfn-applicationinsights-application-windowsevent-loggroupname
	LogGroupName string `json:"LogGroupName,omitempty"`

	// PatternSet AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-applicationinsights-application-windowsevent.html#cfn-applicationinsights-application-windowsevent-patternset
	PatternSet string `json:"PatternSet,omitempty"`

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
func (r *Application_WindowsEvent) AWSCloudFormationType() string {
	return "AWS::ApplicationInsights::Application.WindowsEvent"
}
