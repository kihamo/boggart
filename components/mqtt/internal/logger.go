package internal

import (
	"fmt"
	"strings"

	m "github.com/eclipse/paho.mqtt.golang"
)

const (
	fieldMQTTComponent = "component"
)

type MQTTLogger struct {
	ln func(string, ...interface{})
	f  func(string, ...interface{})
}

func NewMQTTLogger(ln func(string, ...interface{}), f func(string, ...interface{})) *MQTTLogger {
	return &MQTTLogger{
		ln: ln,
		f:  f,
	}
}

func (l MQTTLogger) Println(v ...interface{}) {
	if len(v) == 0 {
		return
	}

	fields := make([]interface{}, 2)
	fields[0] = fieldMQTTComponent

	switch v[0] {
	case m.NET:
		fields[1] = "net"
	case m.PNG:
		fields[1] = "pinger"
	case m.CLI:
		fields[1] = "client"
	case m.DEC:
		fields[1] = "decode"
	case m.MES:
		fields[1] = "message"
	case m.STR:
		fields[1] = "store"
	case m.MID:
		fields[1] = "msgids"
	case m.TST:
		fields[1] = "test"
	case m.STA:
		fields[1] = "state"
	case m.ERR:
		fields[1] = "error"
	}

	if len(fields) == 2 {
		v = v[1:]
	}

	var msg strings.Builder

	if len(v) > 0 {
		for i, value := range v {
			if i != 0 {
				msg.WriteString(" ")
			}

			fmt.Fprintf(&msg, "%v", value)
		}
	}

	l.ln(msg.String(), fields...)
}

func (l MQTTLogger) Printf(format string, v ...interface{}) {
	l.f(format, v...)
}
