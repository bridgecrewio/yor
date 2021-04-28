package structure

import (
	"bridgecrewio/yor/src/common/gitservice"
	"bridgecrewio/yor/src/common/tagging/code2cloud"
	"bridgecrewio/yor/src/common/tagging/gittag"
	"bridgecrewio/yor/src/common/tagging/tags"
	"bridgecrewio/yor/src/common/utils"
	"fmt"
	"io/ioutil"
	"strings"
	"testing"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"

	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/stretchr/testify/assert"
)

func TestTerrraformParser_ParseFile(t *testing.T) {
	t.Run("parse aws eks file", func(t *testing.T) {
		p := &TerrraformParser{}
		p.Init("../../../tests/terraform/resources/", nil)
		filePath := "../../../tests/terraform/resources/eks.tf"
		taggableResources := [][]string{{"aws_vpc", "eks_vpc"}, {"aws_subnet", "eks_subnet1"}, {"aws_subnet", "eks_subnet2"}, {"aws_iam_role", "iam_for_eks"}, {"aws_eks_cluster", "eks_cluster"}}
		expectedTags := map[string]map[string]string{
			"eks_vpc":     {"Name": "${local.resource_prefix.value}-eks-vpc"},
			"eks_subnet1": {"Name": "${local.resource_prefix.value}-eks-subnet", "kubernetes.io/cluster/${local.eks_name.value}": "shared"},
			"eks_subnet2": {"Name": "${local.resource_prefix.value}-eks-subnet2", "kubernetes.io/cluster/${local.eks_name.value}": "shared"},
		}

		expectedLines := map[string]utils.Lines{
			"iam_policy_eks": {Start: 10, End: 19},
			"iam_for_eks":    {Start: 21, End: 24},
			"policy_attachment-AmazonEKSClusterPolicy": {Start: 26, End: 29},
			"policy_attachment-AmazonEKSServicePolicy": {Start: 31, End: 34},
			"eks_vpc":     {Start: 36, End: 43},
			"eks_subnet1": {Start: 45, End: 53},
			"eks_subnet2": {Start: 55, End: 63},
			"eks_cluster": {Start: 65, End: 78},
		}
		parsedBlocks, err := p.ParseFile(filePath)
		if err != nil {
			t.Errorf("failed to read hcl file because %s", err)
		}
		for _, block := range parsedBlocks {
			hclBlock := block.GetRawBlock().(*hclwrite.Block)
			if hclBlock.Type() == "resource" {
				if utils.InSlice(taggableResources, hclBlock.Labels()) {
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

		assert.Equal(t, 7, len(parsedBlocks))
	})

	t.Run("parse complex tags", func(t *testing.T) {
		p := &TerrraformParser{}
		p.Init("../../../tests/terraform/resources", nil)
		filePath := "../../../tests/terraform/resources/complex_tags.tf"
		expectedTags := map[string]map[string]string{
			"vpc_tags_one_line":                         {"Name": "tag-for-s3", "Environment": "prod"},
			"bucket_var_tags":                           {},
			"alb_with_merged_tags":                      {"Name": "tag-for-alb", "Environment": "prod", "yor_trace": "4329587194", "git_org": "bana"},
			"many_instance_tags":                        {"Name": "tag-for-instance", "Environment": "prod", "Owner": "bridgecrew", "yor_trace": "4329587194", "git_org": "bana"},
			"instance_merged_var":                       {"yor_trace": "4329587194", "git_org": "bana"},
			"instance_merged_override":                  {"Environment": "new_env"},
			"aurora_cluster_bastion_auto_scaling_group": {"git_org": "bridgecrewio", "git_repo": "platform", "yor_trace": "48564943-4cfc-403c-88cd-cbb207e0d33e", "Name": "bc-aurora-bastion"},
		}

		parsedBlocks, err := p.ParseFile(filePath)
		if err != nil {
			t.Errorf("failed to read hcl file because %s", err)
		}
		for _, block := range parsedBlocks {
			hclBlock := block.GetRawBlock().(*hclwrite.Block)
			if hclBlock.Type() == "resource" {
				resourceName := hclBlock.Labels()[1]
				expectedTagsForResource := expectedTags[resourceName]
				actualTags := block.GetExistingTags()
				assert.Equal(t, len(expectedTagsForResource), len(actualTags), fmt.Sprintf("failed to extract tags for resource %s\n", hclBlock.Labels()))
				for _, iTag := range actualTags {
					key := iTag.GetKey()
					assert.Equal(t, expectedTagsForResource[key], iTag.GetValue(), fmt.Sprintf("failed to extract tag value for resource %s\n", hclBlock.Labels()))
				}
			}
		}
	})

	t.Run("Skip collision tags block", func(t *testing.T) {
		p := &TerrraformParser{}
		p.Init("../../../tests/terraform/resources", nil)
		filePath := "../../../tests/terraform/resources/collision/main.tf"
		parsedBlocks, err := p.ParseFile(filePath)
		assert.Nil(t, parsedBlocks)
		assert.NotNil(t, err)
	})
}

func TestTerrraformParser_GetSourceFiles(t *testing.T) {
	t.Run("Get all terraform files when having module reference", func(t *testing.T) {
		directory := "../../../tests/terraform/resources/module1"
		terraformParser := TerrraformParser{}
		terraformParser.Init(directory, nil)
		expectedFiles := []string{"module1/main.tf", "module2/main.tf", "module2/outputs.tf", "module3/main.tf", "module3/outputs.tf"}
		actualFiles, err := terraformParser.GetSourceFiles(directory)
		assert.Equal(t, len(expectedFiles), len(actualFiles))
		for _, file := range actualFiles {
			splitFile := strings.Split(file, "/")
			lastTwoParts := splitFile[len(splitFile)-2:]
			assert.True(t, utils.InSlice(expectedFiles, strings.Join(lastTwoParts, "/")), fmt.Sprintf("expected file %s to be in directory\n", file))
		}
		if err != nil {
			t.Error(err)
		}
	})
}

func TestTerrraformParser_WriteFile(t *testing.T) {
	t.Run("Parse a file, tag its blocks, and write them to the file", func(t *testing.T) {
		rootDir := "../../../tests/terraform/resources"
		filePath := "../../../tests/terraform/resources/complex_tags.tf"
		originFileBytes, _ := ioutil.ReadFile(filePath)
		defer func() {
			_ = ioutil.WriteFile(filePath, originFileBytes, 0644)
		}()
		p := &TerrraformParser{}
		blameLines := CreateComplexTagsLines()
		gitService := &gitservice.GitService{
			BlameByFile: map[string]*git.BlameResult{filePath: {
				Lines: blameLines,
			}},
		}
		tagGroup := &gittag.TagGroup{GitService: gitService}
		c2cTagGroup := &code2cloud.TagGroup{}
		tagGroup.InitTagGroup(rootDir, nil)
		c2cTagGroup.InitTagGroup("", nil)
		p.Init(rootDir, nil)
		writeFilePath := "../../../tests/terraform/resources/tagged/complex_tags_tagged.tf"
		writeFileBytes, _ := ioutil.ReadFile(writeFilePath)
		defer func() {
			_ = ioutil.WriteFile(writeFilePath, writeFileBytes, 0644)
		}()
		parsedBlocks, err := p.ParseFile(filePath)
		if err != nil {
			t.Errorf("failed to read hcl file because %s", err)
		}

		for _, block := range parsedBlocks {
			if block.IsBlockTaggable() {
				tagGroup.CreateTagsForBlock(block)
				c2cTagGroup.CreateTagsForBlock(block)
			}
		}

		err = p.WriteFile(filePath, parsedBlocks, writeFilePath)
		if err != nil {
			t.Error(err)
		}
		parsedTaggedFileTags, err := p.ParseFile(writeFilePath)
		if err != nil {
			t.Error(err)
		}

		for _, block := range parsedTaggedFileTags {
			if block.IsBlockTaggable() {
				isYorTagExists := false
				yorTagKey := tags.YorTraceTagKey
				for _, tag := range block.GetExistingTags() {
					if tag.GetKey() == yorTagKey || strings.ReplaceAll(tag.GetKey(), `"`, "") == yorTagKey {
						isYorTagExists = true
					}
				}
				if !isYorTagExists {
					t.Error(fmt.Sprintf("tag not found on merged block %v", yorTagKey))
				}
			}
		}
	})

	t.Run("Test parsing of unsupported blocks", func(t *testing.T) {
		p := &TerrraformParser{}
		p.Init("../../../tests/terraform/mixed", nil)
		blocks, err := p.ParseFile("../../../tests/terraform/mixed/mixed.tf")
		if err != nil {
			t.Fail()
		}
		assert.Equal(t, 1, len(blocks))
		assert.Equal(t, "aws_s3_bucket.test-bucket", blocks[0].GetResourceID())
	})
}

var layout = "2006-01-02 15:04:05"

func getTime() time.Time {
	t, _ := time.Parse(layout, "2020-06-16 17:46:24")
	return t
}

func CreateComplexTagsLines() []*git.Line {
	originFileText, err := ioutil.ReadFile("../../../tests/terraform/resources/complex_tags.tf")
	if err != nil {
		panic(err)
	}
	originLines := strings.Split(string(originFileText), "\n")
	lines := make([]*git.Line, 0)

	for _, line := range originLines {
		lines = append(lines, &git.Line{
			Author: "user@gmail.com",
			Text:   line,
			Date:   getTime(),
			Hash:   plumbing.NewHash("hash"),
		})
	}

	return lines
}
