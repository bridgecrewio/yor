package applicationinsights

import (
	"github.com/awslabs/goformation/v4/cloudformation/policies"
)

// Application_Log AWS CloudFormation Resource (AWS::ApplicationInsights::Application.Log)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-applicationinsights-application-log.html
type Application_Log struct {

	// Encoding AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-applicationinsights-application-log.html#cfn-applicationinsights-application-log-encoding
	Encoding string `json:"Encoding,omitempty"`

	// LogGroupName AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-applicationinsights-application-log.html#cfn-applicationinsights-application-log-loggroupname
	LogGroupName string `json:"LogGroupName,omitempty"`

	// LogPath AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-applicationinsights-application-log.html#cfn-applicationinsights-application-log-logpath
	LogPath string `json:"LogPath,omitempty"`

	// LogType AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-applicationinsights-application-log.html#cfn-applicationinsights-application-log-logtype
	LogType string `json:"LogType,omitempty"`

	// PatternSet AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-applicationinsights-application-log.html#cfn-applicationinsights-application-log-patternset
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
func (r *Application_Log) AWSCloudFormationType() string {
	return "AWS::ApplicationInsights::Application.Log"
}
