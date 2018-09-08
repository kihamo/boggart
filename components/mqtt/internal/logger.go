package internal

import (
	m "github.com/eclipse/paho.mqtt.golang"
)

const (
	fieldMQTTComponent = "mqtt.component"
)

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
	values := make([]interface{}, 0, len(v))
	fields := make(map[string]interface{}, 1)

	for _, value := range v {
		switch value {
		case m.NET:
			fields[fieldMQTTComponent] = "net"
		case m.PNG:
			fields[fieldMQTTComponent] = "pinger"
		case m.CLI:
			fields[fieldMQTTComponent] = "client"
		case m.DEC:
			fields[fieldMQTTComponent] = "decode"
		case m.MES:
			fields[fieldMQTTComponent] = "message"
		case m.STR:
			fields[fieldMQTTComponent] = "store"
		case m.MID:
			fields[fieldMQTTComponent] = "msgids"
		case m.TST:
			fields[fieldMQTTComponent] = "test"
		case m.STA:
			fields[fieldMQTTComponent] = "state"
		case m.ERR:
			fields[fieldMQTTComponent] = "error"

		default:
			values = append(values, value)
		}
	}

	if len(fields) > 0 {
		values = append(values, fields)
	}

	l.ln(values...)
}

func (l MQTTLogger) Printf(format string, v ...interface{}) {
	l.f(format, v...)
}
