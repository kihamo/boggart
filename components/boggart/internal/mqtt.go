package internal

import (
	"context"
	"errors"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/mqtt"
)

const (
	MQTTSubscribeTopicBindReloadByPayload mqtt.Topic = boggart.ComponentName + "/bind/reload"
	MQTTSubscribeTopicBindReloadByTopic   mqtt.Topic = boggart.ComponentName + "/bind/+/reload"
)

func (c *Component) MQTTSubscribers() []mqtt.Subscriber {
	return []mqtt.Subscriber{
		mqtt.NewSubscriber(MQTTSubscribeTopicBindReloadByPayload, 0, c.callbackBindReloadInPayload),
		mqtt.NewSubscriber(MQTTSubscribeTopicBindReloadByTopic, 0, c.callbackBindReloadInTopic),
	}
}

func (c *Component) callbackBindReloadInPayload(ctx context.Context, _ mqtt.Component, message mqtt.Message) error {
	id := message.String()
	if len(id) < 1 {
		return errors.New("bad bind id")
	}

	return c.ReloadConfigByID(id)
}

func (c *Component) callbackBindReloadInTopic(ctx context.Context, _ mqtt.Component, message mqtt.Message) error {
	route := message.Topic().Split()
	if len(route) < 1 {
		return errors.New("bad bind id")
	}

	return c.ReloadConfigByID(route[len(route)-2])
}
