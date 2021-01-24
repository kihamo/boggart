package nut

import (
	"context"
	"errors"

	"github.com/kihamo/boggart/components/mqtt"
)

func (b *Bind) MQTTSubscribers() []mqtt.Subscriber {
	return []mqtt.Subscriber{
		mqtt.NewSubscriber(b.config.TopicVariableSet, 0, b.MQTT().WrapSubscribeDeviceIsOnline(func(_ context.Context, _ mqtt.Component, message mqtt.Message) error {
			if !b.MQTT().CheckSerialNumberInTopic(message.Topic(), 4) {
				return nil
			}

			route := message.Topic().Split()
			if len(route) < 2 {
				return errors.New("bad topic name")
			}

			return b.SetVariable(route[len(route)-2], message.String())
		})),
		mqtt.NewSubscriber(b.config.TopicCommandRun, 0, b.MQTT().WrapSubscribeDeviceIsOnline(func(ctx context.Context, _ mqtt.Component, message mqtt.Message) error {
			if !b.MQTT().CheckSerialNumberInTopic(message.Topic(), 4) {
				return nil
			}

			cmd := message.String()

			if err := b.SendCommand(cmd); err != nil {
				return err
			}

			return b.MQTT().PublishAsyncWithoutCache(ctx, b.config.TopicCommand.Format(b.Meta().SerialNumber(), cmd), "done")
		})),
	}
}
