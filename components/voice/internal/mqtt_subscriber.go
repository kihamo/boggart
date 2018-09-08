package internal

import (
	"bytes"
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
	player voice.Component
}

func NewMQTTSubscribe(player voice.Component) *MQTTSubscribe {
	return &MQTTSubscribe{
		player: player,
	}
}

func (s *MQTTSubscribe) Filters() map[string]byte {
	return map[string]byte{
		voice.MQTTTopicSimpleText:  0,
		voice.MQTTTopicJSONText:    0,
		voice.MQTTTopicPlayerURL:   0,
		voice.MQTTTopicPlayerPause: 0,
		voice.MQTTTopicPlayerStop:  0,
		voice.MQTTTopicPlayerPlay:  0,
	}
}

func (s *MQTTSubscribe) Callback(client mqtt.Component, message m.Message) {
	switch message.Topic() {
	case voice.MQTTTopicJSONText:
		var request SpeechRequest

		if err := json.Unmarshal(message.Payload(), &request); err == nil {
			s.player.SpeechWithOptions(request.Text, request.Volume, request.Speed, request.Speaker)
		}

	case voice.MQTTTopicSimpleText:
		s.player.Speech(string(message.Payload()))

	case voice.MQTTTopicPlayerURL:
		s.player.PlayURL(string(message.Payload()))

	case voice.MQTTTopicPlayerPause:
		if bytes.Compare(message.Payload(), []byte("1")) == 0 {
			s.player.Pause()
		}

	case voice.MQTTTopicPlayerStop:
		if bytes.Compare(message.Payload(), []byte("1")) == 0 {
			s.player.Stop()
		}

	case voice.MQTTTopicPlayerPlay:
		if bytes.Compare(message.Payload(), []byte("1")) == 0 {
			s.player.Play()
		}
	}
}
