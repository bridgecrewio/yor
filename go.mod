module bridgecrewio/yor

go 1.13

require (
	github.com/go-git/go-git/v5 v5.2.0
	github.com/google/uuid v1.2.0
	github.com/hashicorp/go-hclog v0.9.2
	github.com/hashicorp/hcl/v2 v2.8.2
	github.com/hashicorp/terraform v0.14.0
	github.com/minamijoyo/tfschema v0.6.0
	github.com/mitchellh/cli v1.1.0
	github.com/olekukonko/tablewriter v0.0.5
	github.com/stretchr/testify v1.5.1
)

replace github.com/hashicorp/terraform v0.14.0 => github.com/hashicorp/terraform v0.12.29
