package internal

import (
	"bytes"
	"encoding/json"

	"github.com/eclipse/paho.mqtt.golang"
)

var (
	PayloadTrue  = []byte(`true`)
	PayloadFalse = []byte(`false`)
)

type message struct {
	m mqtt.Message
}

func newMessage(m mqtt.Message) *message {
	return &message{
		m: m,
	}
}

func (m *message) Duplicate() bool {
	return m.m.Duplicate()
}

func (m *message) Qos() byte {
	return m.m.Qos()
}

func (m *message) Retained() bool {
	return m.m.Retained()
}

func (m *message) Topic() string {
	return m.m.Topic()
}

func (m *message) MessageID() uint16 {
	return m.m.MessageID()
}

func (m *message) Payload() []byte {
	return m.m.Payload()
}

func (m *message) Ack() {
	m.m.Ack()
}

func (m *message) UnmarshalJSON(v interface{}) error {
	return json.Unmarshal(m.m.Payload(), v)
}

func (m *message) IsTrue() bool {
	return bytes.Equal(m.m.Payload(), PayloadTrue)
}

func (m *message) IsFalse() bool {
	return bytes.Equal(m.m.Payload(), PayloadFalse)
}

func (m *message) String() string {
	return string(m.m.Payload())
}
