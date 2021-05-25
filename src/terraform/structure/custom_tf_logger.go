package structure

import "github.com/bridgecrewio/yor/src/common/logger"

type customTfLogger struct{}

func (c customTfLogger) Ask(_ string) (string, error) {
	return "", nil
}

func (c customTfLogger) AskSecret(_ string) (string, error) {
	return "", nil
}

func (c customTfLogger) Output(s string) {
	logger.Info(s)
}

func (c customTfLogger) Info(s string) {
	logger.Info(s)
}

func (c customTfLogger) Error(s string) {
	logger.Error(s)
}

func (c customTfLogger) Warn(s string) {
	logger.Warning(s)
}
