package structure

import (
	"bridgecrewio/yor/common/structure"
	"bridgecrewio/yor/common/tagging/tags"
	"fmt"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"io/ioutil"
	"strings"
)

var prefixToTagAttribute = map[string]string{"aws": "tags", "azure": "tags", "gcp": "labels"}

type TerrraformParser struct {
}

func (p *TerrraformParser) ParseFile(filePath string) ([]structure.IBlock, error) {
	parsedBlocks := make([]structure.IBlock, 0)
	src, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s because %s", filePath, err)
	}

	hclFile, diagnostics := hclwrite.ParseConfig(src, filePath, hcl.InitialPos)
	if diagnostics != nil && diagnostics.HasErrors() {
		hclErrors := diagnostics.Errs()
		return nil, fmt.Errorf("failed to parse hcl file %s because of errors %s", filePath, hclErrors)
	}
	if hclFile == nil {
		return nil, fmt.Errorf("failed to parse hcl file %s", filePath)
	}

	rawBlocks := hclFile.Body().Blocks()
	for _, block := range rawBlocks {
		terraformBlock, err := p.parseBlock(block)
		if err != nil {
			return nil, fmt.Errorf("failed to pasrse terraform block because %s", err)
		}
		err = terraformBlock.Init(filePath, block)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize terraform block because %s", err)
		}
		parsedBlocks = append(parsedBlocks, terraformBlock)
	}

	return parsedBlocks, nil
}

func (p *TerrraformParser) WriteFile(filePath string, blocks []*structure.IBlock) error {
	// TODO
	return nil
}

func (p *TerrraformParser) parseBlock(hclBlock *hclwrite.Block) (structure.IBlock, error) {
	var existingTags []tags.ITag
	var tagsAttributeName string
	var err error
	isTaggable := false

	if hclBlock.Type() == "resource" {
		tagsAttributeName, err = p.getTagsAttributeName(hclBlock)
		if err != nil {
			return nil, err
		}
		existingTags, isTaggable, err = p.getExistingTags(hclBlock, tagsAttributeName)
		if err != nil {
			return nil, err
		}

		if !isTaggable {
			isTaggable, err = p.isBlockTaggable(hclBlock)
			if err != nil {
				return nil, err
			}

		}
	}

	terraformBlock := TerraformBlock{
		Block: structure.Block{
			ExitingTags:       existingTags,
			IsTaggable:        isTaggable,
			TagsAttributeName: tagsAttributeName,
		},
	}

	return &terraformBlock, nil
}

func (p *TerrraformParser) getTagsAttributeName(hclBlock *hclwrite.Block) (string, error) {
	resourceType := hclBlock.Labels()[0]
	tagsAttributeName, err := getTagAttributeByResourceType(resourceType)
	if err != nil {
		return "", err
	}

	return tagsAttributeName, nil
}

func getTagAttributeByResourceType(resourceType string) (string, error) {
	splitResourceType := strings.Split(resourceType, "_")
	if len(splitResourceType) == 0 {
		return "", fmt.Errorf("failed to split resource type %s", resourceType)
	}
	prefix := prefixToTagAttribute[splitResourceType[0]]
	if prefix == "" {
		return "", fmt.Errorf("failed to find tags attribute name for resource type %s", resourceType)
	}

	return prefix, nil
}

func (p *TerrraformParser) getExistingTags(hclBlock *hclwrite.Block, tagsAttributeName string) ([]tags.ITag, bool, error) {
	isTaggable := false
	existingTags := make([]tags.ITag, 0)

	tagsAttribute := hclBlock.Body().GetAttribute(tagsAttributeName)
	if tagsAttribute != nil {
		isTaggable = true
		tagsTokens := tagsAttribute.Expr().BuildTokens(hclwrite.Tokens{})
		parsedTags := p.parseTagLines(tagsTokens)
		for key := range parsedTags {
			iTag := tags.Init(key, parsedTags[key])
			existingTags = append(existingTags, iTag)
		}
	}

	return existingTags, isTaggable, nil
}

func (p *TerrraformParser) isBlockTaggable(hclBlock *hclwrite.Block) (bool, error) {
	return false, nil
}

func (p *TerrraformParser) parseTagLines(tokens hclwrite.Tokens) map[string]string {
	parsedTags := make(map[string]string)
	entries := make([]hclwrite.Tokens, 0)
	startIndex := 0
	hasEq := false
	for i, token := range tokens {
		if token.Type == hclsyntax.TokenNewline {
			if hasEq {
				entries = append(entries, tokens[startIndex:i])
			}
			startIndex = i + 1
			hasEq = false
		}
		if token.Type == hclsyntax.TokenEqual {
			hasEq = true
		}
	}
	if hasEq {
		entries = append(entries, tokens[startIndex:])
	}

	for _, entry := range entries {
		eqIndex := -1
		var key string
		for j, token := range entry {
			if token.Type == hclsyntax.TokenEqual {
				eqIndex = j + 1
				key = strings.TrimSpace(string(entry[:j].Bytes()))
			}
		}
		value := string(entry[eqIndex:].Bytes())
		value = strings.TrimPrefix(strings.TrimSuffix(value, " "), " ")
		parsedTags[key] = value
	}

	return parsedTags
}
