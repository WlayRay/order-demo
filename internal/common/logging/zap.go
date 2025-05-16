package logging

import (
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func Init() {
	config := zap.NewDevelopmentConfig()
	config.Encoding = "json"

	config.EncoderConfig = zapcore.EncoderConfig{
		TimeKey:        "time",          // 时间字段名称
		LevelKey:       "level",         // 日志级别字段名称
		NameKey:        "logger",        // 日志器名称字段名称（可选）
		CallerKey:      "caller",        // 调用栈字段名称
		FunctionKey:    zapcore.OmitKey, // 省略函数名字段
		MessageKey:     "message",       // 日志消息字段名称
		StacktraceKey:  "stacktrace",    // 堆栈信息字段名称
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder, // 日志级别编码为大写
		EncodeTime:     zapcore.RFC3339TimeEncoder,  // 时间格式化为 RFC3339 标准
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder, // 调用栈路径简化
	}

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

	logger, err := config.Build(
		zap.AddCaller(),
		zap.AddCallerSkip(0),
		zap.AddStacktrace(zapcore.WarnLevel),
	)
	if err != nil {
		panic("failed to initialize logger: " + err.Error())
	}
	zap.ReplaceGlobals(logger)
}
