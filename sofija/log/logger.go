package log

type Logger interface {
	Debug(msg string, fields ...interface{})
	Info(msg string, fields ...interface{})
	Infof(msg string, fields ...interface{})
	Warn(msg string, fields ...interface{})
	Error(msg string, fields ...interface{})
	Errorf(msg string, fields ...interface{})
}

type NilLogger struct {
}

func NewNilLogger() *NilLogger {
	return &NilLogger{}
}

func (*NilLogger) Debug(_ string, _ ...interface{}) {
}

func (*NilLogger) Info(_ string, _ ...interface{}) {
}

func (*NilLogger) Infof(_ string, _ ...interface{}) {
}

func (*NilLogger) Warn(_ string, _ ...interface{}) {
}

func (*NilLogger) Error(_ string, _ ...interface{}) {
}

func (*NilLogger) Errorf(_ string, _ ...interface{}) {
}
