package internal

type MQTTLogger struct {
	ln func(v ...interface{})
	f  func(format string, v ...interface{})
}

func NewMQTTLogger(ln func(v ...interface{}), f func(format string, v ...interface{})) *MQTTLogger {
	return &MQTTLogger{
		ln: ln,
		f:  f,
	}
}

func (l MQTTLogger) Println(v ...interface{}) {
	l.ln(v...)
}

func (l MQTTLogger) Printf(format string, v ...interface{}) {
	l.f(format, v...)
}
