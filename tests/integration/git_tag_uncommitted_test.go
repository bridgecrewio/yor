package integration

//
//import (
//	"bridgecrewio/yor/src/common"
//	"bridgecrewio/yor/src/common/gitservice"
//	"bridgecrewio/yor/src/common/runner"
//	"bridgecrewio/yor/src/common/tagging/gittag"
//	terraformStructure "bridgecrewio/yor/src/terraform/structure"
//	"bridgecrewio/yor/tests/utils"
//	"fmt"
//	"os"
//	"path"
//	"strings"
//	"testing"
//	"time"
//
//	"github.com/hashicorp/hcl/v2/hclsyntax"
//	"github.com/stretchr/testify/assert"
//)
//
//func failIfErr(t *testing.T, err error) {
//	if err != nil {
//		t.Error(err)
//	}
//}
//
//func tagDirectory(t *testing.T, path string) {
//	yorRunner := runner.Runner{}
//	err := yorRunner.Init(&common.Options{
//		Directory: path,
//		ExtraTags: "{}",
//	})
//	failIfErr(t, err)
//	_, err = yorRunner.TagDirectory()
//	failIfErr(t, err)
//}
//
//func TestTagUncommittedResults(t *testing.T) {
//	t.Run("Test terragoat tagging", func(t *testing.T) {
//		terragoatPath := utils.CloneRepo(utils.TerragoatURL)
//		outputPath := "./result_uncommitted.json"
//		defer func() {
//			os.RemoveAll(terragoatPath)
//			os.RemoveAll(outputPath)
//		}()
//
//		terragoatAWSDirectory := path.Join(terragoatPath, "terraform/aws")
//
//		// tag aws directory
//		tagDirectory(t, terragoatAWSDirectory)
//		// tag again, this time the files have uncommitted changes
//		tagDirectory(t, terragoatAWSDirectory)
//
//		terrraformParser := terraformStructure.TerrraformParser{}
//		terrraformParser.Init(terragoatAWSDirectory, nil)
//
//		dbAppFile := path.Join(terragoatAWSDirectory, "db-app.tf")
//		blocks, err := terrraformParser.ParseFile(dbAppFile)
//		failIfErr(t, err)
//		defaultInstanceBlock := blocks[0].(*terraformStructure.TerraformBlock)
//		if defaultInstanceBlock.GetResourceID() != "aws_db_instance.default" {
//			t.Errorf("invalid file structure, the resource id is %s", defaultInstanceBlock.GetResourceID())
//		}
//
//		rawTags := defaultInstanceBlock.HclSyntaxBlock.Body.Attributes["tags"]
//		rawTagsExpr := rawTags.Expr.(*hclsyntax.FunctionCallExpr)
//		assert.Equal(t, "merge", rawTagsExpr.Name)
//		mergeArgs := rawTagsExpr.Args
//		assert.Equal(t, 2, len(mergeArgs))
//		assert.Equal(t, 2, len(mergeArgs[0].(*hclsyntax.ObjectConsExpr).Items))
//		assert.Equal(t, 8, len(mergeArgs[1].(*hclsyntax.ObjectConsExpr).Items))
//
//		currentTags := defaultInstanceBlock.ExitingTags
//
//		expectedTagsValues := map[string]string{
//			"Name":                 "${local.resource_prefix.value}-rds",
//			"Environment":          "local.resource_prefix.value",
//			"git_last_modified_by": gitservice.GetGitUserEmail(),
//			"git_commit":           gittag.CommitUnavailable,
//			"git_file":             strings.TrimPrefix(dbAppFile, terragoatPath+"/"),
//		}
//
//		for _, tag := range currentTags {
//			if expectedVal, ok := expectedTagsValues[tag.GetKey()]; ok {
//				assert.Equal(t, expectedVal, tag.GetValue(), fmt.Sprintf("Missmach in tag %s, expected %s, got %s", tag.GetKey(), expectedVal, tag.GetValue()))
//			}
//			if tag.GetKey() == "git_last_modified_at" {
//				timeTagValue, err := time.Parse("2006-01-02 15:04:05", tag.GetValue())
//				failIfErr(t, err)
//				diff := time.Now().UTC().Sub(timeTagValue)
//				assert.True(t, diff < 2*time.Minute)
//			}
//		}
//
//	})
//}
