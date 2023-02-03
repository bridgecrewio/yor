package customerprofiles

import (
	"github.com/awslabs/goformation/v4/cloudformation/policies"
)

// Integration_ScheduledTriggerProperties AWS CloudFormation Resource (AWS::CustomerProfiles::Integration.ScheduledTriggerProperties)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-customerprofiles-integration-scheduledtriggerproperties.html
type Integration_ScheduledTriggerProperties struct {

	// DataPullMode AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-customerprofiles-integration-scheduledtriggerproperties.html#cfn-customerprofiles-integration-scheduledtriggerproperties-datapullmode
	DataPullMode string `json:"DataPullMode,omitempty"`

	// FirstExecutionFrom AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-customerprofiles-integration-scheduledtriggerproperties.html#cfn-customerprofiles-integration-scheduledtriggerproperties-firstexecutionfrom
	FirstExecutionFrom float64 `json:"FirstExecutionFrom,omitempty"`

	// ScheduleEndTime AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-customerprofiles-integration-scheduledtriggerproperties.html#cfn-customerprofiles-integration-scheduledtriggerproperties-scheduleendtime
	ScheduleEndTime float64 `json:"ScheduleEndTime,omitempty"`

	// ScheduleExpression AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-customerprofiles-integration-scheduledtriggerproperties.html#cfn-customerprofiles-integration-scheduledtriggerproperties-scheduleexpression
	ScheduleExpression string `json:"ScheduleExpression,omitempty"`

	// ScheduleOffset AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-customerprofiles-integration-scheduledtriggerproperties.html#cfn-customerprofiles-integration-scheduledtriggerproperties-scheduleoffset
	ScheduleOffset int `json:"ScheduleOffset,omitempty"`

	// ScheduleStartTime AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-customerprofiles-integration-scheduledtriggerproperties.html#cfn-customerprofiles-integration-scheduledtriggerproperties-schedulestarttime
	ScheduleStartTime float64 `json:"ScheduleStartTime,omitempty"`

	// Timezone AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-customerprofiles-integration-scheduledtriggerproperties.html#cfn-customerprofiles-integration-scheduledtriggerproperties-timezone
	Timezone string `json:"Timezone,omitempty"`

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
func (r *Integration_ScheduledTriggerProperties) AWSCloudFormationType() string {
	return "AWS::CustomerProfiles::Integration.ScheduledTriggerProperties"
}
