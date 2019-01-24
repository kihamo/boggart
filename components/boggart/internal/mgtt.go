package internal

import (
	"context"
	"errors"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/mqtt"
	"github.com/kihamo/shadow/components/messengers"
)

const (
	MQTTSubscribeTopicMessenger mqtt.Topic = "messenger/+/+"
)

func (c *Component) MQTTSubscribers() []mqtt.Subscriber {
	<-c.application.ReadyComponent(c.Name())

	subscribers := make([]mqtt.Subscriber, 0, 2)

	if c.application.HasComponent(messengers.ComponentName) {
		<-c.application.ReadyComponent(messengers.ComponentName)
		cmp := c.application.GetComponent(messengers.ComponentName).(messengers.Component)

		subscribers = append(subscribers, mqtt.NewSubscriber(MQTTSubscribeTopicMessenger.String(), 0, func(_ context.Context, _ mqtt.Component, message mqtt.Message) error {
			if !c.config.Bool(boggart.ConfigMQTTMessengersEnabled) {
				return nil
			}

			parts := mqtt.RouteSplit(message.Topic())
			if len(parts) < 3 {
				return errors.New("bad topic name")
			}

			messenger := cmp.Messenger(parts[len(parts)-2])
			if messenger == nil {
				return errors.New("messenger " + parts[len(parts)-2] + " not found")
			}

			return messenger.SendMessage(parts[len(parts)-1], message.String())
		}))
	}

	return subscribers
}
