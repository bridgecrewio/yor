package structure

import (
	"strings"

	"github.com/ahmetb/go-linq/v3"
	tfjson "github.com/hashicorp/terraform-json"
	awsv2 "github.com/lonegunmanb/terraform-aws-schema/v2/generated"
	awsv3 "github.com/lonegunmanb/terraform-aws-schema/v3/generated"
	awsv4 "github.com/lonegunmanb/terraform-aws-schema/v4/generated"
	awsv5 "github.com/lonegunmanb/terraform-aws-schema/v5/generated"
	azurev2 "github.com/lonegunmanb/terraform-azurerm-schema/v2/generated"
	azurev3 "github.com/lonegunmanb/terraform-azurerm-schema/v3/generated"
	googlev2 "github.com/lonegunmanb/terraform-google-schema/v2/generated"
	googlev3 "github.com/lonegunmanb/terraform-google-schema/v3/generated"
	googlev4 "github.com/lonegunmanb/terraform-google-schema/v4/generated"
)

var TfTaggableResourceTypes []string

func init() {
	linq.From(previousTaggableTypes(awsv4.Resources, isTaggableType, awsv2.Resources, awsv3.Resources, awsv4.Resources)).
		Concat(linq.From(previousTaggableTypes(azurev3.Resources, isTaggableType, azurev2.Resources))).
		Concat(linq.From(previousTaggableTypes(googlev4.Resources, isGoogleTaggableType, googlev2.Resources, googlev3.Resources))).
		Concat(linq.From(taggableTypes(awsv5.Resources, isTaggableType))).
		Concat(linq.From(taggableTypes(azurev3.Resources, isTaggableType))).
		Concat(linq.From(taggableTypes(googlev4.Resources, isGoogleTaggableType))).
		Except(linq.From(unsupportedTerraformBlocks)).
		Distinct().
		Sort(stringLess).
		ToSlice(&TfTaggableResourceTypes)
}

func stringLess(i, j interface{}) bool {
	return strings.Compare(i.(string), j.(string)) < 0
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

func taggableTypes(resources map[string]*tfjson.Schema, taggable func(interface{}) bool) []string {
	var r []string
	for n, s := range resources {
		if taggable(s) {
			r = append(r, n)
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
