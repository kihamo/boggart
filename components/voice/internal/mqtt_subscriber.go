package internal

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"

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
		voice.MQTTTopicSimpleText:      0,
		voice.MQTTTopicJSONText:        0,
		voice.MQTTTopicPlayerURL:       0,
		voice.MQTTTopicPlayerAction:    0,
		voice.MQTTTopicPlayerVolumeSet: 0,
	}
}

func (s *MQTTSubscribe) Callback(ctx context.Context, client mqtt.Component, message m.Message) {
	switch message.Topic() {
	case voice.MQTTTopicJSONText:
		var request SpeechRequest

		if err := json.Unmarshal(message.Payload(), &request); err == nil {
			s.player.SpeechWithOptions(ctx, request.Text, request.Volume, request.Speed, request.Speaker)
		}

	case voice.MQTTTopicSimpleText:
		s.player.Speech(ctx, string(message.Payload()))

	case voice.MQTTTopicPlayerURL:
		s.player.PlayURL(string(message.Payload()))

	case voice.MQTTTopicPlayerAction:
		switch strings.ToLower(string(message.Payload())) {
		case "stop":
			s.player.Stop()

		case "pause":
			s.player.Pause()

		case "play":
			s.player.Play()
		}

	case voice.MQTTTopicPlayerVolumeSet:
		volume, err := strconv.ParseInt(string(message.Payload()), 10, 64)
		if err == nil {
			s.player.SetVolume(volume)
		}
	}
}
