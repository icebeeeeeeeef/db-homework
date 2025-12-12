package bootstrap

import (
	"webook/pkg/logger"

	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// InitLogger 初始化并返回项目抽象的日志接口实现
func InitLogger() logger.LoggerV1 {
	cfg := zap.NewProductionConfig()
	switch viper.GetString("log.level") {
	case "debug":
		cfg.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	case "warn":
		cfg.Level = zap.NewAtomicLevelAt(zapcore.WarnLevel)
	case "error":
		cfg.Level = zap.NewAtomicLevelAt(zapcore.ErrorLevel)
	default:
		cfg.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	}
	cfg.EncoderConfig.TimeKey = "ts"
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	zl, err := cfg.Build(zap.AddCaller())
	if err != nil {
		zl = zap.NewExample()
	}
	return logger.NewZapLogger(zl)
}
