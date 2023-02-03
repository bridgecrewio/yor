package datasync

import (
	"github.com/awslabs/goformation/v4/cloudformation/policies"
)

// Task_Options AWS CloudFormation Resource (AWS::DataSync::Task.Options)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-datasync-task-options.html
type Task_Options struct {

	// Atime AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-datasync-task-options.html#cfn-datasync-task-options-atime
	Atime string `json:"Atime,omitempty"`

	// BytesPerSecond AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-datasync-task-options.html#cfn-datasync-task-options-bytespersecond
	BytesPerSecond int `json:"BytesPerSecond,omitempty"`

	// Gid AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-datasync-task-options.html#cfn-datasync-task-options-gid
	Gid string `json:"Gid,omitempty"`

	// LogLevel AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-datasync-task-options.html#cfn-datasync-task-options-loglevel
	LogLevel string `json:"LogLevel,omitempty"`

	// Mtime AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-datasync-task-options.html#cfn-datasync-task-options-mtime
	Mtime string `json:"Mtime,omitempty"`

	// OverwriteMode AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-datasync-task-options.html#cfn-datasync-task-options-overwritemode
	OverwriteMode string `json:"OverwriteMode,omitempty"`

	// PosixPermissions AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-datasync-task-options.html#cfn-datasync-task-options-posixpermissions
	PosixPermissions string `json:"PosixPermissions,omitempty"`

	// PreserveDeletedFiles AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-datasync-task-options.html#cfn-datasync-task-options-preservedeletedfiles
	PreserveDeletedFiles string `json:"PreserveDeletedFiles,omitempty"`

	// PreserveDevices AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-datasync-task-options.html#cfn-datasync-task-options-preservedevices
	PreserveDevices string `json:"PreserveDevices,omitempty"`

	// TaskQueueing AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-datasync-task-options.html#cfn-datasync-task-options-taskqueueing
	TaskQueueing string `json:"TaskQueueing,omitempty"`

	// TransferMode AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-datasync-task-options.html#cfn-datasync-task-options-transfermode
	TransferMode string `json:"TransferMode,omitempty"`

	// Uid AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-datasync-task-options.html#cfn-datasync-task-options-uid
	Uid string `json:"Uid,omitempty"`

	// VerifyMode AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-datasync-task-options.html#cfn-datasync-task-options-verifymode
	VerifyMode string `json:"VerifyMode,omitempty"`

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
func (r *Task_Options) AWSCloudFormationType() string {
	return "AWS::DataSync::Task.Options"
}
