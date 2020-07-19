package telegram

import (
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/di"
	"github.com/kihamo/boggart/components/mqtt"
)

type Config struct {
	di.ProbesConfig `mapstructure:",squash" yaml:",inline"`
	di.LoggerConfig `mapstructure:",squash" yaml:",inline"`

	Token               string
	Debug               bool
	UpdatesEnabled      bool       `mapstructure:"updates_enabled" yaml:"updates_enabled"`
	UpdatesBuffer       int        `mapstructure:"updates_buffer" yaml:"updates_buffer"`
	UpdatesTimeout      int        `mapstructure:"updates_timeout" yaml:"updates_timeout"`
	TopicSendMessage    mqtt.Topic `mapstructure:"topic_send_message" yaml:"topic_send_message"`
	TopicSendFile       mqtt.Topic `mapstructure:"topic_send_file" yaml:"topic_send_file"`
	TopicSendFileURL    mqtt.Topic `mapstructure:"topic_send_file_url" yaml:"topic_send_file_url"`
	TopicSendFileBase64 mqtt.Topic `mapstructure:"topic_send_file_base64" yaml:"topic_send_file_base64"`
	TopicReceiveMessage mqtt.Topic `mapstructure:"topic_receive_message" yaml:"topic_receive_message"`
	TopicReceiveAudio   mqtt.Topic `mapstructure:"topic_receive_audio" yaml:"topic_receive_audio"`
	TopicReceiveVoice   mqtt.Topic `mapstructure:"topic_receive_voice" yaml:"topic_receive_voice"`
}

func (t Type) Config() interface{} {
	var prefix mqtt.Topic = boggart.ComponentName + "/telegram/+/"

	return &Config{
		ProbesConfig: di.ProbesConfig{
			ReadinessPeriod:  time.Minute,
			ReadinessTimeout: time.Second * 10,
			LivenessPeriod:   time.Second * 30,
			LivenessTimeout:  time.Second * 10,
		},
		LoggerConfig: di.LoggerConfig{
			BufferedRecordsLimit: di.LoggerDefaultBufferedRecordsLimit,
			BufferedRecordsLevel: di.LoggerDefaultBufferedRecordsLevel,
		},
		Debug:               false,
		UpdatesEnabled:      false,
		UpdatesBuffer:       100,
		UpdatesTimeout:      60,
		TopicSendMessage:    prefix + "send/+/message",
		TopicSendFile:       prefix + "send/+/file",
		TopicSendFileURL:    prefix + "send/+/file/url",
		TopicSendFileBase64: prefix + "send/+/file/base64",
		TopicReceiveMessage: prefix + "receive/+/message",
		TopicReceiveAudio:   prefix + "receive/+/audio",
		TopicReceiveVoice:   prefix + "receive/+/voice",
	}
}
