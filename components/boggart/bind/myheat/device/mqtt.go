package device

import (
	"context"

	"github.com/kihamo/boggart/components/mqtt"
	"github.com/kihamo/boggart/providers/myheat/device/client/state"
)

func (b *Bind) MQTTSubscribers() []mqtt.Subscriber {
	cfg := b.config()

	return []mqtt.Subscriber{
		mqtt.NewSubscriber(cfg.TopicSecurityArmed, 0, b.MQTT().WrapSubscribeDeviceIsOnline(b.callbackMQTTSecurityArmed)),
	}
}

func (b *Bind) callbackMQTTSecurityArmed(ctx context.Context, _ mqtt.Component, message mqtt.Message) error {
	if !b.MQTT().CheckSerialNumberInTopic(message.Topic(), -2) {
		return nil
	}

	params := state.NewSetObjStateParamsWithContext(ctx)

	if message.IsTrue() {
		params.Request.SetStateSecurityRequest.Action = "armSecurity"
	} else {
		params.Request.SetStateSecurityRequest.Action = "disarmSecurity"
	}

	_, err := b.client.State.SetObjState(params, nil)

	return err
}
