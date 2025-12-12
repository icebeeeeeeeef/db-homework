package startup

import (
	"webook/pkg/logger"

	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// InitLogger 初始化并返回项目抽象的日志接口实现
func InitLogger() logger.LoggerV1 {
	cfg := zap.NewProductionConfig()
	// 读取日志级别，可选: debug/info/warn/error
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
	// 时间格式更友好
	cfg.EncoderConfig.TimeKey = "ts"
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	zl, err := cfg.Build(zap.AddCaller())
	if err != nil {
		// 兜底：构建失败时使用默认
		zl = zap.NewExample()
	}
	return logger.NewZapLogger(zl)
}
