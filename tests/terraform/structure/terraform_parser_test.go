package structure

import (
	"bridgecrewio/yor/common"
	"bridgecrewio/yor/terraform/structure"
	"fmt"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestTerrraformParser_ParseFile(t *testing.T) {
	t.Run("parse aws eks file", func(t *testing.T) {
		p := &structure.TerrraformParser{}
		p.Init("../resources/", nil)
		filePath := "../resources/eks.tf"
		taggableResources := [][]string{{"aws_vpc", "eks_vpc"}, {"aws_subnet", "eks_subnet1"}, {"aws_subnet", "eks_subnet2"}, {"aws_iam_role", "iam_for_eks"}, {"aws_eks_cluster", "eks_cluster"}}
		expectedTags := map[string]map[string]string{
			"eks_vpc":     {"Name": "\"${local.resource_prefix.value}-eks-vpc\""},
			"eks_subnet1": {"Name": "\"${local.resource_prefix.value}-eks-subnet\"", "\"kubernetes.io/cluster/${local.eks_name.value}\"": "\"shared\""},
			"eks_subnet2": {"Name": "\"${local.resource_prefix.value}-eks-subnet2\"", "\"kubernetes.io/cluster/${local.eks_name.value}\"": "\"shared\""},
		}

		expectedLines := map[string][]int{
			"iam_policy_eks": {10, 19},
			"iam_for_eks":    {21, 24},
			"policy_attachment-AmazonEKSClusterPolicy": {26, 29},
			"policy_attachment-AmazonEKSServicePolicy": {31, 34},
			"eks_vpc":     {36, 43},
			"eks_subnet1": {45, 53},
			"eks_subnet2": {55, 63},
			"eks_cluster": {65, 78},
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
			} else {
				assert.False(t, block.IsBlockTaggable())
			}

			if hclBlock.Type() == "resource" || hclBlock.Type() == "data" {
				name := hclBlock.Labels()[1]
				expectedBlockLines := expectedLines[name]
				actualLines := block.GetLines()
				assert.Equal(t, expectedBlockLines, actualLines)
			}
		}

		assert.Equal(t, 11, len(parsedBlocks))
	})
}

func TestTerrraformParser_GetSourceFiles(t *testing.T) {
	t.Run("Get all terraform files when having module reference", func(t *testing.T) {
		directory := "../resources/module1"
		terraformParser := structure.TerrraformParser{}
		terraformParser.Init(directory, nil)
		expectedFiles := []string{"module1/main.tf", "module2/main.tf", "module2/outputs.tf"}
		actualFiles, err := terraformParser.GetSourceFiles(directory)
		assert.Equal(t, len(expectedFiles), len(actualFiles))
		for _, file := range actualFiles {
			splitFile := strings.Split(file, "/")
			lastTwoParts := splitFile[len(splitFile)-2:]
			assert.True(t, common.InSlice(expectedFiles, strings.Join(lastTwoParts, "/")), fmt.Sprintf("expected file %s to be in directory\n", file))
		}
		if err != nil {
			t.Error(err)
		}
	})
}
