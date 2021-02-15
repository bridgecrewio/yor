package logger

import (
	"bridgecrewio/yor/tests/utils"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLogger(t *testing.T) {
	t.Run("Test logger setup", func(t *testing.T) {
		assert.Equal(t, WARNING, Logger.logLevel)
		logs := utils.CaptureOutput(func() { Warning("Test warning") })
		match, _ := regexp.Match("\\d{4}/\\d{2}/\\d{2} \\d{1,2}:\\d{1,2}:\\d{1,2} Warning: \n\\d{4}/\\d{2}/\\d{2} \\d{1,2}:\\d{1,2}:\\d{1,2} \\[Test warning]", []byte(logs))
		assert.True(t, match)
	})

	t.Run("Test logger DEBUG logLevel setting", func(t *testing.T) {
		Logger.SetLogLevel("DEBUG")
		assert.Equal(t, DEBUG, Logger.logLevel)
		logs := utils.CaptureOutput(func() { Debug("Test debug") })
		match, _ := regexp.Match("\\d{4}/\\d{2}/\\d{2} \\d{1,2}:\\d{1,2}:\\d{1,2} \\[Test debug]", []byte(logs))
		assert.True(t, match)
	})

	t.Run("Test logger INFO logLevel setting", func(t *testing.T) {
		Logger.SetLogLevel("INFO")
		assert.Equal(t, INFO, Logger.logLevel)
		logs := utils.CaptureOutput(func() { Info("Test info") })
		match, _ := regexp.Match("\\d{4}/\\d{2}/\\d{2} \\d{1,2}:\\d{1,2}:\\d{1,2} \\[Test info]", []byte(logs))
		assert.True(t, match)
	})

	t.Run("Test logger WARNING logLevel setting", func(t *testing.T) {
		Logger.SetLogLevel("WARNING")
		assert.Equal(t, WARNING, Logger.logLevel)
		logs := utils.CaptureOutput(func() { Warning("Test warning") })
		match, _ := regexp.Match("\\d{4}/\\d{2}/\\d{2} \\d{1,2}:\\d{1,2}:\\d{1,2} Warning: \n\\d{4}/\\d{2}/\\d{2} \\d{1,2}:\\d{1,2}:\\d{1,2} \\[Test warning]", []byte(logs))
		assert.True(t, match)
	})

	t.Run("Test logger ERROR logLevel setting", func(t *testing.T) {
		Logger.SetLogLevel("ERROR")
		assert.Equal(t, ERROR, Logger.logLevel)
		assert.Panics(t, func() {
			Error("Test error")
		})
	})

	t.Run("Test logger not logging due to logLevel - DEBUG", func(t *testing.T) {
		logs := utils.CaptureOutput(func() { Debug("Test debug 2") })
		assert.Equal(t, "", logs)
	})

	t.Run("Test logger not logging due to logLevel - INFO", func(t *testing.T) {
		logs := utils.CaptureOutput(func() { Info("Test info 2") })
		assert.Equal(t, "", logs)
	})
}
