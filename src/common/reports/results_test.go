package reports

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/bridgecrewio/yor/src/common"
	"github.com/bridgecrewio/yor/src/common/structure"
	"github.com/bridgecrewio/yor/src/common/tagging/code2cloud"
	"github.com/bridgecrewio/yor/src/common/tagging/gittag"
	"github.com/bridgecrewio/yor/src/common/tagging/tags"
	tfStructure "github.com/bridgecrewio/yor/src/terraform/structure"
	"github.com/bridgecrewio/yor/tests/utils"

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

	t.Run("Test report JSON stdout", func(t *testing.T) {
		ReportServiceInst.CreateReport()
		_, _ = ReportServiceInst.report.AsJSONBytes()
		output := utils.CaptureOutput(ReportServiceInst.PrintJSONToStdout)
		lines := strings.Split(output, "\n")
		assert.NotNil(t, output)
		assert.LessOrEqual(t, 100, len(lines))
		match, _ := regexp.Match(" +\"summary\": {", []byte(lines[1]))
		assert.True(t, match)
		match, _ = regexp.Match(" +\"scanned\": \\d,", []byte(lines[2]))
		assert.True(t, match)
		match, _ = regexp.Match(" +\"newResources\": \\d,", []byte(lines[3]))
		assert.True(t, match)
		match, _ = regexp.Match(" +\"updatedResources\": \\d", []byte(lines[4]))
		assert.True(t, match)
		match, _ = regexp.Match(" +},", []byte(lines[5]))
		assert.True(t, match)
		match, _ = regexp.Match(" \"newResourceTags\": \\[", []byte(lines[6]))
		assert.True(t, match)
		match, _ = regexp.Match(" +{", []byte(lines[7]))
		assert.True(t, match)
		match, _ = regexp.Match(" +\"file\": \".*?\",$", []byte(lines[8]))
		assert.True(t, match)
		match, _ = regexp.Match(" +\"resourceId\": \".*?\",$", []byte(lines[9]))
		assert.True(t, match)
		match, _ = regexp.Match(" +\"key\": \".*?\",$", []byte(lines[10]))
		assert.True(t, match)
		match, _ = regexp.Match(" +\"oldValue\": \"\",$", []byte(lines[11]))
		assert.True(t, match)
		match, _ = regexp.Match(" +\"updatedValue\": \".*?\",$", []byte(lines[12]))
		assert.True(t, match)
		match, _ = regexp.Match(" +\"yorTraceId\": \".*?\"$", []byte(lines[13]))
		assert.True(t, match)
		match, _ = regexp.Match(" },$", []byte(lines[14]))
		assert.True(t, match)
	})

	t.Run("Test report JSON file", func(t *testing.T) {
		ReportServiceInst.CreateReport()
		reportFileName := "test.json"
		defer func() {
			err := os.Remove(reportFileName)
			if err != nil {
				assert.Fail(t, "Failed to delete the report file")
			}
		}()

		_, _ = ReportServiceInst.report.AsJSONBytes()
		ReportServiceInst.PrintJSONToFile(reportFileName)
		content, _ := ioutil.ReadFile(reportFileName)
		result := Report{}
		_ = json.Unmarshal(content, &result)

		assert.Equal(t, 5, result.Summary.Scanned)
		assert.Equal(t, 2, result.Summary.NewResources)
		assert.Equal(t, 2, result.Summary.UpdatedResources)
		assert.LessOrEqual(t, 8, len(result.NewResourceTags))
		assert.LessOrEqual(t, 4, len(result.UpdatedResourceTags))
	})

	t.Run("Test report structure", func(t *testing.T) {
		ReportServiceInst.CreateReport()
		report := ReportServiceInst.report
		assert.Equal(t, len(accumulator.GetScannedBlocks()), report.Summary.Scanned)
		assert.Equal(t, 2, report.Summary.NewResources)
		for _, tr := range report.NewResourceTags {
			assert.NotEqual(t, "", tr.YorTraceID)
			assert.NotEqual(t, "", tr.UpdatedValue)
			assert.NotEqual(t, "", tr.File)
			assert.NotEqual(t, "", tr.TagKey)
			assert.NotEqual(t, "", tr.ResourceID)
			assert.Equal(t, "", tr.OldValue)
		}

		assert.Equal(t, 2, report.Summary.UpdatedResources)
		for _, tr := range report.UpdatedResourceTags {
			assert.NotEqual(t, "", tr.YorTraceID)
			assert.NotEqual(t, "", tr.UpdatedValue)
			assert.NotEqual(t, "", tr.File)
			assert.NotEqual(t, "", tr.TagKey)
			assert.NotEqual(t, "", tr.ResourceID)
		}
	})

	t.Run("Test CLI output structure", func(t *testing.T) {
		ReportServiceInst.CreateReport()

		output := utils.CaptureOutput(ReportServiceInst.PrintToStdout(noColorBool bool))
                colors := noColorCheck(noColorBool)
		lines := strings.Split(output, "\n")
		// Verify banner
		assert.Equal(t, fmt.Sprintf("%v%vv%v", common.YorLogo, colors.Purple, common.Version), strings.Join(lines[0:6], "\n"))

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

	t.Run("Test list-tags result", func(t *testing.T) {
		grt := &gittag.GitRepoTag{}
		grt.Init()

		got := &gittag.GitOrgTag{}
		got.Init()

		ytt := &code2cloud.YorTraceTag{}
		ytt.Init()

		o := utils.CaptureOutput(func() {
			ReportServiceInst.PrintTagGroupTags(map[string][]tags.ITag{
				"git": {
					grt,
					got,
				},
				"code2cloud": {
					ytt,
				},
			})
		})

		lines := strings.Split(o, "\n")
		match, _ := regexp.Match(".*\\bGROUP\\b.*\\bTAG KEY\\b.*\\bDESCRIPTION\\b.*", []byte(lines[1]))
		assert.True(t, match)
		match, _ = regexp.Match(".*\\b(code2cloud|git)\\b.*\\b(yor_trace|git_.*?)\\b.*\\b[A-Za-z .]+\\b", []byte(lines[3]))
		assert.True(t, match)
	})
}

func setupAccumulator() *TagChangeAccumulator {
	accumulator := TagChangeAccumulatorInstance
	accumulator.AccumulateChanges(&tfStructure.TerraformBlock{
		Block: structure.Block{
			FilePath:    "/module/regional/mock.tf",
			ExitingTags: nil,
			NewTags: []tags.ITag{
				&code2cloud.YorTraceTag{
					Tag: tags.Tag{
						Key:   "yor_trace",
						Value: "mock-uuid",
					},
				},
				&gittag.GitOrgTag{
					Tag: tags.Tag{
						Key:   "git_org",
						Value: "bridgecrewio",
					},
				},
				&gittag.GitRepoTag{
					Tag: tags.Tag{
						Key:   "git_repository",
						Value: "terragoat",
					},
				},
				&gittag.GitModifiersTag{
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
				&code2cloud.YorTraceTag{
					Tag: tags.Tag{
						Key:   "yor_trace",
						Value: "mock-uuid",
					},
				},
				&gittag.GitOrgTag{
					Tag: tags.Tag{
						Key:   "git_org",
						Value: "bridgecrewio",
					},
				},
				&gittag.GitRepoTag{
					Tag: tags.Tag{
						Key:   "git_repository",
						Value: "terragoat",
					},
				},
				&gittag.GitModifiersTag{
					Tag: tags.Tag{
						Key:   "git_modifiers",
						Value: "shati",
					},
				},
				&gittag.GitLastModifiedAtTag{
					Tag: tags.Tag{
						Key:   "git_last_modified_at",
						Value: "2021-02-11T09:00:00.000Z",
					},
				},
				&gittag.GitLastModifiedByTag{
					Tag: tags.Tag{
						Key:   "git_last_modified_by",
						Value: "shati",
					},
				},
			},
			NewTags: []tags.ITag{
				&code2cloud.YorTraceTag{
					Tag: tags.Tag{
						Key:   "yor_trace",
						Value: "mock-uuid",
					},
				},
				&gittag.GitOrgTag{
					Tag: tags.Tag{
						Key:   "git_org",
						Value: "bridgecrewio",
					},
				},
				&gittag.GitRepoTag{
					Tag: tags.Tag{
						Key:   "git_repository",
						Value: "terragoat",
					},
				},
				&gittag.GitModifiersTag{
					Tag: tags.Tag{
						Key:   "git_modifiers",
						Value: "shati",
					},
				},
				&gittag.GitLastModifiedAtTag{
					Tag: tags.Tag{
						Key:   "git_last_modified_at",
						Value: "2021-02-11T10:00:00.000Z",
					},
				},
				&gittag.GitLastModifiedByTag{
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
				&code2cloud.YorTraceTag{
					Tag: tags.Tag{
						Key:   "yor_trace",
						Value: "another-uuid",
					},
				},
				&gittag.GitOrgTag{
					Tag: tags.Tag{
						Key:   "git_org",
						Value: "bridgecrewio",
					},
				},
				&gittag.GitRepoTag{
					Tag: tags.Tag{
						Key:   "git_repository",
						Value: "terragoat",
					},
				},
				&gittag.GitModifiersTag{
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
				&code2cloud.YorTraceTag{
					Tag: tags.Tag{
						Key:   "yor_trace",
						Value: "yet-another-uuid",
					},
				},
				&gittag.GitOrgTag{
					Tag: tags.Tag{
						Key:   "git_org",
						Value: "bridgecrewio",
					},
				},
				&gittag.GitRepoTag{
					Tag: tags.Tag{
						Key:   "git_repository",
						Value: "terragoat",
					},
				},
				&gittag.GitModifiersTag{
					Tag: tags.Tag{
						Key:   "git_modifiers",
						Value: "shati",
					},
				},
				&gittag.GitLastModifiedAtTag{
					Tag: tags.Tag{
						Key:   "git_last_modified_at",
						Value: "2021-02-11T09:00:00.000Z",
					},
				},
				&gittag.GitLastModifiedByTag{
					Tag: tags.Tag{
						Key:   "git_last_modified_by",
						Value: "shati",
					},
				},
			},
			NewTags: []tags.ITag{
				&gittag.GitOrgTag{
					Tag: tags.Tag{
						Key:   "git_org",
						Value: "bridgecrewio",
					},
				},
				&gittag.GitRepoTag{
					Tag: tags.Tag{
						Key:   "git_repository",
						Value: "terragoat",
					},
				},
				&gittag.GitModifiersTag{
					Tag: tags.Tag{
						Key:   "git_modifiers",
						Value: "gandalf/shati",
					},
				},
				&gittag.GitLastModifiedAtTag{
					Tag: tags.Tag{
						Key:   "git_last_modified_at",
						Value: "2021-02-11T09:15:00.000Z",
					},
				},
				&gittag.GitLastModifiedByTag{
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
				&code2cloud.YorTraceTag{
					Tag: tags.Tag{
						Key:   "yor_trace",
						Value: "yet-another-uuid-2",
					},
				},
				&gittag.GitOrgTag{
					Tag: tags.Tag{
						Key:   "git_org",
						Value: "bridgecrewio",
					},
				},
				&gittag.GitRepoTag{
					Tag: tags.Tag{
						Key:   "git_repository",
						Value: "terragoat",
					},
				},
				&gittag.GitModifiersTag{
					Tag: tags.Tag{
						Key:   "git_modifiers",
						Value: "shati",
					},
				},
				&gittag.GitLastModifiedAtTag{
					Tag: tags.Tag{
						Key:   "git_last_modified_at",
						Value: "2021-02-11T09:00:00.000Z",
					},
				},
				&gittag.GitLastModifiedByTag{
					Tag: tags.Tag{
						Key:   "git_last_modified_by",
						Value: "shati",
					},
				},
			},
			NewTags: []tags.ITag{
				&code2cloud.YorTraceTag{
					Tag: tags.Tag{
						Key:   "yor_trace",
						Value: "yet-another-uuid-2",
					},
				},
				&gittag.GitOrgTag{
					Tag: tags.Tag{
						Key:   "git_org",
						Value: "bridgecrewio",
					},
				},
				&gittag.GitRepoTag{
					Tag: tags.Tag{
						Key:   "git_repository",
						Value: "terragoat",
					},
				},
				&gittag.GitModifiersTag{
					Tag: tags.Tag{
						Key:   "git_modifiers",
						Value: "shati",
					},
				},
				&gittag.GitLastModifiedAtTag{
					Tag: tags.Tag{
						Key:   "git_last_modified_at",
						Value: "2021-02-11T09:00:00.000Z",
					},
				},
				&gittag.GitLastModifiedByTag{
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
