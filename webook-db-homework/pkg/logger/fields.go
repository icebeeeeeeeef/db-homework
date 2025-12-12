package logger

func Error(err error) Field {
	return Field{
		Key:   "error",
		Value: err.Error(),
	}
}

func String(key string, value string) Field {
	return Field{
		Key:   key,
		Value: value,
	}
}

func Int(key string, value int) Field {
	return Field{
		Key:   key,
		Value: value,
	}
}
