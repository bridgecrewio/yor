package networkfirewall

import (
	"github.com/awslabs/goformation/v4/cloudformation/policies"
)

// LoggingConfiguration_LogDestinationConfigs AWS CloudFormation Resource (AWS::NetworkFirewall::LoggingConfiguration.LogDestinationConfigs)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-networkfirewall-loggingconfiguration-logdestinationconfigs.html
type LoggingConfiguration_LogDestinationConfigs struct {

	// LogDestinationConfigs AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-networkfirewall-loggingconfiguration-logdestinationconfigs.html#cfn-networkfirewall-loggingconfiguration-logdestinationconfigs-logdestinationconfigs
	LogDestinationConfigs []LoggingConfiguration_LogDestinationConfig `json:"LogDestinationConfigs,omitempty"`

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
func (r *LoggingConfiguration_LogDestinationConfigs) AWSCloudFormationType() string {
	return "AWS::NetworkFirewall::LoggingConfiguration.LogDestinationConfigs"
}
