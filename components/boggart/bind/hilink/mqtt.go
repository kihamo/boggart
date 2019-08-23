package hilink

import (
	"context"
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/providers/hilink/client/device"
	"github.com/kihamo/boggart/components/boggart/providers/hilink/client/ussd"
	"github.com/kihamo/boggart/components/boggart/providers/hilink/models"
	"github.com/kihamo/boggart/components/mqtt"
)

const (
	MQTTSubscribeTopicUSSDSend   mqtt.Topic = boggart.ComponentName + "/hilink/+/ussd/send"
	MQTTSubscribeTopicUSSDResult mqtt.Topic = boggart.ComponentName + "/hilink/+/ussd"
	MQTTSubscribeTopicReboot     mqtt.Topic = boggart.ComponentName + "/hilink/+/reboot"
	MQTTPublishTopicSMS          mqtt.Topic = boggart.ComponentName + "/hilink/+/sms"
)

func (b *Bind) MQTTPublishes() []mqtt.Topic {
	return []mqtt.Topic{
		MQTTPublishTopicSMS,
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

	if message.String() == "" {
		return nil
	}

	params := ussd.NewSendUSSDParamsWithContext(ctx).
		WithRequest(&models.USSD{
			Content: message.String(),
		})

	_, err := b.client.Ussd.SendUSSD(params)
	if err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			return nil

		default:
			response, err := b.client.Ussd.GetUSSD(ussd.NewGetUSSDParamsWithContext(ctx))
			if err == nil && response.Payload.Content != "" {
				return b.MQTTPublishAsync(
					ctx,
					MQTTSubscribeTopicUSSDResult.Format(mqtt.NameReplace(b.SerialNumber())),
					response.Payload.Content)
			}

			time.Sleep(time.Second)
		}
	}

	return nil
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
