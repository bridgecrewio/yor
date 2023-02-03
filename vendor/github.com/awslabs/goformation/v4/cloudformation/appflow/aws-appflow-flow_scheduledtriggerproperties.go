package appflow

import (
	"github.com/awslabs/goformation/v4/cloudformation/policies"
)

// Flow_ScheduledTriggerProperties AWS CloudFormation Resource (AWS::AppFlow::Flow.ScheduledTriggerProperties)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appflow-flow-scheduledtriggerproperties.html
type Flow_ScheduledTriggerProperties struct {

	// DataPullMode AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appflow-flow-scheduledtriggerproperties.html#cfn-appflow-flow-scheduledtriggerproperties-datapullmode
	DataPullMode string `json:"DataPullMode,omitempty"`

	// ScheduleEndTime AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appflow-flow-scheduledtriggerproperties.html#cfn-appflow-flow-scheduledtriggerproperties-scheduleendtime
	ScheduleEndTime float64 `json:"ScheduleEndTime,omitempty"`

	// ScheduleExpression AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appflow-flow-scheduledtriggerproperties.html#cfn-appflow-flow-scheduledtriggerproperties-scheduleexpression
	ScheduleExpression string `json:"ScheduleExpression,omitempty"`

	// ScheduleStartTime AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appflow-flow-scheduledtriggerproperties.html#cfn-appflow-flow-scheduledtriggerproperties-schedulestarttime
	ScheduleStartTime float64 `json:"ScheduleStartTime,omitempty"`

	// TimeZone AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appflow-flow-scheduledtriggerproperties.html#cfn-appflow-flow-scheduledtriggerproperties-timezone
	TimeZone string `json:"TimeZone,omitempty"`

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
func (r *Flow_ScheduledTriggerProperties) AWSCloudFormationType() string {
	return "AWS::AppFlow::Flow.ScheduledTriggerProperties"
}
