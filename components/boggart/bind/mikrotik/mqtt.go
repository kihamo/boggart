package mikrotik

import (
	"context"
	"errors"
	"net"

	"github.com/kihamo/boggart/components/mqtt"
)

func (b *Bind) MQTTSubscribers() []mqtt.Subscriber {
	subscribers := []mqtt.Subscriber{
		mqtt.NewSubscriber(b.config.TopicWiFiMACState, 0, b.callbackMQTTWiFiSync),
	}

	if b.config.TopicSyslog != "" {
		subscribers = append(subscribers, mqtt.NewSubscriber(b.config.TopicSyslog, 0, b.callbackMQTTSyslog))
	}

	return subscribers
}

func (b *Bind) callbackMQTTWiFiSync(ctx context.Context, client mqtt.Component, message mqtt.Message) error {
	sn := b.SerialNumberWait()
	if sn == "" || !b.MQTT().CheckSerialNumberInTopic(message.Topic(), 5) {
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
			return b.MQTT().PublishAsyncRaw(ctx, topic, 1, true, "")
		}

		return nil
	}

	// если state отличается от того что в списке, значит отправляем тот что из списка, так как там мастер данные
	state := value.(bool)
	if state != message.Bool() {
		return b.MQTT().PublishAsync(ctx, topic, state)
	}

	return nil
}

func (b *Bind) callbackMQTTSyslog(ctx context.Context, _ mqtt.Component, message mqtt.Message) error {
	var request map[string]interface{}

	if err := message.JSONUnmarshal(&request); err != nil {
		return err
	}

	tag, ok := request["tag"]
	if !ok {
		return nil
	}

	content, ok := request["content"]
	if !ok {
		return nil
	}

	switch tag {
	case "wifi":
		sn := b.SerialNumberWait()
		if sn == "" {
			return errors.New("serial number is empty")
		}

		check := wifiClientRegexp.FindStringSubmatch(content.(string))
		if len(check) < 4 {
			return nil
		}

		if _, err := net.ParseMAC(check[1]); err != nil {
			return err
		}

		login := mqtt.NameReplace(check[1])

		switch check[3] {
		case "connected":
			// TODO:
			_ = b.MQTT().PublishAsyncWithoutCache(ctx, b.config.TopicWiFiConnectedMAC.Format(sn), login)
			b.updateWiFiClient(ctx)
		case "disconnected":
			// TODO:
			_ = b.MQTT().PublishAsyncWithoutCache(ctx, b.config.TopicWiFiDisconnectedMAC.Format(sn), login)
			b.updateWiFiClient(ctx)
		}

	case "vpn":
		sn := b.SerialNumberWait()
		if sn == "" {
			return errors.New("serial number is empty")
		}

		check := vpnClientRegexp.FindStringSubmatch(content.(string))
		if len(check) < 2 {
			return nil
		}

		login := mqtt.NameReplace(check[1])

		switch check[2] {
		case "in":
			// TODO:
			_ = b.MQTT().PublishAsyncWithoutCache(ctx, b.config.TopicVPNConnectedLogin.Format(sn), login)
			b.updateVPNClient(ctx)
		case "out":
			// TODO:
			_ = b.MQTT().PublishAsyncWithoutCache(ctx, b.config.TopicVPNDisconnectedLogin.Format(sn), login)
			b.updateVPNClient(ctx)
		}
	}

	return nil
}
