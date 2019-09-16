package hilink

import (
	"context"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/mqtt"
	"github.com/kihamo/boggart/providers/hilink/client/device"
)

func (b *Bind) MQTTPublishes() []mqtt.Topic {
	return []mqtt.Topic{
		b.config.TopicSMS,
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
		mqtt.NewSubscriber(b.config.TopicUSSDSend, 0, boggart.WrapMQTTSubscribeDeviceIsOnline(b.Status, b.callbackMQTTUSSDSend)),
		mqtt.NewSubscriber(b.config.TopicReboot, 0, boggart.WrapMQTTSubscribeDeviceIsOnline(b.Status, b.callbackMQTTReboot)),
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

	return b.MQTTPublishAsync(ctx, b.config.TopicUSSDResult.Format(b.SerialNumber()), content)
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
