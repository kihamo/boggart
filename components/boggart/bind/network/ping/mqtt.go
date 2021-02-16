package ping

import (
	"context"

	"github.com/kihamo/boggart/components/mqtt"
)

func (b *Bind) MQTTSubscribers() []mqtt.Subscriber {
	cfg := b.config()

	return []mqtt.Subscriber{
		mqtt.NewSubscriber(cfg.TopicCheck.Format(cfg.Hostname), 0, func(ctx context.Context, _ mqtt.Component, _ mqtt.Message) error {
			return b.Check(ctx)
		}),
	}
}
