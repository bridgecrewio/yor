package integration

import (
	"bridgecrewio/yor/src/common/reports"
	"encoding/json"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRunResults(t *testing.T) {
	t.Run("Test terragoat tagging", func(t *testing.T) {
		content, _ := ioutil.ReadFile("../../result.json")
		report := &reports.Report{}
		err := json.Unmarshal(content, &report)
		if err != nil {
			assert.Fail(t, "Failed to parse json result")
		}
		assert.LessOrEqual(t, 63, report.Summary.Scanned)
		assert.LessOrEqual(t, 63, report.Summary.NewResources)
		assert.Equal(t, 0, report.Summary.UpdatedResources)

		var taggedAWS, taggedGCP, taggedAzure int
		resourceSet := make(map[string]bool)

		for _, tr := range report.NewResourceTags {
			switch {
			case strings.HasPrefix(tr.ResourceID, "aws"):
				taggedAWS++
			case strings.HasPrefix(tr.ResourceID, "google_"):
				taggedGCP++
			case strings.HasPrefix(tr.ResourceID, "azurerm"):
				taggedAzure++
			}

			assert.NotEqual(t, "", tr.ResourceID)
			assert.NotEqual(t, "", tr.File)
			assert.NotEqual(t, "", tr.UpdatedValue)
			assert.NotEqual(t, "", tr.TagKey)
			assert.NotEqual(t, "", tr.YorTraceID)
			assert.Equal(t, "", tr.OldValue)

			resourceSet[tr.ResourceID] = true
		}

		assert.LessOrEqual(t, 312, taggedAWS)
		assert.LessOrEqual(t, 32, taggedGCP)
		assert.LessOrEqual(t, 160, taggedAzure)
		assert.Equal(t, report.Summary.Scanned, len(resourceSet))
	})
}
