package tfschema

import (
	"github.com/hashicorp/terraform/configs/configschema"
)

// Attribute is wrapper for configschema.Attribute.
type Attribute struct {
	// Type is a type of the attribute's value.
	// Note that Type is not cty.Type to customize string representation.
	Type Type `json:"type"`
	// Required is a flag whether this attribute is required.
	Required bool `json:"required"`
	// Optional is a flag whether this attribute is optional.
	// This field conflicts with Required.
	Optional bool `json:"optional"`
	// Computed is a flag whether this attribute is computed.
	// If true, the value may come from provider rather than configuration.
	// If combined with Optional, then the config may optionally provide an
	// overridden value.
	Computed bool `json:"computed"`
	// Sensitive is a flag whether this attribute may contain sensitive information.
	Sensitive bool `json:"sensitive"`
}

// NewAttribute creates a new Attribute instance.
func NewAttribute(a *configschema.Attribute) *Attribute {
	return &Attribute{
		Type:      *NewType(a.Type),
		Required:  a.Required,
		Optional:  a.Optional,
		Computed:  a.Computed,
		Sensitive: a.Sensitive,
	}
}
