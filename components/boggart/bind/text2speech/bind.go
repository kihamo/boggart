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

	"github.com/kihamo/boggart/components/boggart/di"
	"github.com/kihamo/boggart/protocols/swagger"
	"github.com/kihamo/boggart/providers/yandex_speechkit_cloud"
	"github.com/kihamo/boggart/providers/yandex_speechkit_cloud/client/generate"
)

type Bind struct {
	di.ConfigBind
	di.LoggerBind
	di.MetaBind
	di.MQTTBind
	di.WidgetBind

	provider *speechkit.Client
}

const (
	separator = string(os.PathSeparator)
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

func (b *Bind) config() *Config {
	return b.Config().Bind().(*Config)
}

func (b *Bind) Run() error {
	b.provider = speechkit.New(b.config().Debug, swagger.NewLogger(
		func(message string) {
			b.Logger().Info(message)
		},
		func(message string) {
			b.Logger().Debug(message)
		}))

	return nil
}

func (b *Bind) GenerateURL(ctx context.Context, text, format, quality, language, speaker, emotion string, speed float64, force bool) (*url.URL, error) {
	vs := map[string]string{
		"text": text,
	}
	if force {
		vs["force"] = "1"
	}

	if language != "" {
		vs["language"] = language
	}

	if speaker != "" {
		vs["speaker"] = speaker
	}

	if emotion != "" {
		vs["emotion"] = emotion
	}

	if format != "" {
		vs["format"] = format
	}

	if quality != "" {
		vs["quality"] = quality
	}

	if speed > 0 {
		vs["quality"] = strconv.FormatFloat(speed, 'f', -1, 64)
	}

	return b.Widget().URL(vs)
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

	cfg := b.config()
	cacheKey := md5.New()

	params := generate.NewGenerateParamsWithContext(ctx).
		WithKey(cfg.Key).
		WithText(text)

	format = strings.ToLower(format)
	if format == "" {
		format = cfg.Format
	}

	switch format {
	case speechkit.FormatMP3, speechkit.FormatOPUS:
		params.SetFormat(&format)
		cacheKey.Write([]byte(separator + format))

	case speechkit.FormatWAV:
		quality = strings.ToLower(quality)
		if quality == "" {
			quality = cfg.Quality
		}

		switch quality {
		case speechkit.QualityHi, speechkit.QualityLo:
			params.SetQuality(&quality)
			cacheKey.Write([]byte(separator + quality))

		default:
			return errors.New("unknown quality " + quality)
		}

		params.SetFormat(&format)
		cacheKey.Write([]byte(separator + format))

	default:
		return errors.New("unknown format " + format)
	}

	language = strings.ToLower(language)
	if language == "" {
		language = cfg.Language
	}

	switch language {
	case speechkit.LanguageEnglish, speechkit.LanguageRussian, speechkit.LanguageTurkish, speechkit.LanguageUkrainian:
		params.SetLang(&language)
		cacheKey.Write([]byte(separator + language))

	default:
		return errors.New("unknown language " + language)
	}

	speaker = strings.ToLower(speaker)
	if speaker == "" {
		speaker = cfg.Speaker
	}

	switch speaker {
	case speechkit.SpeakerAlyss, speechkit.SpeakerErmil, speechkit.SpeakerJane, speechkit.SpeakerOksana, speechkit.SpeakerOmazh, speechkit.SpeakerZahar:
		params.SetSpeaker(&speaker)
		cacheKey.Write([]byte(separator + speaker))

	default:
		return errors.New("unknown speaker " + speaker)
	}

	emotion = strings.ToLower(emotion)
	if emotion == "" {
		emotion = cfg.Emotion
	}

	switch emotion {
	case speechkit.EmotionNeutral, speechkit.EmotionEvil, speechkit.EmotionGood:
		params.SetEmotion(&emotion)
		cacheKey.Write([]byte(separator + emotion))

	default:
		return errors.New("unknown emotion " + emotion)
	}

	if speed == 0 {
		speed = cfg.Speed
	}

	if speed > 1 || speed < 0.1 {
		return fmt.Errorf("value of speed is wrong %.2f", speed)
	}

	params.SetSpeed(&speed)
	cacheKey.Write([]byte(separator + strconv.FormatFloat(speed, 'f', -1, 64)))
	cacheKey.Write([]byte(separator + text))

	var wrapBuffer io.Writer

	fileName := cfg.CacheDirectory + separator + hex.EncodeToString(cacheKey.Sum(nil))

	// cache
	if cfg.CacheEnable && !force {
		if _, err := os.Stat(fileName); err == nil {
			if f, e := os.Open(fileName); e == nil {
				_, err = io.Copy(writer, f)
			} else {
				err = e
			}

			return err
		}
	}

	if cfg.CacheEnable || force {
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

	if cfg.CacheEnable {
		f, e := os.Create(fileName)
		if e != nil {
			return e
		}

		_, err = io.Copy(io.MultiWriter(f, writer), wrapBuffer.(io.Reader))
	}

	return err
}
