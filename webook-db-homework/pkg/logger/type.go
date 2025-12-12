package logger

type Logger interface {
	Info(msg string, fields ...any)

	Debug(msg string, fields ...any)

	Warn(msg string, fields ...any)

	Error(msg string, fields ...any)
}

//这种的实现要求传入的参数是key-value key-value的顺序，必须是偶数，下面的则不用
//这种实现的兼容性更好，也就是限制不多
//但是V1的实现必须要求参数是有名字的

//这样就可以把zap 的logger 抽象成一个实现

type LoggerV1 interface {
	Info(msg string, fields ...Field)

	Debug(msg string, fields ...Field)

	Warn(msg string, fields ...Field)

	Error(msg string, fields ...Field)
}

type Field struct {
	Key   string
	Value any
}
