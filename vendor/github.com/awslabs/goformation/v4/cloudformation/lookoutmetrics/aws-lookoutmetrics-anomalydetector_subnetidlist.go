package lookoutmetrics

import (
	"github.com/awslabs/goformation/v4/cloudformation/policies"
)

// AnomalyDetector_SubnetIdList AWS CloudFormation Resource (AWS::LookoutMetrics::AnomalyDetector.SubnetIdList)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-lookoutmetrics-anomalydetector-subnetidlist.html
type AnomalyDetector_SubnetIdList struct {

	// SubnetIdList AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-lookoutmetrics-anomalydetector-subnetidlist.html#cfn-lookoutmetrics-anomalydetector-subnetidlist-subnetidlist
	SubnetIdList []string `json:"SubnetIdList,omitempty"`

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
func (r *AnomalyDetector_SubnetIdList) AWSCloudFormationType() string {
	return "AWS::LookoutMetrics::AnomalyDetector.SubnetIdList"
}
