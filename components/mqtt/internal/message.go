package internal

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"strconv"
	"time"

	m "github.com/eclipse/paho.mqtt.golang"
	"github.com/kihamo/boggart/components/mqtt"
)

var (
	PayloadTrue  = []byte(`true`)
	PayloadFalse = []byte(`false`)
)

type message struct {
	msg      m.Message
	datetime time.Time
}

func newMessage(msg m.Message) *message {
	return &message{
		msg:      msg,
		datetime: time.Now(),
	}
}

func (m *message) Datetime() time.Time {
	return m.datetime
}

func (m *message) Duplicate() bool {
	return m.msg.Duplicate()
}

func (m *message) Qos() byte {
	return m.msg.Qos()
}

func (m *message) Retained() bool {
	return m.msg.Retained()
}

func (m *message) Topic() mqtt.Topic {
	return mqtt.Topic(m.msg.Topic())
}

func (m *message) MessageID() uint16 {
	return m.msg.MessageID()
}

func (m *message) Payload() []byte {
	return m.msg.Payload()
}

func (m *message) Ack() {
	m.msg.Ack()
}

func (m *message) JSONUnmarshal(v interface{}) error {
	return json.Unmarshal(m.msg.Payload(), v)
}

func (m *message) Bool() bool {
	return m.IsTrue()
}

func (m *message) IsTrue() bool {
	return bytes.Equal(m.msg.Payload(), PayloadTrue)
}

func (m *message) IsFalse() bool {
	return bytes.Equal(m.msg.Payload(), PayloadFalse)
}

func (m *message) String() string {
	return string(m.msg.Payload())
}

func (m *message) Float64() float64 {
	v, _ := strconv.ParseFloat(m.String(), 64)
	return v
}

func (m *message) Base64() ([]byte, error) {
	payload := m.msg.Payload()

	buf := make([]byte, base64.StdEncoding.DecodedLen(len(payload)))
	_, err := base64.StdEncoding.Decode(buf, payload)

	return buf, err
}

func (m *message) Len() int {
	return len(m.msg.Payload())
}
