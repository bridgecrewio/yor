package lookoutmetrics

import (
	"github.com/awslabs/goformation/v4/cloudformation/policies"
)

// AnomalyDetector_SecurityGroupIdList AWS CloudFormation Resource (AWS::LookoutMetrics::AnomalyDetector.SecurityGroupIdList)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-lookoutmetrics-anomalydetector-securitygroupidlist.html
type AnomalyDetector_SecurityGroupIdList struct {

	// SecurityGroupIdList AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-lookoutmetrics-anomalydetector-securitygroupidlist.html#cfn-lookoutmetrics-anomalydetector-securitygroupidlist-securitygroupidlist
	SecurityGroupIdList []string `json:"SecurityGroupIdList,omitempty"`

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
func (r *AnomalyDetector_SecurityGroupIdList) AWSCloudFormationType() string {
	return "AWS::LookoutMetrics::AnomalyDetector.SecurityGroupIdList"
}
