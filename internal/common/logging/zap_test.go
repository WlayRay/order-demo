package logging

import (
	"testing"

	"go.uber.org/zap"
)

func TestZap(t *testing.T) {
	Init()

	message := "This is a test message"
	zap.L().Debug(message, zap.String("key", "value"), zap.Int("number", 42))
	zap.L().Info(message, zap.String("key", "value"), zap.Int("number", 42))
	zap.L().Warn(message, zap.String("key", "value"), zap.Int("number", 42))
	zap.L().Error(message, zap.String("key", "value"), zap.Int("number", 42))

	// zap.L().Fatal(message, zap.String("key", "value"), zap.Int("number", 42))

	zap.L().DPanic(message, zap.String("key", "value"), zap.Int("number", 42))
	zap.L().Panic(message, zap.String("key", "value"), zap.Int("number", 42))
	zap.L().Info(message, zap.String("key", "value"), zap.Int("number", 42))

	zap.L().Info("SADD", zap.Bool("key", true))
	zap.L().Sync()
}
