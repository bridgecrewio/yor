package structure

import (
	"bridgecrewio/yor/common"
	"bridgecrewio/yor/terraform/structure"
	"fmt"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTerrraformParser_ParseFile(t *testing.T) {
	t.Run("parse aws file", func(t *testing.T) {
		p := &structure.TerrraformParser{}
		filePath := "../resources/eks.tf"
		taggableResources := [][]string{{"aws_vpc", "eks_vpc"}, {"aws_subnet", "eks_subnet1"}, {"aws_subnet", "eks_subnet2"}}
		expectedTags := map[string]map[string]string{
			"eks_vpc":     {"Name": "\"${local.resource_prefix.value}-eks-vpc\""},
			"eks_subnet1": {"Name": "\"${local.resource_prefix.value}-eks-subnet\"", "\"kubernetes.io/cluster/${local.eks_name.value}\"": "\"shared\""},
			"eks_subnet2": {"Name": "\"${local.resource_prefix.value}-eks-subnet2\"", "\"kubernetes.io/cluster/${local.eks_name.value}\"": "\"shared\""},
		}
		parsedBlocks, err := p.ParseFile(filePath)
		if err != nil {
			t.Errorf("failed to read hcl file because %s", err)
		}
		for _, block := range parsedBlocks {
			hclBlock := block.GetRawBlock().(*hclwrite.Block)
			if hclBlock.Type() == "resource" {
				if common.InSlice(taggableResources, hclBlock.Labels()) {
					assert.True(t, block.IsBlockTaggable(), fmt.Sprintf("expected block %s to be taggable", hclBlock.Labels()))
					resourceName := hclBlock.Labels()[1]
					expectedTagsForResource := expectedTags[resourceName]
					actualTags := block.GetExistingTags()
					assert.Equal(t, len(expectedTagsForResource), len(actualTags))
					for _, iTag := range actualTags {
						key := iTag.GetKey()
						assert.Equal(t, expectedTagsForResource[key], iTag.GetValue())
					}

				} else {
					assert.False(t, block.IsBlockTaggable(), fmt.Sprintf("expected block %s not to be taggable", hclBlock.Labels()))
				}
			}
		}

		assert.Equal(t, 11, len(parsedBlocks))
	})
}
