package modbus

import (
	"context"

	"github.com/kihamo/boggart/components/mqtt"
)

func (b *Bind) MQTTSubscribers() []mqtt.Subscriber {
	cfg := b.config()
	id := b.Meta().ID()

	return []mqtt.Subscriber{
		mqtt.NewSubscriber(cfg.TopicPower.Format(id), 0, func(ctx context.Context, _ mqtt.Component, message mqtt.Message) error {
			value := message.Bool()
			err := b.Provider().Status(value)

			if err == nil {
				err = b.MQTT().PublishAsync(ctx, cfg.TopicPowerState.Format(id), value)
			}

			return err
		}),
		mqtt.NewSubscriber(cfg.TopicSetTemperature.Format(id), 0, func(ctx context.Context, _ mqtt.Component, message mqtt.Message) error {
			value := message.Float64()
			err := b.Provider().SetTemperature(value)

			if err == nil {
				// устанавливаемое значение всегда кратно 0.5 и округляется в меньшую сторону
				// даже на устройстве шаг 0.5, поэтому принудительно округляем
				val := int(value * 10)
				val -= val % 5
				value = float64(val) / 10

				err = b.MQTT().PublishAsync(ctx, cfg.TopicSetTemperatureState.Format(id), value)
			}

			return err
		}),
		mqtt.NewSubscriber(cfg.TopicAway.Format(id), 0, func(ctx context.Context, _ mqtt.Component, message mqtt.Message) error {
			value := message.Bool()
			err := b.Provider().Away(value)

			if err == nil {
				err = b.MQTT().PublishAsync(ctx, cfg.TopicAwayState.Format(id), value)
			}

			return err
		}),
		mqtt.NewSubscriber(cfg.TopicAwayTemperature.Format(id), 0, func(ctx context.Context, _ mqtt.Component, message mqtt.Message) error {
			value := message.Uint64()
			err := b.Provider().AwayTemperature(uint16(value))

			if err == nil {
				err = b.MQTT().PublishAsync(ctx, cfg.TopicAwayTemperatureState.Format(id), value)
			}

			return err
		}),
	}
}
