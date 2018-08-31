package internal

import (
	"github.com/kihamo/boggart/components/voice"
	yandex "github.com/kihamo/boggart/components/voice/providers/yandex_speechkit_cloud"
	"github.com/kihamo/shadow/components/config"
)

func (c *Component) ConfigVariables() []config.Variable {
	return []config.Variable{
		config.NewVariable(voice.ConfigYandexSpeechKitCloudKey, config.ValueTypeString).
			WithUsage("API key").
			WithGroup("Yandex SpeechKit Cloud"),
		config.NewVariable(voice.ConfigYandexSpeechKitCloudFormat, config.ValueTypeString).
			WithUsage("Default format").
			WithGroup("Yandex SpeechKit Cloud").
			WithEditable(true).
			WithView([]string{config.ViewEnum}).
			WithViewOptions(map[string]interface{}{
				config.ViewOptionEnumOptions: [][]interface{}{
					{yandex.FormatMP3, "MP3"},
					{yandex.FormatWAV, "WAV"},
					{yandex.FormatOPUS, "OPUS"},
				},
			}).
			WithDefault(yandex.FormatMP3),
		config.NewVariable(voice.ConfigYandexSpeechKitCloudQuality, config.ValueTypeString).
			WithUsage("Default quality for WAV format").
			WithGroup("Yandex SpeechKit Cloud").
			WithEditable(true).
			WithView([]string{config.ViewEnum}).
			WithViewOptions(map[string]interface{}{
				config.ViewOptionEnumOptions: [][]interface{}{
					{yandex.QualityHi, "Hi"},
					{yandex.QualityLo, "Lo"},
				},
			}).
			WithDefault(yandex.QualityHi),
		config.NewVariable(voice.ConfigYandexSpeechKitCloudLanguage, config.ValueTypeString).
			WithUsage("Default language").
			WithGroup("Yandex SpeechKit Cloud").
			WithEditable(true).
			WithView([]string{config.ViewEnum}).
			WithViewOptions(map[string]interface{}{
				config.ViewOptionEnumOptions: [][]interface{}{
					{yandex.LanguageRussian, "Russian"},
					{yandex.LanguageEnglish, "English"},
					{yandex.LanguageUkrainian, "Ukrainian"},
					{yandex.LanguageTurkish, "Turkish"},
				},
			}).
			WithDefault(yandex.LanguageRussian),
		config.NewVariable(voice.ConfigYandexSpeechKitCloudSpeaker, config.ValueTypeString).
			WithUsage("Default speaker").
			WithGroup("Yandex SpeechKit Cloud").
			WithEditable(true).
			WithView([]string{config.ViewEnum}).
			WithViewOptions(map[string]interface{}{
				config.ViewOptionEnumOptions: [][]interface{}{
					{yandex.SpeakerJane, "Jane (female)"},
					{yandex.SpeakerOksana, "Oksana (female)"},
					{yandex.SpeakerAlyss, "Alyss (female)"},
					{yandex.SpeakerOmazh, "Omazh (female)"},
					{yandex.SpeakerZahar, "Zahar (male)"},
					{yandex.SpeakerErmil, "Ermil (male)"},
				},
			}).
			WithDefault(yandex.SpeakerAlyss),
		config.NewVariable(voice.ConfigYandexSpeechKitCloudSpeed, config.ValueTypeFloat64).
			WithUsage("Default speed").
			WithGroup("Yandex SpeechKit Cloud").
			WithEditable(true).
			WithDefault(1.0),
		config.NewVariable(voice.ConfigYandexSpeechKitCloudEmotion, config.ValueTypeString).
			WithUsage("Default emotion").
			WithGroup("Yandex SpeechKit Cloud").
			WithEditable(true).
			WithView([]string{config.ViewEnum}).
			WithViewOptions(map[string]interface{}{
				config.ViewOptionEnumOptions: [][]interface{}{
					{yandex.EmotionGood, "Good"},
					{yandex.EmotionNeutral, "Neutral"},
					{yandex.EmotionEvil, "Evil"},
				},
			}).
			WithDefault(yandex.EmotionNeutral),
	}
}
