package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLogger(t *testing.T) {
	t.Run("Test logger setup", func(t *testing.T) {
		assert.Equal(t, WARNING, Logger.logLevel)
		Logger.Warning("Test warning")
	})

	t.Run("Test logger DEBUG logLevel setting", func(t *testing.T) {
		Logger.SetLogLevel("DEBUG")
		assert.Equal(t, DEBUG, Logger.logLevel)
		Logger.Debug("Test debug")
	})

	t.Run("Test logger INFO logLevel setting", func(t *testing.T) {
		Logger.SetLogLevel("INFO")
		assert.Equal(t, INFO, Logger.logLevel)
		Logger.Info("Test info")
	})

	t.Run("Test logger WARNING logLevel setting", func(t *testing.T) {
		Logger.SetLogLevel("WARNING")
		assert.Equal(t, WARNING, Logger.logLevel)
		Logger.Warning("Test warning")
	})

	t.Run("Test logger ERROR logLevel setting", func(t *testing.T) {
		Logger.SetLogLevel("ERROR")
		assert.Equal(t, ERROR, Logger.logLevel)
		assert.Panics(t, func() {
			Logger.Error("Test error")
		})
	})
}
