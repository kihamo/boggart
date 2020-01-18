package text2speech

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/di"
	speechkit "github.com/kihamo/boggart/providers/yandex_speechkit_cloud"
	"github.com/kihamo/boggart/providers/yandex_speechkit_cloud/client/generate"
)

type Bind struct {
	di.ConfigBind
	di.MetaBind
	di.MQTTBind
	di.LoggerBind

	config   *Config
	provider *speechkit.Client
}

const (
	cacheKeySeparator = string(os.PathSeparator)
)

var messageReplacer = strings.NewReplacer(
	":", "",
	"[", "",
	"]", "",
	"-", "",
	";", "",
	")", "",
	"(", "",
	"\"", "",
	"'", "",
)

func trim(message string) string {
	message = strings.TrimSpace(message)
	message = strings.ToLower(message)
	message = strings.TrimFunc(message, func(r rune) bool {
		switch r {
		case '.', '!', '?', ',', ';', ':', '-', '(', ')', '"', '\'', '[', ']':
			return true
		}

		return false
	})
	message = messageReplacer.Replace(message)
	message = strings.Join(strings.Fields(message), " ")

	return message
}

func (b *Bind) GenerateURL(ctx context.Context, text, format, quality, language, speaker, emotion string, speed float64, force bool) (*url.URL, error) {
	externalURL := b.Config().App().String(boggart.ConfigExternalURL)
	if externalURL == "" {
		return nil, errors.New("config external URL ins't set")
	}

	u, err := url.Parse(externalURL)
	if err != nil {
		return nil, err
	}

	u.Path = "/" + boggart.ComponentName + "/widget/" + b.Meta().ID()

	values := u.Query()
	values.Add("text", text)
	values.Add("speed", strconv.FormatFloat(speed, 'f', -1, 64))
	values.Add("force", strconv.FormatBool(force))
	values.Add("language", language)
	values.Add("speaker", speaker)
	values.Add("emotion", emotion)
	values.Add("format", format)
	values.Add("quality", quality)

	if keysConfig := b.Config().App().String(boggart.ConfigAccessKeys); keysConfig != "" {
		if keys := strings.Split(keysConfig, ","); len(keys) > 0 {
			values.Add(boggart.AccessKeyName, keys[0])
		}
	}

	u.RawQuery = values.Encode()

	return u, nil
}

func (b *Bind) Generate(ctx context.Context, text, format, quality, language, speaker, emotion string, speed float64, force bool) (io.Reader, error) {
	writer := bytes.NewBuffer(nil)

	err := b.GenerateWriter(ctx, text, format, quality, language, speaker, emotion, speed, force, writer)
	if err != nil {
		return nil, err
	}

	return writer, nil
}

func (b *Bind) GenerateWriter(ctx context.Context, text, format, quality, language, speaker, emotion string, speed float64, force bool, writer io.Writer) error {
	text = trim(text)
	if text == "" {
		return errors.New("text is empty")
	}

	cacheKey := md5.New()

	params := generate.NewGenerateParamsWithContext(ctx).
		WithKey(b.config.Key).
		WithText(text)

	format = strings.ToLower(format)
	if format == "" {
		format = b.config.Format
	}

	switch format {
	case speechkit.FormatMP3, speechkit.FormatOPUS:
		params.SetFormat(&format)
		cacheKey.Write([]byte(cacheKeySeparator + format))

	case speechkit.FormatWAV:
		quality = strings.ToLower(quality)
		if quality == "" {
			quality = b.config.Quality
		}

		switch quality {
		case speechkit.QualityHi, speechkit.QualityLo:
			params.SetQuality(&quality)
			cacheKey.Write([]byte(cacheKeySeparator + quality))

		default:
			return errors.New("unknown quality " + quality)
		}

		params.SetFormat(&format)
		cacheKey.Write([]byte(cacheKeySeparator + format))

	default:
		return errors.New("unknown format " + format)
	}

	language = strings.ToLower(language)
	if language == "" {
		language = b.config.Language
	}

	switch language {
	case speechkit.LanguageEnglish, speechkit.LanguageRussian, speechkit.LanguageTurkish, speechkit.LanguageUkrainian:
		params.SetLang(&language)
		cacheKey.Write([]byte(cacheKeySeparator + language))

	default:
		return errors.New("unknown language " + language)
	}

	speaker = strings.ToLower(speaker)
	if speaker == "" {
		speaker = b.config.Speaker
	}

	switch speaker {
	case speechkit.SpeakerAlyss, speechkit.SpeakerErmil, speechkit.SpeakerJane, speechkit.SpeakerOksana, speechkit.SpeakerOmazh, speechkit.SpeakerZahar:
		params.SetSpeaker(&speaker)
		cacheKey.Write([]byte(cacheKeySeparator + speaker))

	default:
		return errors.New("unknown speaker " + speaker)
	}

	emotion = strings.ToLower(emotion)
	if emotion == "" {
		emotion = b.config.Emotion
	}

	switch emotion {
	case speechkit.EmotionNeutral, speechkit.EmotionEvil, speechkit.EmotionGood:
		params.SetEmotion(&emotion)
		cacheKey.Write([]byte(cacheKeySeparator + emotion))

	default:
		return errors.New("unknown emotion " + emotion)
	}

	if speed == 0 {
		speed = b.config.Speed
	}

	if speed > 1 || speed < 0.1 {
		return fmt.Errorf("value of speed is wrong %.2f", speed)
	}

	params.SetSpeed(&speed)
	cacheKey.Write([]byte(cacheKeySeparator + strconv.FormatFloat(speed, 'f', -1, 64)))
	cacheKey.Write([]byte(cacheKeySeparator + text))

	var wrapBuffer io.Writer
	fileName := b.config.CacheDirectory + string(os.PathSeparator) + hex.EncodeToString(cacheKey.Sum(nil))

	// cache
	if b.config.CacheEnable && !force {
		if _, err := os.Stat(fileName); err == nil {
			if f, err := os.Open(fileName); err == nil {
				_, err = io.Copy(writer, f)
			}

			return err
		}
	}

	if b.config.CacheEnable || force {
		wrapBuffer = bytes.NewBuffer(nil)
	} else {
		wrapBuffer = writer
	}

	_, err := b.provider.Generate.Generate(params, wrapBuffer)

	if err != nil {
		switch v := err.(type) {
		case *generate.GenerateBadRequest:
			err = errors.New(v.GetPayload())
		case *generate.GenerateLocked:
			err = errors.New(v.GetPayload())
		case *generate.GenerateDefault:
			err = errors.New(v.GetPayload())
		}

		return err
	}

	if b.config.CacheEnable {
		f, err := os.Create(fileName)
		if err != nil {
			return err
		}

		_, err = io.Copy(io.MultiWriter(f, writer), wrapBuffer.(io.Reader))
	}

	return err
}
