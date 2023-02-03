package serverless

// Function definition
type Function struct {
	Handler      string                 `json:"handler,omitempty"`
	Name         string                 `json:"name,omitempty"`
	Description  string                 `json:"description,omitempty"`
	MemorySize   int                    `json:"memorySize,omitempty"`
	Runtime      string                 `json:"runtime,omitempty"`
	Timeout      int                    `json:"timeout,omitempty"`
	Role         string                 `json:"role,omitempty"`
	OnError      string                 `json:"onError,omitempty"`
	AwsKmsKeyArn string                 `json:"awsKmsKeyArn,omitempty"`
	Environment  map[string]interface{} `json:"environment,omitempty"`
	Tags         map[string]interface{} `json:"tags,omitempty"`
	Events       []Events               `json:"events,omitempty"`
	Package      Package                `json:"package,omitempty"`
	VPC          VPC                    `json:"vpc,omitempty"`
}

// VPC definition
type VPC struct {
	SecurityGroupIds []string `json:"securityGroupIds,omitempty"`
	SubnetIds        []string `json:"subnetIds,omitempty"`
}
