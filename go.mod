module github.com/bridgecrewio/yor

go 1.13

require (
	github.com/awslabs/goformation/v4 v4.19.1
	github.com/go-git/go-git/v5 v5.2.0
	github.com/google/uuid v1.2.0
	github.com/hashicorp/go-hclog v0.9.2
	github.com/hashicorp/go-version v1.2.0
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/hashicorp/hcl/v2 v2.8.2
	github.com/hashicorp/terraform v0.14.0
	github.com/hashicorp/terraform-config-inspect v0.0.0-20191212124732-c6ae6269b9d7
	github.com/minamijoyo/tfschema v0.6.0
	github.com/mitchellh/cli v1.1.0
	github.com/olekukonko/tablewriter v0.0.5
	github.com/pmezard/go-difflib v1.0.0
	github.com/sanathkr/yaml v1.0.0
	github.com/stretchr/testify v1.6.1
	github.com/thepauleh/goserverless v0.0.0-20210509094513-618ed3018bcd
	github.com/urfave/cli/v2 v2.3.0
	github.com/zclconf/go-cty v1.7.0
	go.opencensus.io v0.22.0
	golang.org/x/sys v0.0.0-20210119212857-b64e53b001e4 // indirect
	gopkg.in/validator.v2 v2.0.0-20200605151824-2b28d334fa05
	gopkg.in/yaml.v2 v2.4.0
)

replace github.com/hashicorp/terraform v0.14.0 => github.com/hashicorp/terraform v0.12.31
