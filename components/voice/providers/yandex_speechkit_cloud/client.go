package yandex_speechkit_cloud

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	h "net/http"
	"net/url"

	"github.com/kihamo/boggart/components/boggart/protocols/http"
	tracing "github.com/kihamo/shadow/components/tracing/http"
)

const (
	ComponentName = "yandex_speechkit_cloud"

	FormatMP3  = "mp3"
	FormatWAV  = "wav"
	FormatOPUS = "opus"

	QualityHi = "hi"
	QualityLo = "lo"

	LanguageRussian   = "ru-RU"
	LanguageEnglish   = "en-US"
	LanguageUkrainian = "uk-UK"
	LanguageTurkish   = "tr-TR"

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

type YandexSpeechKitCloud struct {
	client  *http.Client
	key     string
	baseURL *url.URL
}

func NewYandexSpeechKitCloud(key string, debug bool) *YandexSpeechKitCloud {
	values := url.Values{}
	values.Add("key", key)

	return &YandexSpeechKitCloud{
		client: http.NewClient().WithDebug(debug),
		key:    key,
		baseURL: &url.URL{
			Scheme:   "https",
			Host:     "tts.voicetech.yandex.net",
			RawQuery: values.Encode(),
		},
	}
}

func (c *YandexSpeechKitCloud) Generate(ctx context.Context, text, lang, speaker, emotion, format, quality string, speed float64) ([]byte, error) {
	u := *(c.baseURL)
	u.Path = "generate"

	values := u.Query()
	values.Add("text", text)
	values.Add("format", format)

	if format == FormatWAV {
		values.Add("quality", quality)
	}

	values.Add("lang", lang)
	values.Add("speaker", speaker)
	values.Add("speed", fmt.Sprintf("%.1f", speed))
	values.Add("emotion", emotion)
	u.RawQuery = values.Encode()

	ctx = tracing.ComponentNameToContext(ctx, ComponentName)
	ctx = tracing.OperationNameToContext(ctx, ComponentName+".generate")

	response, err := c.client.Get(ctx, u.String())
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	switch response.StatusCode {
	case h.StatusOK:
		return ioutil.ReadAll(response.Body)

	case h.StatusLocked:
		return nil, errors.New("API key is locked, please contact Yandex support team")
	}

	return nil, errors.New("Returns not 200 OK response")
}
