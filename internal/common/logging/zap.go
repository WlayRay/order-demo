package logging

import (
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func Init() {
	var config zap.Config
	//if viper.GetBool("development") {
	config = zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	//} else {
	//	config.Encoding = "json"
	//	config = zap.NewProductionConfig()
	//	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	//}

	zapLevel := viper.GetString("zap-level")
	switch zapLevel {
	case "debug":
		config.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	case "info":
		config.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	case "warn":
		config.Level = zap.NewAtomicLevelAt(zapcore.WarnLevel)
	case "error":
		config.Level = zap.NewAtomicLevelAt(zapcore.ErrorLevel)
	case "fatal":
		config.Level = zap.NewAtomicLevelAt(zapcore.FatalLevel)
	default:
		config.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	}

	//// 配置输出到文件
	//if logFile := viper.GetString("log-file"); logFile != "" {
	//	config.OutputPaths = append(config.OutputPaths, logFile)
	//	config.ErrorOutputPaths = append(config.ErrorOutputPaths, logFile)
	//}

	logger, err := config.Build(
		zap.AddCaller(),
		zap.AddCallerSkip(1),
		zap.AddStacktrace(zapcore.WarnLevel),
	)
	if err != nil {
		panic("failed to initialize logger: " + err.Error())
	}
	zap.ReplaceGlobals(logger)
}
