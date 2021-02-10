package tagging

import (
	"bridgecrewio/yor/common/git_service"
	"bridgecrewio/yor/common/structure"
	"bridgecrewio/yor/common/tagging"
	tfStructure "bridgecrewio/yor/terraform/structure"
	"fmt"
	"reflect"
)

type TerraformTagger struct {
	tagging.Tagger
}

func (t *TerraformTagger) IsBlockTaggable(block interface{}) bool {
	tfBlock, ok := block.(*tfStructure.TerraformBlock)
	if !ok {
		return false
	}
	return tfBlock.IsBlockTaggable()
}

func (t *TerraformTagger) CreateTagsForBlock(block structure.IBlock, gitBlame *git_service.GitBlame) error {
	tfBlock, ok := block.(*tfStructure.TerraformBlock)
	if !ok {
		return fmt.Errorf("failed to convert data to *tfStructure.TerraformBlock. Type of block: %s", reflect.TypeOf(block))
	}
	for _, tag := range t.Tags {
		err := tag.CalculateValue(gitBlame)
		if err != nil {
			return fmt.Errorf("failed to calculate tag value of tag %v, err: %s", tag, err)
		}
	}
	tfBlock.AddNewTags(t.Tags)
	return nil
}
