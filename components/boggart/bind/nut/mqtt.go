package nut

import (
	"context"
	"errors"

	"github.com/kihamo/boggart/components/mqtt"
)

func (b *Bind) MQTTPublishes() []mqtt.Topic {
	return []mqtt.Topic{
		b.config.TopicVariable,
	}
}

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

			variable := route[len(route)-2]

			variables, err := b.Variables()
			if err != nil {
				return err
			}

			for _, v := range variables {
				if mqtt.NameReplace(v.Name) == variable {
					result, err := b.SetVariable(v.Name, message.String())
					if !result {
						return errors.New("nut returned not OK result")
					}

					return err
				}
			}

			return errors.New("variable name " + variable + " not found")
		})),
		mqtt.NewSubscriber(b.config.TopicCommand, 0, b.MQTT().WrapSubscribeDeviceIsOnline(func(_ context.Context, _ mqtt.Component, message mqtt.Message) error {
			if !b.MQTT().CheckSerialNumberInTopic(message.Topic(), 2) {
				return nil
			}

			result, err := b.SendCommand(message.String())
			if !result {
				return errors.New("nut returned not OK result")
			}

			return err
		})),
	}
}
