package internal

import (
	m "github.com/eclipse/paho.mqtt.golang"
	"github.com/kihamo/boggart/components/mqtt"
	"github.com/kihamo/boggart/components/voice"
)

type Speaker interface {
	Speech(text string) error
}

type MQTTSubscribe struct {
	speaker Speaker
}

func NewMQTTSubscribe(speaker Speaker) *MQTTSubscribe {
	return &MQTTSubscribe{
		speaker: speaker,
	}
}

func (s *MQTTSubscribe) Filters() map[string]byte {
	return map[string]byte{
		voice.MQTTTopicSimpleText: 0,
	}
}

func (s *MQTTSubscribe) Callback(client mqtt.Component, message m.Message) {
	s.speaker.Speech(string(message.Payload()))
}
