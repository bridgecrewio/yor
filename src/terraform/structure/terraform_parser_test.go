package structure

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/bridgecrewio/yor/src/common/gitservice"
	"github.com/bridgecrewio/yor/src/common/structure"
	"github.com/bridgecrewio/yor/src/common/tagging/code2cloud"
	"github.com/bridgecrewio/yor/src/common/tagging/gittag"
	"github.com/bridgecrewio/yor/src/common/tagging/tags"
	"github.com/bridgecrewio/yor/src/common/utils"
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

		expectedLines := map[string]structure.Lines{
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
			if hclBlock.Type() == ResourceBlockType {
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

			if hclBlock.Type() == ResourceBlockType || hclBlock.Type() == DataBlockType {
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
			"instance_null_tags":                        nil,
		}

		parsedBlocks, err := p.ParseFile(filePath)
		if err != nil {
			t.Errorf("failed to read hcl file because %s", err)
		}
		for _, block := range parsedBlocks {
			hclBlock := block.GetRawBlock().(*hclwrite.Block)
			if hclBlock.Type() == ResourceBlockType {
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

func TestTerrraformParser(t *testing.T) {
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

func TestTerrraformParser_Module(t *testing.T) {
	t.Run("Parse a file, tag its blocks, and write them to the file", func(t *testing.T) {
		rootDir := "../../../tests/terraform/resources"
		filePath := "../../../tests/terraform/resources/complex_tags.tf"
		originFileBytes, _ := ioutil.ReadFile(filePath)
		defer func() {
			_ = ioutil.WriteFile(filePath, originFileBytes, 0644)
		}()
		p := &TerrraformParser{}
		blameLines := CreateComplexTagsLines()
		gitService := &gitservice.GitService{}
		var blameByFile sync.Map
		blameByFile.Store(filePath, &git.BlameResult{Lines: blameLines})
		gitService.BlameByFile = &blameByFile
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
			if utils.InSlice([]string{"aws_autoscaling_group.autoscaling_group", "aws_autoscaling_group.autoscaling_group_tagged"}, block.GetResourceID()) {
				assert.False(t, block.IsBlockTaggable())
			}
			if block.IsBlockTaggable() {
				_ = tagGroup.CreateTagsForBlock(block)
				_ = c2cTagGroup.CreateTagsForBlock(block)
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

	t.Run("Parse a file with escaped tags, tag its blocks, and write them to the file", func(t *testing.T) {
		rootDir := "../../../tests/terraform/resources/k8s_tf"
		filePath := "../../../tests/terraform/resources/k8s_tf/main.tf"
		originFileBytes, _ := ioutil.ReadFile(filePath)
		defer func() {
			_ = ioutil.WriteFile(filePath, originFileBytes, 0644)
		}()
		p := &TerrraformParser{}
		c2cTagGroup := &code2cloud.TagGroup{}
		c2cTagGroup.InitTagGroup("", nil)
		p.Init(rootDir, nil)
		writeFilePath := "../../../tests/terraform/resources/k8s_tf/main.tf"
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
				_ = c2cTagGroup.CreateTagsForBlock(block)
			} else {
				assert.Fail(t, fmt.Sprintf("Block %v should be taggable!", block.GetResourceID()))
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
					assert.NotEqualf(t, "kubernetes.io/cluster/$${local.prefix}", tag.GetKey(), "Bad tag exists!")
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

	t.Run("Test reading & writing of module block", func(t *testing.T) {
		p := &TerrraformParser{}
		p.Init("../../../tests/terraform/module/module_with_tags", nil)
		sourceFilePath := "../../../tests/terraform/module/module_with_tags/main.tf"
		expectedFileName := "../../../tests/terraform/module/module_with_tags/expected.txt"
		blocks, err := p.ParseFile(sourceFilePath)
		if err != nil {
			t.Fail()
		}
		assert.Equal(t, 1, len(blocks))
		mb := blocks[0]
		assert.Equal(t, "complete_sg", mb.GetResourceID())
		assert.Equal(t, "tags", mb.(*TerraformBlock).TagsAttributeName)
		mb.AddNewTags([]tags.ITag{
			&tags.Tag{Key: tags.YorTraceTagKey, Value: "some-uuid"},
			&tags.Tag{Key: "mock_tag", Value: "mock_value"},
		})

		resultFileName := "result.txt"
		defer func() {
			_ = os.Remove(resultFileName)
		}()
		_ = p.WriteFile(sourceFilePath, blocks, resultFileName)
		resultStr, _ := ioutil.ReadFile(resultFileName)
		expectedStr, _ := ioutil.ReadFile(expectedFileName)
		assert.Equal(t, string(resultStr), string(expectedStr))
	})

	t.Run("Test taggable unaccessible module", func(t *testing.T) {
		p := &TerrraformParser{}
		p.Init("../../../tests/terraform/module/tfe_module", nil)
		sourceFilePath := "../../../tests/terraform/module/tfe_module/main.tf"
		blocks, err := p.ParseFile(sourceFilePath)
		if err != nil {
			t.Fail()
		}

		assert.Equal(t, 1, len(blocks))
		moduleBlock := blocks[0]
		assert.True(t, moduleBlock.IsBlockTaggable())
		assert.NotNil(t, 2, len(moduleBlock.GetExistingTags()))
	})

	t.Run("Test reading & writing of module block without tags", func(t *testing.T) {
		p := &TerrraformParser{}
		p.Init("../../../tests/terraform/module/module", nil)
		sourceFilePath := "../../../tests/terraform/module/module/main.tf"
		expectedFileName := "../../../tests/terraform/module/module/expected.txt"
		blocks, err := p.ParseFile(sourceFilePath)
		if err != nil {
			t.Fail()
		}
		assert.Equal(t, 1, len(blocks))
		mb := blocks[0]
		assert.Equal(t, "complete_sg", mb.GetResourceID())
		assert.True(t, mb.IsBlockTaggable())
		assert.Equal(t, "tags", mb.(*TerraformBlock).TagsAttributeName)
		mb.AddNewTags([]tags.ITag{
			&tags.Tag{Key: tags.YorTraceTagKey, Value: "some-uuid"},
			&tags.Tag{Key: "mock_tag", Value: "mock_value"},
		})

		resultFileName := "result.txt"
		defer func() {
			_ = os.Remove(resultFileName)
		}()
		_ = p.WriteFile(sourceFilePath, blocks, resultFileName)
		resultStr, _ := ioutil.ReadFile(resultFileName)
		expectedStr, _ := ioutil.ReadFile(expectedFileName)
		assert.Equal(t, string(expectedStr), string(resultStr))
	})

	t.Run("TestTagsAttributeScenarios", func(t *testing.T) {
		p := &TerrraformParser{}
		p.Init("../../../tests/terraform/resources/attributescenarios", nil)
		filePath := "../../../tests/terraform/resources/attributescenarios/main.tf"
		resultFilePath := "../../../tests/terraform/resources/attributescenarios/main_result.tf"
		expectedFilePath := "../../../tests/terraform/resources/attributescenarios/expected.txt"
		blocks, _ := p.ParseFile(filePath)
		assert.Equal(t, 8, len(blocks))
		for _, block := range blocks {
			if block.IsBlockTaggable() {
				block.AddNewTags([]tags.ITag{
					&tags.Tag{Key: "git_repo", Value: "yor"},
					&tags.Tag{Key: "git_org", Value: "bridgecrewio"},
				})
			}
		}

		_ = p.WriteFile(filePath, blocks, resultFilePath)
		defer func() {
			_ = os.Remove(resultFilePath)
		}()

		result, _ := ioutil.ReadFile(resultFilePath)
		expected, _ := ioutil.ReadFile(expectedFilePath)
		assert.Equal(t, string(expected), string(result))
	})

	t.Run("Module isTaggable local/remote", func(t *testing.T) {
		directory := "../../../tests/terraform/resources/local_module"
		terraformParser := TerrraformParser{}
		terraformParser.Init(directory, nil)
		expectedFiles := []string{"main.tf", "sub_local_module/main.tf", "sub_local_module/variables.tf"}
		for _, file := range expectedFiles {
			filePath := filepath.Join(directory, file)
			fileBlocks, err := terraformParser.ParseFile(filePath)
			if err != nil {
				assert.Fail(t, fmt.Sprintf("Failed to parse file %v", filePath))
			}
			for _, b := range fileBlocks {
				switch b.GetResourceID() {
				case "sub_module":
					assert.False(t, b.IsBlockTaggable())
				case "sg":
					assert.True(t, b.IsBlockTaggable())
				}
			}
		}
	})

	t.Run("Test isModuleTaggable on remote modules", func(t *testing.T) {
		directory := "../../../tests/terraform/module/provider_modules"
		terraformParser := TerrraformParser{}
		terraformParser.Init(directory, nil)
		blocks, _ := terraformParser.ParseFile(directory + "/main.tf")
		assert.Equal(t, 8, len(blocks))
		for _, block := range blocks {
			assert.True(t, block.IsBlockTaggable(), fmt.Sprintf("Block %v should be taggable", block.GetResourceID()))
		}
	})
}

func TestExtractProviderFromModuleSrc(t *testing.T) {
	tests := []struct {
		name   string
		source string
		want   string
	}{
		{name: "registry_aws_module", source: "terraform-aws-modules/security-group/aws", want: "aws"},
		{name: "git_aws_module", source: "git@github.com:terraform-aws-modules/terraform-aws-vpc.git", want: "aws"},
		{name: "github_aws_module", source: "github.com/terraform-aws-modules/terraform-aws-vpc", want: "aws"},
		{name: "private_registry_aws_module", source: "app.terraform.io/acme/rds/aws", want: "aws"},
		{name: "registry_google_module", source: "terraform-google-modules/network/google", want: "google"},
		{name: "git_google_module", source: "git@github.com:terraform-google-modules/terraform-google-network.git", want: "google"},
		{name: "github_google_module", source: "github.com/terraform-google-modules/terraform-google-network.git", want: "google"},
		{name: "private_registry_google_module", source: "app.terraform.io/acme/network/google", want: "google"},
		{name: "azure_github", source: "git@github.com:aztfmod/terraform-azurerm-caf.git", want: "azurerm"},
		{name: "repo_with_ref", source: "claranet/run-common/azurerm//modules/logs", want: "azurerm"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ExtractProviderFromModuleSrc(tt.source); got != tt.want {
				t.Errorf("ExtractProviderFromModuleSrc() = %v, want %v", got, tt.want)
			}
		})
	}
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
