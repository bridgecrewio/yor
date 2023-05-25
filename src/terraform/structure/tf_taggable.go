package structure

import (
	"github.com/ahmetb/go-linq/v3"
	tfjson "github.com/hashicorp/terraform-json"
	aws_v2 "github.com/lonegunmanb/terraform-aws-schema/v2/generated"
	aws_v3 "github.com/lonegunmanb/terraform-aws-schema/v3/generated"
	aws_v4 "github.com/lonegunmanb/terraform-aws-schema/v4/generated"
	azure_v2 "github.com/lonegunmanb/terraform-azurerm-schema/v2/generated"
	azure_v3 "github.com/lonegunmanb/terraform-azurerm-schema/v3/generated"
	google_v2 "github.com/lonegunmanb/terraform-google-schema/v2/generated"
	google_v3 "github.com/lonegunmanb/terraform-google-schema/v3/generated"
	google_v4 "github.com/lonegunmanb/terraform-google-schema/v4/generated"
)

var TfTaggableResourceTypes []string

func init() {
	linq.From(previousTaggableTypes(aws_v4.Resources, isTaggableType, aws_v2.Resources, aws_v3.Resources)).
		Concat(linq.From(previousTaggableTypes(azure_v3.Resources, isTaggableType, azure_v2.Resources))).
		Concat(linq.From(previousTaggableTypes(google_v4.Resources, isGoogleTaggableType, google_v2.Resources, google_v3.Resources))).
		Concat(linq.From(aws_v4.Resources).Where(isTaggableType)).
		Concat(linq.From(azure_v3.Resources).Where(isTaggableType)).
		Concat(linq.From(google_v4.Resources).Where(isGoogleTaggableType))
}

func previousTaggableTypes(latestTypes map[string]*tfjson.Schema, taggable func(interface{}) bool, previousTypes ...map[string]*tfjson.Schema) []string {
	var r []string
	for _, types := range previousTypes {
		for n, schema := range types {
			if taggable(schema) && !tagsHasBeenRemoved(n, latestTypes) {
				r = append(r, n)
			}
		}
	}
	return r
}

func tagsHasBeenRemoved(name string, latestTypes map[string]*tfjson.Schema) bool {
	s, ok := latestTypes[name]
	if !ok {
		return false
	}
	_, stillExist := s.Block.Attributes["tags"]
	return !stillExist
}

func isTaggableType(s interface{}) bool {
	schema := s.(*tfjson.Schema)
	a, ok := schema.Block.Attributes["tags"]
	return ok && a.AttributeType.GoString() == "cty.Map(cty.String)"
}

func isGoogleTaggableType(s interface{}) bool {
	schema := s.(*tfjson.Schema)
	a, ok := schema.Block.Attributes["tags"]
	return ok && a.AttributeType.GoString() == "cty.Set(cty.String)"
}
