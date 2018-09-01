package yandex_speechkit_cloud

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	h "net/http"
	"net/url"

	"github.com/kihamo/boggart/components/boggart/protocols/http"
)

const (
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

	response, err := c.client.Get(context.Background(), u.String())
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != h.StatusOK {
		return nil, errors.New("Returns not 200 OK response")
	}

	return ioutil.ReadAll(response.Body)
}
