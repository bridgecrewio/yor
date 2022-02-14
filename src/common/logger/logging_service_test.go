package logger

import (
	"regexp"
	"strings"
	"testing"

	"github.com/bridgecrewio/yor/tests/utils"

	"github.com/stretchr/testify/assert"
)

func TestLogger(t *testing.T) {
	t.Run("Test logger setup", func(t *testing.T) {
		assert.Equal(t, WARNING, Logger.logLevel)
		logs := utils.CaptureOutput(func() { Warning("Test warning") })
		match, _ := regexp.Match("\\d{4}/\\d{2}/\\d{2} \\d{1,2}:\\d{1,2}:\\d{1,2} \\[WARNING] Test warning", []byte(logs))
		assert.True(t, match)
	})

	t.Run("Test logger DEBUG logLevel setting", func(t *testing.T) {
		Logger.SetLogLevel("DEBUG")
		assert.Equal(t, DEBUG, Logger.logLevel)
		logs := utils.CaptureOutput(func() { Debug("Test debug") })
		match, _ := regexp.Match("\\d{4}/\\d{2}/\\d{2} \\d{1,2}:\\d{1,2}:\\d{1,2} \\[DEBUG] Test debug", []byte(logs))
		assert.True(t, match)
	})

	t.Run("Test logger INFO logLevel setting", func(t *testing.T) {
		Logger.SetLogLevel("INFO")
		assert.Equal(t, INFO, Logger.logLevel)
		logs := utils.CaptureOutput(func() { Info("Test info") })
		match, _ := regexp.Match("\\d{4}/\\d{2}/\\d{2} \\d{1,2}:\\d{1,2}:\\d{1,2} \\[INFO] Test info", []byte(logs))
		assert.True(t, match)
	})

	t.Run("Test logger WARNING logLevel setting", func(t *testing.T) {
		Logger.SetLogLevel("WARNING")
		assert.Equal(t, WARNING, Logger.logLevel)
		logs := utils.CaptureOutput(func() { Warning("Test warning") })
		match, _ := regexp.Match("\\d{4}/\\d{2}/\\d{2} \\d{1,2}:\\d{1,2}:\\d{1,2} \\[WARNING] Test warning", []byte(logs))
		assert.True(t, match)
	})

	t.Run("Test logger ERROR logLevel setting", func(t *testing.T) {
		Logger.SetLogLevel("ERROR")
		assert.Equal(t, ERROR, Logger.logLevel)
	})

	t.Run("Test logger not logging due to logLevel - DEBUG", func(t *testing.T) {
		logs := utils.CaptureOutput(func() { Debug("Test debug 2") })
		assert.Equal(t, "", logs)
	})

	t.Run("Test logger not logging due to logLevel - INFO", func(t *testing.T) {
		logs := utils.CaptureOutput(func() { Info("Test info 2") })
		assert.Equal(t, "", logs)
	})

	t.Run("Test mute and unmute", func(t *testing.T) {
		Logger.SetLogLevel("WARNING")

		var result string
		infoMsg := "Test muted INFO"
		warningMsg := "Test muted WARNING"
		debugMsg := "Test muted DEBUG"
		MuteOutputBlock(func() {
			result = utils.CaptureOutput(func() { Info(infoMsg) })
			assert.Equal(t, "", result)
			result = utils.CaptureOutput(func() { Warning(warningMsg) })
			assert.Equal(t, "", result)
			result = utils.CaptureOutput(func() { Debug(debugMsg) })
			assert.Equal(t, "", result)
		})
		Logger.SetLogLevel("DEBUG")
		result = utils.CaptureOutput(func() { Info(infoMsg) })
		assert.True(t, strings.Contains(result, infoMsg))
		result = utils.CaptureOutput(func() { Warning(warningMsg) })
		assert.True(t, strings.Contains(result, warningMsg))
		result = utils.CaptureOutput(func() { Debug(debugMsg) })
		assert.True(t, strings.Contains(result, debugMsg))
		Logger.SetLogLevel("WARNING")
	})
}
