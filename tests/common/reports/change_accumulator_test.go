package reports

import (
	"bridgecrewio/yor/common/reports"
	"bridgecrewio/yor/common/structure"
	"bridgecrewio/yor/common/tagging/tags"
	tfStructure "bridgecrewio/yor/terraform/structure"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTagChangeAccumulator(t *testing.T) {
	t.Run("Test changes accumulator", func(t *testing.T) {
		accumulator := reports.TagChangeAccumulatorInstance
		accumulator.AccumulateChanges(&tfStructure.TerraformBlock{
			Block: structure.Block{
				FilePath:    "/mock.tf",
				ExitingTags: nil,
				NewTags: []tags.ITag{
					&tags.YorTraceTag{
						Tag: tags.Tag{
							Key:   "yor_trace",
							Value: "mock-uuid",
						},
					},
					&tags.GitOrgTag{
						Tag: tags.Tag{
							Key:   "git_org",
							Value: "bridgecrewio",
						},
					},
					&tags.GitRepoTag{
						Tag: tags.Tag{
							Key:   "git_repository",
							Value: "terragoat",
						},
					},
					&tags.GitModifiersTag{
						Tag: tags.Tag{
							Key:   "git_modifiers",
							Value: "shati",
						},
					},
				},
				RawBlock:          nil,
				IsTaggable:        true,
				TagsAttributeName: "tag",
			},
			HclSyntaxBlock: &hclsyntax.Block{
				Type:            "",
				Labels:          []string{"aws_s3_bucket", "my_bucket"},
				Body:            nil,
				TypeRange:       hcl.Range{},
				LabelRanges:     nil,
				OpenBraceRange:  hcl.Range{},
				CloseBraceRange: hcl.Range{},
			},
		})
		accumulator.AccumulateChanges(&tfStructure.TerraformBlock{
			Block: structure.Block{
				FilePath:    "/eks.tf",
				ExitingTags: nil,
				NewTags: []tags.ITag{
					&tags.YorTraceTag{
						Tag: tags.Tag{
							Key:   "yor_trace",
							Value: "another-uuid",
						},
					},
					&tags.GitOrgTag{
						Tag: tags.Tag{
							Key:   "git_org",
							Value: "bridgecrewio",
						},
					},
					&tags.GitRepoTag{
						Tag: tags.Tag{
							Key:   "git_repository",
							Value: "terragoat",
						},
					},
					&tags.GitModifiersTag{
						Tag: tags.Tag{
							Key:   "git_modifiers",
							Value: "gandalf",
						},
					},
				},
				RawBlock:          nil,
				IsTaggable:        true,
				TagsAttributeName: "tag",
			},
			HclSyntaxBlock: &hclsyntax.Block{
				Type:            "",
				Labels:          []string{"aws_eks_cluster", "etl_jobs"},
				Body:            nil,
				TypeRange:       hcl.Range{},
				LabelRanges:     nil,
				OpenBraceRange:  hcl.Range{},
				CloseBraceRange: hcl.Range{},
			},
		})
		accumulator.AccumulateChanges(&tfStructure.TerraformBlock{
			Block: structure.Block{
				FilePath: "/iam.tf",
				ExitingTags: []tags.ITag{
					&tags.YorTraceTag{
						Tag: tags.Tag{
							Key:   "yor_trace",
							Value: "yet-another-uuid",
						},
					},
					&tags.GitOrgTag{
						Tag: tags.Tag{
							Key:   "git_org",
							Value: "bridgecrewio",
						},
					},
					&tags.GitRepoTag{
						Tag: tags.Tag{
							Key:   "git_repository",
							Value: "terragoat",
						},
					},
					&tags.GitModifiersTag{
						Tag: tags.Tag{
							Key:   "git_modifiers",
							Value: "shati",
						},
					},
				},
				NewTags: []tags.ITag{
					&tags.GitOrgTag{
						Tag: tags.Tag{
							Key:   "git_org",
							Value: "bridgecrewio",
						},
					},
					&tags.GitRepoTag{
						Tag: tags.Tag{
							Key:   "git_repository",
							Value: "terragoat",
						},
					},
					&tags.GitModifiersTag{
						Tag: tags.Tag{
							Key:   "git_modifiers",
							Value: "gandalf/shati",
						},
					},
				},
				RawBlock:          nil,
				IsTaggable:        true,
				TagsAttributeName: "tag",
			},
			HclSyntaxBlock: &hclsyntax.Block{
				Type:            "",
				Labels:          []string{"aws_iam_role", "eks_node_role"},
				Body:            nil,
				TypeRange:       hcl.Range{},
				LabelRanges:     nil,
				OpenBraceRange:  hcl.Range{},
				CloseBraceRange: hcl.Range{},
			},
		})

		assert.Equal(t, 3, len(accumulator.GetScannedBlocks()))

		reports.ReportServiceInst.CreateReport()
		//reports.ReportServiceInst.PrintToStdout()
		newBlocks, updatedBlocks := accumulator.GetBlockChanges()
		assert.Equal(t, 2, len(newBlocks))
		assert.Equal(t, 1, len(updatedBlocks))
	})
}
