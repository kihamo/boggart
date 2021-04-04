package text2speech

import (
	"context"

	"github.com/kihamo/boggart/components/mqtt"
)

func (b *Bind) MQTTSubscribers() []mqtt.Subscriber {
	id := b.Meta().ID()
	cfg := b.config()

	return []mqtt.Subscriber{
		mqtt.NewSubscriber(cfg.TopicGenerateBinaryOptions.Format(id), 0, b.callbackMQTTGenerateOptions(true)),
		mqtt.NewSubscriber(cfg.TopicGenerateBinaryText.Format(id), 0, b.callbackMQTTGenerateText(true)),
		mqtt.NewSubscriber(cfg.TopicGenerateURLOptions.Format(id), 0, b.callbackMQTTGenerateOptions(false)),
		mqtt.NewSubscriber(cfg.TopicGenerateURLText.Format(id), 0, b.callbackMQTTGenerateText(false)),
	}
}

func (b *Bind) callbackMQTTGenerateOptions(binary bool) func(ctx context.Context, _ mqtt.Component, message mqtt.Message) error {
	cfg := b.config()

	return func(ctx context.Context, _ mqtt.Component, message mqtt.Message) (err error) {
		var r struct {
			Text     string  `json:"text"`
			Format   string  `json:"format,omitempty"`
			Quality  string  `json:"quality,omitempty"`
			Language string  `json:"language,omitempty"`
			Speaker  string  `json:"speaker,omitempty"`
			Emotion  string  `json:"emotion,omitempty"`
			Speed    float64 `json:"speed,omitempty"`
			Force    bool    `json:"force,omitempty"`
		}

		if err = message.JSONUnmarshal(&r); err != nil {
			return err
		}

		// binary
		if binary {
			reader, err := b.Generate(ctx, r.Text, r.Format, r.Quality, r.Language, r.Speaker, r.Emotion, r.Speed, r.Force)
			if err != nil {
				return err
			}

			return b.MQTT().PublishAsyncWithoutCache(ctx, cfg.TopicResponseBinary.Format(b.Meta().ID()), reader)
		}

		// URL
		u, err := b.GenerateURL(ctx, r.Text, r.Format, r.Quality, r.Language, r.Speaker, r.Emotion, r.Speed, r.Force)
		if u == nil {
			return err
		}

		return b.MQTT().PublishAsyncWithoutCache(ctx, cfg.TopicResponseURL.Format(b.Meta().ID()), u.String())
	}
}

func (b *Bind) callbackMQTTGenerateText(binary bool) func(ctx context.Context, _ mqtt.Component, message mqtt.Message) error {
	cfg := b.config()

	return func(ctx context.Context, _ mqtt.Component, message mqtt.Message) (err error) {
		// binary
		if binary {
			reader, err := b.Generate(ctx, message.String(), "", "", "", "", "", 0, false)
			if err != nil {
				return err
			}

			return b.MQTT().PublishAsyncWithoutCache(ctx, cfg.TopicResponseBinary.Format(b.Meta().ID()), reader)
		}

		// URL
		u, err := b.GenerateURL(ctx, message.String(), "", "", "", "", "", 0, false)
		if u == nil {
			return err
		}

		return b.MQTT().PublishAsyncWithoutCache(ctx, cfg.TopicResponseURL.Format(b.Meta().ID()), u.String())
	}
}
