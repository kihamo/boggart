package text2speech

import (
	"os"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/di"
	"github.com/kihamo/boggart/components/mqtt"
	speechkit "github.com/kihamo/boggart/providers/yandex_speechkit_cloud"
)

type Config struct {
	di.LoggerConfig `mapstructure:",squash" yaml:",inline"`

	Debug                      bool
	CacheEnable                bool   `mapstructure:"cache_enabled" yaml:"cache_enabled"`
	CacheDirectory             string `mapstructure:"cache_directory" yaml:"cache_directory"`
	Key                        string `valid:"required"`
	Format                     string
	Quality                    string
	Language                   string
	Speaker                    string
	Emotion                    string
	Speed                      float64
	TopicGenerateBinaryOptions mqtt.Topic `mapstructure:"topic_generate_binary_options" yaml:"topic_generate_binary_options"`
	TopicGenerateBinaryText    mqtt.Topic `mapstructure:"topic_generate_binary_text" yaml:"topic_generate_binary_text"`
	TopicGenerateURLOptions    mqtt.Topic `mapstructure:"topic_generate_url_options" yaml:"topic_generate_url_options"`
	TopicGenerateURLText       mqtt.Topic `mapstructure:"topic_generate_url_text" yaml:"topic_generate_url_text"`
	TopicURL                   mqtt.Topic `mapstructure:"topic_url" yaml:"topic_url"`
	TopicBinary                mqtt.Topic `mapstructure:"topic_binary" yaml:"topic_binary"`
}

func (Type) Config() interface{} {
	var prefix mqtt.Topic = boggart.ComponentName + "/text2speech/+/"

	cacheDir, _ := os.UserCacheDir()
	if cacheDir == "" {
		cacheDir = os.TempDir()
	}

	if cacheDir != "" {
		cacheDirBind := cacheDir + string(os.PathSeparator) + boggart.ComponentName + "_text2speech"

		err := os.Mkdir(cacheDirBind, 0700)
		if err == nil || os.IsExist(err) {
			cacheDir = cacheDirBind
		}
	}

	return &Config{
		LoggerConfig: di.LoggerConfig{
			BufferedRecordsLimit: di.LoggerDefaultBufferedRecordsLimit,
			BufferedRecordsLevel: di.LoggerDefaultBufferedRecordsLevel,
		},
		CacheDirectory:             cacheDir,
		Format:                     speechkit.FormatWAV,
		Quality:                    speechkit.QualityHi,
		Language:                   speechkit.LanguageRussian,
		Speaker:                    speechkit.SpeakerOksana,
		Emotion:                    speechkit.EmotionNeutral,
		Speed:                      0.9,
		TopicGenerateBinaryOptions: prefix + "generate/binary/options",
		TopicGenerateBinaryText:    prefix + "generate/binary/text",
		TopicGenerateURLOptions:    prefix + "generate/url/options",
		TopicGenerateURLText:       prefix + "generate/url/text",
		TopicURL:                   prefix + "url",
		TopicBinary:                prefix + "binary",
	}
}
