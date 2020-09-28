package timelapse

import (
	"context"

	"github.com/kihamo/boggart/components/mqtt"
)

func (b *Bind) MQTTSubscribers() []mqtt.Subscriber {
	return []mqtt.Subscriber{
		mqtt.NewSubscriber(b.config.TopicCapture.Format(b.Meta().ID()), 0, func(ctx context.Context, _ mqtt.Component, message mqtt.Message) error {
			return b.Capture(ctx, nil)
		}),
	}
}
