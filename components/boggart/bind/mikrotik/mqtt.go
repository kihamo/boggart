package mikrotik

import (
	"context"

	"github.com/kihamo/boggart/components/mqtt"
)

func (b *Bind) MQTTPublishes() []mqtt.Topic {
	return []mqtt.Topic{
		b.config.TopicWiFiMACState,
		b.config.TopicWiFiConnectedMAC,
		b.config.TopicWiFiDisconnectedMAC,
		b.config.TopicVPNLoginState,
		b.config.TopicVPNConnectedLogin,
		b.config.TopicVPNDisconnectedLogin,
	}
}

func (b *Bind) MQTTSubscribers() []mqtt.Subscriber {
	return []mqtt.Subscriber{
		mqtt.NewSubscriber(b.config.TopicWiFiMACState, 0, b.callbackMQTTWiFiSync),
	}
}

func (b *Bind) callbackMQTTWiFiSync(ctx context.Context, client mqtt.Component, message mqtt.Message) error {
	sn := b.SerialNumberWait()
	if sn == "" || !b.MQTTContainer().CheckSerialNumberInTopic(message.Topic(), 5) {
		return nil
	}

	parts := message.Topic().Split()
	key := parts[len(parts)-2]

	// проверяем наличие в списке, дождавшись первоначальной загрузки
	value, ok := b.clientWiFi.LoadWait(key)

	topic := b.config.TopicWiFiMACState.Format(sn, key)

	// если в списке нет, значит удаляем из mqtt (если там не удалено)
	if !ok {
		if message.IsTrue() {
			return b.MQTTContainer().PublishAsyncRaw(ctx, topic, 1, true, "")
		}

		return nil
	}

	// если state отличается от того что в списке, значит отправляем тот что из списка, так как там мастер данные
	state := value.(bool)
	if state != message.Bool() {
		return b.MQTTContainer().PublishAsync(ctx, topic, state)
	}

	return nil
}
