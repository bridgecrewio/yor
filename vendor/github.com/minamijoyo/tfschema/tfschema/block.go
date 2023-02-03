package tfschema

import (
	"github.com/hashicorp/terraform/configs/configschema"
)

// Block is wrapper for configschema.Block.
// This ia a layer for customization not enough for Terraform's core.
// Most of the structure is the smae as the core, but some are different.
type Block struct {
	// Attributes is a map of any attributes.
	Attributes map[string]*Attribute `json:"attributes"`
	// BlockTypes is a map of any nested block types.
	BlockTypes map[string]*NestedBlock `json:"block_types"`
}

// NewBlock creates a new Block instance.
func NewBlock(b *configschema.Block) *Block {
	return &Block{
		Attributes: NewAttributes(b.Attributes),
		BlockTypes: NewBlockTypes(b.BlockTypes),
	}
}

// NewAttributes creates a new map of Attributes.
func NewAttributes(as map[string]*configschema.Attribute) map[string]*Attribute {
	m := make(map[string]*Attribute)

	for k, v := range as {
		m[k] = NewAttribute(v)
	}

	return m
}

// NewBlockTypes creates a new map of NestedBlocks.
func NewBlockTypes(bs map[string]*configschema.NestedBlock) map[string]*NestedBlock {
	m := make(map[string]*NestedBlock)

	for k, v := range bs {
		m[k] = NewNestedBlock(v)
	}

	return m
}
