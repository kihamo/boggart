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
				// устанавливаемое значение всегда кратно 0.5 и округляется в меньшую сторону
				// даже на устройстве шаг 0.5, поэтому принудительно округляем
				val := int(value * 10)
				val -= val % 5
				value = float64(val) / 10

				err = b.MQTT().PublishAsync(ctx, cfg.TopicSetTemperatureState.Format(id), value)
			}

			return err
		}),
	}
}
