package telegram

import (
	"os"
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/di"
	"github.com/kihamo/boggart/components/mqtt"
	"github.com/kihamo/boggart/types"
)

type Config struct {
	di.ProbesConfig `mapstructure:",squash" yaml:",inline"`
	di.LoggerConfig `mapstructure:",squash" yaml:",inline"`

	FileURLPrefix       types.URL `mapstructure:"file_url_prefix" yaml:"file_url_prefix"`
	Token               string
	FileDirectory       string     `mapstructure:"file_directory" yaml:"file_directory"`
	TopicSendMessage    mqtt.Topic `mapstructure:"topic_send_message" yaml:"topic_send_message"`
	TopicSendFile       mqtt.Topic `mapstructure:"topic_send_file" yaml:"topic_send_file"`
	TopicSendFileURL    mqtt.Topic `mapstructure:"topic_send_file_url" yaml:"topic_send_file_url"`
	TopicSendFileBase64 mqtt.Topic `mapstructure:"topic_send_file_base64" yaml:"topic_send_file_base64"`
	TopicReceiveMessage mqtt.Topic `mapstructure:"topic_receive_message" yaml:"topic_receive_message"`
	TopicReceiveAudio   mqtt.Topic `mapstructure:"topic_receive_audio" yaml:"topic_receive_audio"`
	TopicReceiveVoice   mqtt.Topic `mapstructure:"topic_receive_voice" yaml:"topic_receive_voice"`
	UpdatesBuffer       int        `mapstructure:"updates_buffer" yaml:"updates_buffer"`
	UpdatesTimeout      int        `mapstructure:"updates_timeout" yaml:"updates_timeout"`
	Debug               bool
	UpdatesEnabled      bool `mapstructure:"updates_enabled" yaml:"updates_enabled"`
	UseURLForSendFile   bool `mapstructure:"use_url_for_send_file" yaml:"use_url_for_send_file"`
	FileAutoClean       bool `mapstructure:"file_auto_clean" yaml:"file_auto_clean"`
}

func (t Type) Config() interface{} {
	var prefix mqtt.Topic = boggart.ComponentName + "/telegram/+/"

	cacheDir, _ := os.UserCacheDir()
	if cacheDir == "" {
		cacheDir = os.TempDir()
	}

	if cacheDir != "" {
		cacheDirBind := cacheDir + string(os.PathSeparator) + boggart.ComponentName + "_telegram"

		err := os.Mkdir(cacheDirBind, 0700)
		if err == nil || os.IsExist(err) {
			cacheDir = cacheDirBind
		}
	}

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
		UseURLForSendFile:   true,
		FileDirectory:       cacheDir,
		FileAutoClean:       true,
		TopicSendMessage:    prefix + "send/+/message",
		TopicSendFile:       prefix + "send/+/file",
		TopicSendFileURL:    prefix + "send/+/file/url",
		TopicSendFileBase64: prefix + "send/+/file/base64",
		TopicReceiveMessage: prefix + "receive/+/message",
		TopicReceiveAudio:   prefix + "receive/+/audio",
		TopicReceiveVoice:   prefix + "receive/+/voice",
	}
}
