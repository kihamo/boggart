package yandex_speechkit_cloud

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	h "net/http"
	"net/url"

	"github.com/kihamo/boggart/components/boggart/protocols/http"
	"github.com/kihamo/shadow/components/tracing"
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

func NewYandexSpeechKitCloud(key string) *YandexSpeechKitCloud {
	values := url.Values{}
	values.Add("key", key)

	return &YandexSpeechKitCloud{
		client: http.NewClient(),
		key:    key,
		baseURL: &url.URL{
			Scheme:   "https",
			Host:     "tts.voicetech.yandex.net",
			RawQuery: values.Encode(),
		},
	}
}

func (c *YandexSpeechKitCloud) WithDebug(debug bool) *YandexSpeechKitCloud {
	c.client.WithDebug(debug)
	return c
}

func (c *YandexSpeechKitCloud) GenerateURL(ctx context.Context, text, lang, speaker, emotion, format, quality string, speed float64) string {
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

	return u.String()
}

func (c *YandexSpeechKitCloud) Generate(ctx context.Context, text, lang, speaker, emotion, format, quality string, speed float64) (content io.Reader, err error) {
	span, ctx := tracing.StartSpanFromContext(ctx, ComponentName, "generate")
	defer span.Finish()

	u := c.GenerateURL(ctx, text, lang, speaker, emotion, format, quality, speed)

	response, err := c.client.Get(ctx, u)
	if err == nil {
		defer response.Body.Close()

		switch response.StatusCode {
		case h.StatusOK:
			b := &bytes.Buffer{}
			_, err = io.Copy(b, response.Body)
			content = b

		case h.StatusLocked:
			err = errors.New("API key is locked, please contact Yandex support team")

		default:
			err = errors.New("returns not 200 OK response")
		}
	}

	if err != nil {
		tracing.SpanError(span, err)
	}

	return content, err
}
