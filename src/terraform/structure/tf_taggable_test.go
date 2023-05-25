package structure

import (
	aws_v2 "github.com/lonegunmanb/terraform-aws-schema/v2/generated"
	gcp_v2 "github.com/lonegunmanb/terraform-google-schema/v2/generated"
	"testing"
)

func TestGetType(t *testing.T) {
	//previousResources := make(map[string]struct{}, 0)
	for _, schema := range aws_v2.Resources {
		tags, ok := schema.Block.Attributes["tags"]
		if !ok {
			continue
		}
		atype := tags.AttributeType
		s := atype.GoString()
		println(s)
	}
}

func TestGetGoogleType(t *testing.T) {
	//previousResources := make(map[string]struct{}, 0)
	for _, schema := range gcp_v2.Resources {
		tags, ok := schema.Block.Attributes["tags"]
		if !ok {
			continue
		}
		atype := tags.AttributeType
		s := atype.GoString()
		println(s)
	}
}
