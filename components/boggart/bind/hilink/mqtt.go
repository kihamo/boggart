package hilink

import (
	"context"

	"github.com/kihamo/boggart/components/mqtt"
	"github.com/kihamo/boggart/providers/hilink/client/device"
)

func (b *Bind) MQTTPublishes() []mqtt.Topic {
	return []mqtt.Topic{
		b.config.TopicSMS,
		b.config.TopicSMSUnread,
		b.config.TopicSMSInbox,
		b.config.TopicBalance,
		b.config.TopicOperator,
		b.config.TopicLimitInternetTraffic,
		b.config.TopicSignalRSSI,
		b.config.TopicSignalRSRP,
		b.config.TopicSignalRSRQ,
		b.config.TopicSignalSINR,
		b.config.TopicSignalLevel,
		b.config.TopicConnectionTime,
		b.config.TopicConnectionDownload,
		b.config.TopicConnectionUpload,
	}
}

func (b *Bind) MQTTSubscribers() []mqtt.Subscriber {
	return []mqtt.Subscriber{
		mqtt.NewSubscriber(b.config.TopicUSSDSend, 0, b.MQTTContainer().WrapSubscribeDeviceIsOnline(b.callbackMQTTUSSDSend)),
		mqtt.NewSubscriber(b.config.TopicReboot, 0, b.MQTTContainer().WrapSubscribeDeviceIsOnline(b.callbackMQTTReboot)),
	}
}

func (b *Bind) callbackMQTTUSSDSend(ctx context.Context, _ mqtt.Component, message mqtt.Message) error {
	if !b.MQTTContainer().CheckSerialNumberInTopic(message.Topic(), 3) {
		return nil
	}

	content, err := b.USSD(ctx, message.String())
	if err != nil {
		return err
	}

	return b.MQTTContainer().PublishAsync(ctx, b.config.TopicUSSDResult.Format(b.SerialNumber()), content)
}

func (b *Bind) callbackMQTTReboot(ctx context.Context, _ mqtt.Component, message mqtt.Message) error {
	if message.IsFalse() || !b.MQTTContainer().CheckSerialNumberInTopic(message.Topic(), 2) {
		return nil
	}

	params := device.NewDeviceControlParamsWithContext(ctx)
	params.Request.Control = 1

	_, err := b.client.Device.DeviceControl(params)
	return err
}
