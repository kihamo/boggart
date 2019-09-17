package esphome

import (
	"context"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/mqtt"
	"github.com/kihamo/boggart/providers/esphome/native_api"
)

func (b *Bind) MQTTPublishes() []mqtt.Topic {
	return []mqtt.Topic{
		b.config.TopicState,
	}
}

func (b *Bind) MQTTSubscribers() []mqtt.Subscriber {
	return []mqtt.Subscriber{
		mqtt.NewSubscriber(b.config.TopicStateSet, 0, boggart.WrapMQTTSubscribeDeviceIsOnline(b.Status, func(ctx context.Context, _ mqtt.Component, message mqtt.Message) error {
			if !boggart.CheckSerialNumberInMQTTTopic(b, message.Topic(), 4) {
				return nil
			}

			parts := message.Topic().Split()

			entity, err := b.EntityByObjectID(ctx, parts[len(parts)-3])
			if err != nil {
				return err
			}

			switch native_api.EntityType(entity) {
			case native_api.EntityTypeFan:
				return b.provider.FanCommand(ctx, &native_api.FanCommandRequest{
					Key:      entity.(native_api.MessageEntity).GetKey(),
					HasState: true,
					State:    message.Bool(),
				})

			case native_api.EntityTypeLight:
				return b.provider.LightCommand(ctx, &native_api.LightCommandRequest{
					Key:      entity.(native_api.MessageEntity).GetKey(),
					HasState: true,
					State:    message.Bool(),
				})

			case native_api.EntityTypeSwitch:
				return b.provider.SwitchCommand(ctx, &native_api.SwitchCommandRequest{
					Key:   entity.(native_api.MessageEntity).GetKey(),
					State: message.Bool(),
				})
			}

			return nil
		})),
	}
}
