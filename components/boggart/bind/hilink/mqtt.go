package hilink

import (
	"context"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/providers/hilink/client/device"
	"github.com/kihamo/boggart/components/mqtt"
)

const (
	MQTTSubscribeTopicUSSDSend   mqtt.Topic = boggart.ComponentName + "/hilink/+/ussd/send"
	MQTTSubscribeTopicUSSDResult mqtt.Topic = boggart.ComponentName + "/hilink/+/ussd"
	MQTTSubscribeTopicReboot     mqtt.Topic = boggart.ComponentName + "/hilink/+/reboot"
	MQTTPublishTopicSMS          mqtt.Topic = boggart.ComponentName + "/hilink/+/sms"
	MQTTPublishTopicBalance      mqtt.Topic = boggart.ComponentName + "/hilink/+/balance"
)

func (b *Bind) MQTTPublishes() []mqtt.Topic {
	return []mqtt.Topic{
		MQTTPublishTopicSMS,
		MQTTPublishTopicBalance,
	}
}

func (b *Bind) MQTTSubscribers() []mqtt.Subscriber {
	return []mqtt.Subscriber{
		mqtt.NewSubscriber(MQTTSubscribeTopicUSSDSend.String(), 0, boggart.WrapMQTTSubscribeDeviceIsOnline(b.Status, b.callbackMQTTUSSDSend)),
		mqtt.NewSubscriber(MQTTSubscribeTopicReboot.String(), 0, boggart.WrapMQTTSubscribeDeviceIsOnline(b.Status, b.callbackMQTTReboot)),
	}
}

func (b *Bind) callbackMQTTUSSDSend(ctx context.Context, _ mqtt.Component, message mqtt.Message) error {
	if !boggart.CheckSerialNumberInMQTTTopic(b, message.Topic(), 3) {
		return nil
	}

	content, err := b.USSD(ctx, message.String())
	if err != nil {
		return err
	}

	return b.MQTTPublishAsync(ctx, MQTTSubscribeTopicUSSDResult.Format(mqtt.NameReplace(b.SerialNumber())), content)
}

func (b *Bind) callbackMQTTReboot(ctx context.Context, _ mqtt.Component, message mqtt.Message) error {
	if message.IsFalse() || !boggart.CheckSerialNumberInMQTTTopic(b, message.Topic(), 2) {
		return nil
	}

	params := device.NewDeviceControlParamsWithContext(ctx)
	params.Request.Control = 1

	_, err := b.client.Device.DeviceControl(params)
	return err
}
