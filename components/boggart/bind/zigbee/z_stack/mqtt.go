package zstack

import (
	"context"

	"github.com/kihamo/boggart/components/mqtt"
)

func (b *Bind) syncPermitJoin() {
	if sn := b.Meta().SerialNumber(); sn != "" {
		ctx := context.Background()
		cfg := b.config()

		b.MQTT().PublishAsync(ctx, cfg.TopicStatePermitJoin.Format(sn), b.client.PermitJoinEnabled())
		b.MQTT().PublishAsync(ctx, cfg.TopicStatePermitJoinDuration.Format(sn), b.permitJoinDuration())
	}
}

func (b *Bind) MQTTSubscribers() []mqtt.Subscriber {
	return []mqtt.Subscriber{
		mqtt.NewSubscriber(b.config().TopicPermitJoin, 0, b.MQTT().WrapSubscribeDeviceIsOnline(
			func(ctx context.Context, _ mqtt.Component, message mqtt.Message) (err error) {
				if !b.MQTT().CheckSerialNumberInTopic(message.Topic(), -2) {
					return nil
				}

				if message.IsTrue() {
					err = b.client.PermitJoin(ctx, b.permitJoinDuration())
				} else {
					err = b.client.PermitJoinDisable(ctx)
				}

				return err
			})),
	}
}
