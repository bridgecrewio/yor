package tfschema

import (
	"github.com/hashicorp/hcl2/ext/typeexpr"
	"github.com/zclconf/go-cty/cty"
)

// Type is a type of the attribute's value.
type Type struct {
	// We embed cty.Type to customize string representation.
	cty.Type
}

// NewType creates a new Type instance.
func NewType(t cty.Type) *Type {
	return &Type{
		Type: t,
	}
}

// MarshalJSON returns a encoded string in JSON.
func (t *Type) MarshalJSON() ([]byte, error) {
	return []byte(`"` + t.Name() + `"`), nil
}

// Name returns a name of type.
// Terraform v0.12 introduced a new `SchemaConfigModeAttr` feature.
// Most attributes have simple types, but if `SchemaConfigModeAttr` is set for
// an attribute, it is syntactically NestedBlock but semantically interpreted
// as an Attribute. In this case, Attribute has a complex data type. It is
// reasonable to use the same notation as the type annotation in HCL2 to
// represent the correct data type.
// So we use typeexpr.TypeString(cty.Type) in HCL2
//
// See also:
// - https://github.com/minamijoyo/tfschema/issues/9
// - https://github.com/terraform-providers/terraform-provider-aws/pull/8187
// - https://github.com/hashicorp/terraform/pull/20626
// - https://www.terraform.io/docs/configuration/types.html
func (t *Type) Name() string {
	return typeexpr.TypeString(t.Type)
}
