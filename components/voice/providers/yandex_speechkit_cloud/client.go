package yandex_speechkit_cloud

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	h "net/http"
	"net/url"
	"strings"

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

func (c *YandexSpeechKitCloud) GenerateURL(ctx context.Context, text, lang, speaker, emotion, format, quality string, speed float64) (string, error) {
	u := *(c.baseURL)
	u.Path = "generate"

	values := u.Query()
	values.Add("text", text)

	format = strings.ToLower(format)
	switch format {
	case FormatMP3, FormatOPUS:
		values.Add("format", format)

	case FormatWAV:
		quality = strings.ToLower(quality)
		switch quality {
		case QualityHi, QualityLo:
			values.Add("quality", quality)

		default:
			return "", errors.New("unknown quality " + quality)
		}

		values.Add("format", format)

	default:
		return "", errors.New("unknown format " + format)
	}

	lang = strings.ToLower(lang)
	switch lang {
	case LanguageEnglish, LanguageRussian, LanguageTurkish, LanguageUkrainian:
		values.Add("lang", lang)

	default:
		return "", errors.New("unknown language " + lang)
	}

	speaker = strings.ToLower(speaker)
	switch speaker {
	case SpeakerAlyss, SpeakerErmil, SpeakerJane, SpeakerOksana, SpeakerOmazh, SpeakerZahar:
		values.Add("speaker", speaker)

	default:
		return "", errors.New("unknown speaker " + speaker)
	}

	emotion = strings.ToLower(emotion)
	switch emotion {
	case EmotionNeutral, EmotionEvil, EmotionGood:
		values.Add("emotion", emotion)

	default:
		return "", errors.New("unknown emotion " + emotion)
	}

	if speed > 1 || speed < 0 {
		return "", fmt.Errorf("value of speed is wrong %.2f", speed)
	}

	values.Add("speed", fmt.Sprintf("%.1f", speed))

	u.RawQuery = values.Encode()

	return u.String(), nil
}

func (c *YandexSpeechKitCloud) Generate(ctx context.Context, text, lang, speaker, emotion, format, quality string, speed float64) (content io.Reader, err error) {
	span, ctx := tracing.StartSpanFromContext(ctx, ComponentName, "generate")
	defer span.Finish()

	u, err := c.GenerateURL(ctx, text, lang, speaker, emotion, format, quality, speed)
	if err == nil {
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
	}

	if err != nil {
		tracing.SpanError(span, err)
	}

	return content, err
}
