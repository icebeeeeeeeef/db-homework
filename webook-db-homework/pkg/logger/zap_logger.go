package logger

import "go.uber.org/zap"

type ZapLogger struct {
	logger *zap.Logger
}

func NewZapLogger(logger *zap.Logger) *ZapLogger {
	return &ZapLogger{logger: logger}
}

func (l *ZapLogger) Info(msg string, fields ...Field) {
	l.logger.Info(msg, l.toZapFields(fields)...)
}

func (l *ZapLogger) Debug(msg string, fields ...Field) {
	l.logger.Debug(msg, l.toZapFields(fields)...)
}

func (l *ZapLogger) Warn(msg string, fields ...Field) {
	l.logger.Warn(msg, l.toZapFields(fields)...)
}

func (l *ZapLogger) Error(msg string, fields ...Field) {
	l.logger.Error(msg, l.toZapFields(fields)...)
}
func (l *ZapLogger) toZapFields(fields []Field) []zap.Field {
	var zapFields []zap.Field
	for _, field := range fields {
		zapFields = append(zapFields, zap.Any(field.Key, field.Value))
	}
	return zapFields
}
