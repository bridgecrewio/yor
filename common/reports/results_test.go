package reports

import (
	"bridgecrewio/yor/common"
	"bridgecrewio/yor/common/structure"
	"bridgecrewio/yor/common/tagging/tags"
	tfStructure "bridgecrewio/yor/terraform/structure"
	"bridgecrewio/yor/tests/utils"
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/stretchr/testify/assert"
)

func TestResultsGeneration(t *testing.T) {
	accumulator := setupAccumulator()
	t.Run("Test change accumulator", func(t *testing.T) {
		assert.Equal(t, 5, len(accumulator.GetScannedBlocks()))
		newBlocks, updatedBlocks := accumulator.GetBlockChanges()
		assert.Equal(t, 2, len(newBlocks))
		assert.Equal(t, 2, len(updatedBlocks))
	})

	t.Run("Test report structure", func(t *testing.T) {
		ReportServiceInst.CreateReport()
		report := ReportServiceInst.report
		assert.Equal(t, len(accumulator.GetScannedBlocks()), report.ScannedResources)
		assert.Equal(t, 2, len(report.NewResources))
		for _, newRes := range report.NewResources {
			assert.NotNil(t, newRes.GetTraceID())
			assert.NotNil(t, newRes.MergeTags())
		}

		assert.Equal(t, 2, len(report.UpdatedResources))
		for _, updatedRes := range report.UpdatedResources {
			tagDiff := updatedRes.CalculateTagsDiff()
			for _, diff := range tagDiff.Updated {
				assert.NotNil(t, diff.Key)
				assert.NotNil(t, diff.PrevValue)
				assert.NotNil(t, diff.NewValue)
			}

			for _, diff := range tagDiff.Added {
				assert.NotNil(t, diff.GetKey())
				assert.NotNil(t, diff.GetValue())
			}
		}
	})

	t.Run("Test CLI output structure", func(t *testing.T) {
		ReportServiceInst.CreateReport()

		output := utils.CaptureOutput(ReportServiceInst.PrintToStdout)
		lines := strings.Split(output, "\n")
		// Verify banner
		assert.Equal(t, fmt.Sprintf("%v%vv%v", common.YorLogo, colorPurple, common.Version), strings.Join(lines[0:6], "\n"))

		// Verify counts
		lines = lines[7:]
		matched, _ := regexp.Match(".*?Scanned Resources:.*?\\b\\d\\b", []byte(lines[0]))
		assert.True(t, matched)
		matched, _ = regexp.Match(".*?New Resources Traced:.*?\\b\\d\\b", []byte(lines[1]))
		assert.True(t, matched)
		matched, _ = regexp.Match(".*?Updated Resources:.*?\\b\\d\\b", []byte(lines[2]))
		assert.True(t, matched)
		assert.Equal(t, "", lines[3])

		// Verify New Resources Table
		lines = lines[4:]
		matched, _ = regexp.Match(".*?New Resources Traced \\(\\d\\):", []byte(lines[0]))
		assert.True(t, matched)
		matched, _ = regexp.Match("[|\\s]+FILE[|\\s]+RESOURCE[|\\s]+TAG KEY[|\\s]+TAG VALUE[|\\s]+YOR ID[|\\s]+", []byte(lines[2]))
		assert.True(t, matched)
		matched, _ = regexp.Match("[|\\s]+[a-z./]+[|\\s]+[a-z\\d._]+", []byte(lines[4]))
		assert.True(t, matched)
		matched, _ = regexp.Match("(|[\\s]+){3}[a-z._]+[\\s]+|[\\s]+[a-z]+", []byte(lines[6]))
		assert.True(t, matched)
		matched, _ = regexp.Match("(|[\\s]+){3}[a-z._]+[\\s]+|[\\s]+[a-z]+", []byte(lines[8]))
		assert.True(t, matched)

		// Verify Updated Resources Table
		lines = lines[21:]
		matched, _ = regexp.Match(".*?Updated Resource Traces \\(\\d\\):", []byte(lines[0]))
		assert.True(t, matched)
		matched, _ = regexp.Match("[|\\s]+FILE[|\\s]+RESOURCE[|\\s]+TAG KEY[|\\s]+OLD VALUE[|\\s]+UPDATED VALUE[|\\s]+YOR ID[|\\s]+", []byte(lines[2]))
		assert.True(t, matched)
		matched, _ = regexp.Match("[|\\s]+[a-z./]+[|\\s]+[a-z\\d._]+[|\\s]+.*?[a-z\\d._:\\-]+[|\\s]+.*?[a-z\\d._:\\-]+[|\\s]+.*?[a-z\\d._:\\-]+[|\\s]+[a-z\\d-]+[|\\s]+", []byte(lines[4]))
		assert.True(t, matched)
	})
}

func setupAccumulator() *TagChangeAccumulator {
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
	return accumulator
}
