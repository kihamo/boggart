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

	fields := make([]interface{}, 0, 4)

	switch v[0] {
	case m.NET:
		fields = append(fields, fieldMQTTComponent, "net")
	case m.PNG:
		fields = append(fields, fieldMQTTComponent, "pinger")
	case m.CLI:
		fields = append(fields, fieldMQTTComponent, "client")
	case m.DEC:
		fields = append(fields, fieldMQTTComponent, "decode")
	case m.MES:
		fields = append(fields, fieldMQTTComponent, "message")
	case m.STR:
		fields = append(fields, fieldMQTTComponent, "store")
	case m.MID:
		fields = append(fields, fieldMQTTComponent, "msgids")
	case m.TST:
		fields = append(fields, fieldMQTTComponent, "test")
	case m.STA:
		fields = append(fields, fieldMQTTComponent, "state")
	case m.ERR:
		fields = append(fields, fieldMQTTComponent, "error")
	}

	if len(fields) == 2 {
		v = v[1:]
	}

	var msg strings.Builder

	if len(v) > 0 {
		for i, value := range v {
			if i != 0 {
				fmt.Fprint(&msg, " ")
			}

			fmt.Fprintf(&msg, "%v", value)
		}
	}

	l.ln(msg.String(), fields...)
}

func (l MQTTLogger) Printf(format string, v ...interface{}) {
	l.f(format, v...)
}
