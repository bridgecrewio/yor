package appmesh

import (
	"github.com/awslabs/goformation/v4/cloudformation/policies"
)

// VirtualGateway_VirtualGatewayHealthCheckPolicy AWS CloudFormation Resource (AWS::AppMesh::VirtualGateway.VirtualGatewayHealthCheckPolicy)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appmesh-virtualgateway-virtualgatewayhealthcheckpolicy.html
type VirtualGateway_VirtualGatewayHealthCheckPolicy struct {

	// HealthyThreshold AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appmesh-virtualgateway-virtualgatewayhealthcheckpolicy.html#cfn-appmesh-virtualgateway-virtualgatewayhealthcheckpolicy-healthythreshold
	HealthyThreshold int `json:"HealthyThreshold"`

	// IntervalMillis AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appmesh-virtualgateway-virtualgatewayhealthcheckpolicy.html#cfn-appmesh-virtualgateway-virtualgatewayhealthcheckpolicy-intervalmillis
	IntervalMillis int `json:"IntervalMillis"`

	// Path AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appmesh-virtualgateway-virtualgatewayhealthcheckpolicy.html#cfn-appmesh-virtualgateway-virtualgatewayhealthcheckpolicy-path
	Path string `json:"Path,omitempty"`

	// Port AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appmesh-virtualgateway-virtualgatewayhealthcheckpolicy.html#cfn-appmesh-virtualgateway-virtualgatewayhealthcheckpolicy-port
	Port int `json:"Port,omitempty"`

	// Protocol AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appmesh-virtualgateway-virtualgatewayhealthcheckpolicy.html#cfn-appmesh-virtualgateway-virtualgatewayhealthcheckpolicy-protocol
	Protocol string `json:"Protocol,omitempty"`

	// TimeoutMillis AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appmesh-virtualgateway-virtualgatewayhealthcheckpolicy.html#cfn-appmesh-virtualgateway-virtualgatewayhealthcheckpolicy-timeoutmillis
	TimeoutMillis int `json:"TimeoutMillis"`

	// UnhealthyThreshold AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appmesh-virtualgateway-virtualgatewayhealthcheckpolicy.html#cfn-appmesh-virtualgateway-virtualgatewayhealthcheckpolicy-unhealthythreshold
	UnhealthyThreshold int `json:"UnhealthyThreshold"`

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
func (r *VirtualGateway_VirtualGatewayHealthCheckPolicy) AWSCloudFormationType() string {
	return "AWS::AppMesh::VirtualGateway.VirtualGatewayHealthCheckPolicy"
}
