package modbus

import (
	"context"
	"github.com/kihamo/boggart/components/mqtt"
)

func (b *Bind) MQTTSubscribers() []mqtt.Subscriber {
	cfg := b.config()
	id := b.Meta().ID()

	return []mqtt.Subscriber{
		mqtt.NewSubscriber(cfg.TopicSetTemperature.Format(id), 0, func(ctx context.Context, _ mqtt.Component, message mqtt.Message) error {
			value := message.Float64()
			err := b.Provider().SetTemperature(value)

			if err == nil {
				err = b.MQTT().PublishAsync(ctx, cfg.TopicSetTemperatureState.Format(id), value)
			}

			return err
		}),
	}
}
