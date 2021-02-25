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

		var taggedAWS, taggedGCP, taggedAzure bool

		for _, tr := range report.NewResourceTags {
			switch {
			case strings.HasPrefix(tr.ResourceID, "aws"):
				taggedAWS = true
			case strings.HasPrefix(tr.ResourceID, "google_"):
				taggedGCP = true
			case strings.HasPrefix(tr.ResourceID, "azurerm"):
				taggedAzure = true
			}

			if taggedAWS && taggedGCP && taggedAzure {
				break
			}
		}

		assert.True(t, taggedAWS && taggedGCP && taggedAzure)
	})
}
