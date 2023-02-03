package appmesh

import (
	"github.com/awslabs/goformation/v4/cloudformation/policies"
)

// GatewayRoute_HttpGatewayRouteMatch AWS CloudFormation Resource (AWS::AppMesh::GatewayRoute.HttpGatewayRouteMatch)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appmesh-gatewayroute-httpgatewayroutematch.html
type GatewayRoute_HttpGatewayRouteMatch struct {

	// Prefix AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appmesh-gatewayroute-httpgatewayroutematch.html#cfn-appmesh-gatewayroute-httpgatewayroutematch-prefix
	Prefix string `json:"Prefix,omitempty"`

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
func (r *GatewayRoute_HttpGatewayRouteMatch) AWSCloudFormationType() string {
	return "AWS::AppMesh::GatewayRoute.HttpGatewayRouteMatch"
}
