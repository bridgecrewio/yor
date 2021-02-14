module bridgecrewio/yor

go 1.13

require (
	github.com/go-git/go-git/v5 v5.2.0
	github.com/google/uuid v1.2.0
	github.com/hashicorp/go-hclog v0.9.2
	github.com/hashicorp/hcl/v2 v2.8.2
	github.com/hashicorp/terraform v0.14.0
	github.com/jpoles1/gopherbadger v2.4.0+incompatible // indirect
	github.com/minamijoyo/tfschema v0.6.0
	github.com/mitchellh/cli v1.1.0
	github.com/olekukonko/tablewriter v0.0.5
	github.com/stretchr/testify v1.5.1
	golang.org/x/net v0.0.0-20201021035429-f5854403a974 // indirect
	golang.org/x/sys v0.0.0-20210119212857-b64e53b001e4 // indirect
)

replace github.com/hashicorp/terraform v0.14.0 => github.com/hashicorp/terraform v0.12.29