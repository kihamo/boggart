package text2speech

import (
	"context"

	"github.com/kihamo/boggart/components/mqtt"
)

func (b *Bind) MQTTPublishes() []mqtt.Topic {
	return []mqtt.Topic{
		b.config.TopicURL,
		b.config.TopicBinary,
	}
}

func (b *Bind) MQTTSubscribers() []mqtt.Subscriber {
	return []mqtt.Subscriber{
		mqtt.NewSubscriber(b.config.TopicGenerateBinaryOptions, 0, b.callbackMQTTGenerateBinaryOptions(true)),
		mqtt.NewSubscriber(b.config.TopicGenerateBinaryText, 0, b.callbackMQTTGenerateBinaryText(true)),
		mqtt.NewSubscriber(b.config.TopicGenerateURLOptions, 0, b.callbackMQTTGenerateBinaryOptions(false)),
		mqtt.NewSubscriber(b.config.TopicGenerateURLText, 0, b.callbackMQTTGenerateBinaryText(false)),
	}
}

func (b *Bind) callbackMQTTGenerateBinaryOptions(binary bool) func(ctx context.Context, _ mqtt.Component, message mqtt.Message) error {
	return func(ctx context.Context, _ mqtt.Component, message mqtt.Message) (err error) {
		var r struct {
			Text     string  `json:"text"`
			Format   string  `json:"format"`
			Quality  string  `json:"quality"`
			Language string  `json:"language"`
			Speaker  string  `json:"speaker"`
			Emotion  string  `json:"emotion"`
			Speed    float64 `json:"speed"`
			Force    bool    `json:"force"`
		}

		err = message.UnmarshalJSON(&r)
		if err != nil {
			return err
		}

		reader, err := b.Generate(ctx, r.Text, r.Format, r.Quality, r.Language, r.Speaker, r.Emotion, r.Speed, r.Force)
		if err != nil {
			return err
		}

		if binary {
			return b.MQTT().PublishAsync(ctx, b.config.TopicBinary, reader)
		}

		u := b.GenerateURL(ctx, r.Text, r.Format, r.Quality, r.Language, r.Speaker, r.Emotion, r.Speed, r.Force)
		if u == nil {
			return nil
		}

		return b.MQTT().PublishAsync(ctx, b.config.TopicURL, u.String())
	}
}

func (b *Bind) callbackMQTTGenerateBinaryText(binary bool) func(ctx context.Context, _ mqtt.Component, message mqtt.Message) error {
	return func(ctx context.Context, _ mqtt.Component, message mqtt.Message) (err error) {
		reader, err := b.Generate(ctx, message.String(), "", "", "", "", "", 0, false)
		if err != nil {
			return err
		}

		if binary {
			return b.MQTT().PublishAsync(ctx, b.config.TopicBinary, reader)
		}

		u := b.GenerateURL(ctx, message.String(), "", "", "", "", "", 0, false)
		if u == nil {
			return nil
		}

		return b.MQTT().PublishAsync(ctx, b.config.TopicURL, u.String())
	}
}
