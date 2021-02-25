package integration

import (
	"bridgecrewio/yor/src/common/reports"
	"encoding/json"
	"io/ioutil"
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
		assert.LessOrEqual(t, 39, report.Summary.Scanned)
		assert.LessOrEqual(t, 39, report.Summary.NewResources)
		assert.Equal(t, 0, report.Summary.UpdatedResources)
	})
}
