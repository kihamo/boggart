package speechkit

import (
	"github.com/go-openapi/runtime"
	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/runtime/logger"
	"github.com/kihamo/boggart/providers/yandex_speechkit_cloud/client"
)

const (
	FormatMP3  = "mp3"
	FormatWAV  = "wav"
	FormatOPUS = "opus"

	QualityHi = "hi"
	QualityLo = "lo"

	LanguageRussian   = "ru-ru"
	LanguageEnglish   = "en-us"
	LanguageUkrainian = "uk-uk"
	LanguageTurkish   = "tr-tr"

	SpeakerJane   = "jane"
	SpeakerOksana = "oksana"
	SpeakerAlyss  = "alyss"
	SpeakerOmazh  = "omazh"
	SpeakerZahar  = "zahar"
	SpeakerErmil  = "ermil"

	EmotionGood    = "good"
	EmotionNeutral = "neutral"
	EmotionEvil    = "evil"
)

type Client struct {
	*client.YandexSpeechKitCloud
}

func New(debug bool, logger logger.Logger) *Client {
	cl := &Client{
		YandexSpeechKitCloud: client.Default,
	}

	if rt, ok := cl.YandexSpeechKitCloud.Transport.(*httptransport.Runtime); ok {
		rt.Consumers["audio/wav"] = runtime.ByteStreamConsumer()
		rt.Consumers["audio/mpeg"] = runtime.ByteStreamConsumer()
		rt.Consumers["audio/opus"] = runtime.ByteStreamConsumer()
		rt.Consumers["text/html"] = runtime.TextConsumer()

		if logger != nil {
			rt.SetLogger(logger)
		}

		rt.SetDebug(debug)
	}

	return cl
}
