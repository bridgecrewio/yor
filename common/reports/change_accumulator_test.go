package reports

import (
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
		accumulator := TagChangeAccumulatorInstance
		accumulator.AccumulateChanges(&tfStructure.TerraformBlock{
			Block: structure.Block{
				FilePath:    "/module/regional/mock.tf",
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
				FilePath: "/module/regional/mock.tf",
				ExitingTags: []tags.ITag{
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
					&tags.GitLastModifiedAtTag{
						Tag: tags.Tag{
							Key:   "git_last_modified_at",
							Value: "2021-02-11T09:00:00.000Z",
						},
					},
					&tags.GitLastModifiedByTag{
						Tag: tags.Tag{
							Key:   "git_last_modified_by",
							Value: "shati",
						},
					},
				},
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
					&tags.GitLastModifiedAtTag{
						Tag: tags.Tag{
							Key:   "git_last_modified_at",
							Value: "2021-02-11T10:00:00.000Z",
						},
					},
					&tags.GitLastModifiedByTag{
						Tag: tags.Tag{
							Key:   "git_last_modified_by",
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
				Labels:          []string{"aws_s3_bucket", "data_bucket"},
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
					&tags.GitLastModifiedAtTag{
						Tag: tags.Tag{
							Key:   "git_last_modified_at",
							Value: "2021-02-11T09:00:00.000Z",
						},
					},
					&tags.GitLastModifiedByTag{
						Tag: tags.Tag{
							Key:   "git_last_modified_by",
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
					&tags.GitLastModifiedAtTag{
						Tag: tags.Tag{
							Key:   "git_last_modified_at",
							Value: "2021-02-11T09:15:00.000Z",
						},
					},
					&tags.GitLastModifiedByTag{
						Tag: tags.Tag{
							Key:   "git_last_modified_by",
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
				Labels:          []string{"aws_iam_role", "eks_node_role"},
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
							Value: "yet-another-uuid-2",
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
					&tags.GitLastModifiedAtTag{
						Tag: tags.Tag{
							Key:   "git_last_modified_at",
							Value: "2021-02-11T09:00:00.000Z",
						},
					},
					&tags.GitLastModifiedByTag{
						Tag: tags.Tag{
							Key:   "git_last_modified_by",
							Value: "shati",
						},
					},
				},
				NewTags: []tags.ITag{
					&tags.YorTraceTag{
						Tag: tags.Tag{
							Key:   "yor_trace",
							Value: "yet-another-uuid-2",
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
					&tags.GitLastModifiedAtTag{
						Tag: tags.Tag{
							Key:   "git_last_modified_at",
							Value: "2021-02-11T09:00:00.000Z",
						},
					},
					&tags.GitLastModifiedByTag{
						Tag: tags.Tag{
							Key:   "git_last_modified_by",
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
				Labels:          []string{"aws_iam_role", "eks_master_role"},
				Body:            nil,
				TypeRange:       hcl.Range{},
				LabelRanges:     nil,
				OpenBraceRange:  hcl.Range{},
				CloseBraceRange: hcl.Range{},
			},
		})

		assert.Equal(t, 5, len(accumulator.GetScannedBlocks()))

		ReportServiceInst.CreateReport()
		ReportServiceInst.PrintToStdout()
		newBlocks, updatedBlocks := accumulator.GetBlockChanges()
		assert.Equal(t, 2, len(newBlocks))
		assert.Equal(t, 2, len(updatedBlocks))
	})
}
