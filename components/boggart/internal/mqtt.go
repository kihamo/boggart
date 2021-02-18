package internal

import (
	"bytes"
	"context"
	"errors"

	"github.com/kihamo/boggart/components/boggart/di"
	"github.com/kihamo/boggart/components/mqtt"
)

var (
	payloadReload = []byte(`reload`)
	payloadDone   = []byte(`done`)
)

// TODO: run probes of bind

func (c *Component) callbackBindReloadInPayload(_ context.Context, _ mqtt.Component, message mqtt.Message) error {
	id := message.String()
	if len(id) < 1 {
		return errors.New("bad bind id")
	}

	return c.ReloadConfigByID(id)
}

func (c *Component) callbackBindReloadInTopic(ctx context.Context, _ mqtt.Component, message mqtt.Message) error {
	if !bytes.Equal(message.Payload(), payloadReload) {
		return nil
	}

	topic := message.Topic()

	route := topic.Split()
	if len(route) < 1 {
		return errors.New("bad bind id")
	}

	bindID := route[len(route)-2]

	if err := c.ReloadConfigByID(bindID); err != nil {
		return err
	}

	item := c.Bind(bindID)
	if item == nil {
		return errors.New("bind " + bindID + "isn't registered")
	}

	if bindSupport, ok := di.MQTTContainerBind(item.Bind()); ok && bindSupport != nil {
		return bindSupport.PublishAsyncWithoutCache(ctx, topic, payloadDone)
	}

	c.mqtt.PublishAsyncWithoutCache(ctx, topic, 1, false, payloadDone)
	return nil
}
