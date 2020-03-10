package hilink

import (
	"context"

	"github.com/kihamo/boggart/components/mqtt"
	"github.com/kihamo/boggart/providers/hilink/client/device"
)

func (b *Bind) MQTTSubscribers() []mqtt.Subscriber {
	return []mqtt.Subscriber{
		mqtt.NewSubscriber(b.config.TopicUSSDSend, 0, b.MQTT().WrapSubscribeDeviceIsOnline(b.callbackMQTTUSSDSend)),
		mqtt.NewSubscriber(b.config.TopicReboot, 0, b.MQTT().WrapSubscribeDeviceIsOnline(b.callbackMQTTReboot)),
	}
}

func (b *Bind) callbackMQTTUSSDSend(ctx context.Context, _ mqtt.Component, message mqtt.Message) error {
	if !b.MQTT().CheckSerialNumberInTopic(message.Topic(), 3) {
		return nil
	}

	content, err := b.USSD(ctx, message.String())
	if err != nil {
		return err
	}

	return b.MQTT().PublishAsync(ctx, b.config.TopicUSSDResult.Format(b.Meta().SerialNumber()), content)
}

func (b *Bind) callbackMQTTReboot(ctx context.Context, _ mqtt.Component, message mqtt.Message) error {
	if message.IsFalse() || !b.MQTT().CheckSerialNumberInTopic(message.Topic(), 2) {
		return nil
	}

	params := device.NewDeviceControlParamsWithContext(ctx)
	params.Request.Control = 1

	_, err := b.client.Device.DeviceControl(params)
	return err
}
