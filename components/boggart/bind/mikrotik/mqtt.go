package mikrotik

import (
	"context"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/mqtt"
)

const (
	MQTTPublishTopicWiFiMACState         mqtt.Topic = boggart.ComponentName + "/router/+/wifi/clients/+/state"
	MQTTPublishTopicWiFiConnectedMAC     mqtt.Topic = boggart.ComponentName + "/router/+/wifi/clients/last/on/mac"
	MQTTPublishTopicWiFiDisconnectedMAC  mqtt.Topic = boggart.ComponentName + "/router/+/wifi/clients/last/off/mac"
	MQTTPublishTopicVPNLoginState        mqtt.Topic = boggart.ComponentName + "/router/+/vpn/clients/+/state"
	MQTTPublishTopicVPNConnectedLogin    mqtt.Topic = boggart.ComponentName + "/router/+/vpn/clients/last/on/login"
	MQTTPublishTopicVPNDisconnectedLogin mqtt.Topic = boggart.ComponentName + "/router/+/vpn/clients/last/off/login"
)

func (b *Bind) MQTTPublishes() []mqtt.Topic {
	return []mqtt.Topic{
		MQTTPublishTopicWiFiMACState,
		MQTTPublishTopicWiFiConnectedMAC,
		MQTTPublishTopicWiFiDisconnectedMAC,
		MQTTPublishTopicVPNLoginState,
		MQTTPublishTopicVPNConnectedLogin,
		MQTTPublishTopicVPNDisconnectedLogin,
	}
}

func (b *Bind) MQTTSubscribers() []mqtt.Subscriber {
	return []mqtt.Subscriber{
		mqtt.NewSubscriber(MQTTPublishTopicWiFiMACState.String(), 0, b.callbackMQTTWiFiSync),
	}
}

func (b *Bind) callbackMQTTWiFiSync(ctx context.Context, client mqtt.Component, message mqtt.Message) error {
	sn := b.SerialNumberWait()
	if sn == "" || !boggart.CheckSerialNumberInMQTTTopic(b, message.Topic(), 5) {
		return nil
	}

	parts := mqtt.RouteSplit(message.Topic())
	key := parts[len(parts)-2]

	// проверяем наличие в списке, дождавшись первоначальной загрузки
	value, ok := b.clientWiFi.LoadWait(key)

	topic := MQTTPublishTopicWiFiMACState.Format(sn, key)

	// если в списке нет, значит удаляем из mqtt (если там не удалено)
	if !ok {
		if message.IsTrue() {
			return b.MQTTPublishAsyncRaw(ctx, topic, 1, true, "")
		}

		return nil
	}

	// если state отличается от того что в списке, значит отправляем тот что из списка, так как там мастер данные
	state := value.(bool)
	if state != message.Bool() {
		return b.MQTTPublishAsync(ctx, topic, state)
	}

	return nil
}
