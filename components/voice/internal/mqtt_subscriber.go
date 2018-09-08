package internal

import (
	"encoding/json"

	m "github.com/eclipse/paho.mqtt.golang"
	"github.com/kihamo/boggart/components/mqtt"
	"github.com/kihamo/boggart/components/voice"
)

type SpeechRequest struct {
	Text    string  `json:"text"`
	Volume  int64   `json:"volume"`
	Speed   float64 `json:"speed"`
	Speaker string  `json:"speaker"`
}

type MQTTSubscribe struct {
	speaker voice.Component
}

func NewMQTTSubscribe(speaker voice.Component) *MQTTSubscribe {
	return &MQTTSubscribe{
		speaker: speaker,
	}
}

func (s *MQTTSubscribe) Filters() map[string]byte {
	return map[string]byte{
		voice.MQTTTopicSimpleText: 0,
		voice.MQTTTopicJSONText:   0,
	}
}

func (s *MQTTSubscribe) Callback(client mqtt.Component, message m.Message) {
	switch message.Topic() {
	case voice.MQTTTopicJSONText:
		var request SpeechRequest

		if err := json.Unmarshal(message.Payload(), &request); err == nil {
			s.speaker.SpeechWithOptions(request.Text, request.Volume, request.Speed, request.Speaker)
		}

	default:
		s.speaker.Speech(string(message.Payload()))
	}
}
