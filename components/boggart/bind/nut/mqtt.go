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
		mqtt.NewSubscriber(b.config.TopicCommand, 0, b.MQTT().WrapSubscribeDeviceIsOnline(func(_ context.Context, _ mqtt.Component, message mqtt.Message) error {
			if !b.MQTT().CheckSerialNumberInTopic(message.Topic(), 2) {
				return nil
			}

			return b.SendCommand(message.String())
		})),
	}
}
