package hilink

import (
	"context"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/mqtt"
	"github.com/kihamo/boggart/providers/hilink/client/device"
)

const (
	MQTTPrefix mqtt.Topic = boggart.ComponentName + "/hilink/+/"

	MQTTSubscribeTopicUSSDSend           = MQTTPrefix + "ussd/send"
	MQTTSubscribeTopicUSSDResult         = MQTTPrefix + "ussd"
	MQTTSubscribeTopicReboot             = MQTTPrefix + "reboot"
	MQTTPublishTopicSMS                  = MQTTPrefix + "sms"
	MQTTPublishTopicBalance              = MQTTPrefix + "balance"
	MQTTPublishTopicOperator             = MQTTPrefix + "operator"
	MQTTPublishTopicLimitInternetTraffic = MQTTPrefix + "limits/internet-traffic"
	MQTTPublishSignalRSSI                = MQTTPrefix + "signal/rssi"
	MQTTPublishSignalRSRP                = MQTTPrefix + "signal/rsrp"
	MQTTPublishSignalRSRQ                = MQTTPrefix + "signal/rsrq"
	MQTTPublishSignalSINR                = MQTTPrefix + "signal/sinr"
	MQTTPublishSignalLevel               = MQTTPrefix + "signal/level"
	MQTTPublishConnectionTime            = MQTTPrefix + "connection/time"
	MQTTPublishConnectionDownload        = MQTTPrefix + "connection/download"
	MQTTPublishConnectionUpload          = MQTTPrefix + "connection/upload"
)

func (b *Bind) MQTTPublishes() []mqtt.Topic {
	return []mqtt.Topic{
		MQTTPublishTopicSMS,
		MQTTPublishTopicBalance,
		MQTTPublishTopicOperator,
		MQTTPublishTopicLimitInternetTraffic,
		MQTTPublishSignalRSSI,
		MQTTPublishSignalRSRP,
		MQTTPublishSignalRSRQ,
		MQTTPublishSignalSINR,
		MQTTPublishSignalLevel,
		MQTTPublishConnectionTime,
		MQTTPublishConnectionDownload,
		MQTTPublishConnectionUpload,
	}
}

func (b *Bind) MQTTSubscribers() []mqtt.Subscriber {
	return []mqtt.Subscriber{
		mqtt.NewSubscriber(MQTTSubscribeTopicUSSDSend, 0, boggart.WrapMQTTSubscribeDeviceIsOnline(b.Status, b.callbackMQTTUSSDSend)),
		mqtt.NewSubscriber(MQTTSubscribeTopicReboot, 0, boggart.WrapMQTTSubscribeDeviceIsOnline(b.Status, b.callbackMQTTReboot)),
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
