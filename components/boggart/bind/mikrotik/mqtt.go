package mikrotik

import (
	"context"
	"net"

	"github.com/kihamo/boggart/components/mqtt"
)

// подписываемся на собственные топики, это нужно для того чтобы при первом старте
// выключить активность соединений, которые удалены в самом микротике и соответственно
// они будут не доступны в последующий обновлениях, а в топике и них ретеншен сообщение с true
func (b *Bind) callbackMQTTInterfacesZombies(ctx context.Context, _ mqtt.Component, message mqtt.Message) error {
	if message.IsFalse() {
		return nil
	}
	/*
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
	*/
	return nil
}

// Для ускорения проверки изменения в подключениях, чекаем syslog сообщения и парсим оттуда необходиму инфу
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
	case b.config.SyslogTagWireless:
		check := wifiClientRegexp.FindStringSubmatch(content.(string))
		if len(check) < 4 {
			return nil
		}

		if _, err := net.ParseMAC(check[1]); err != nil {
			return err
		}

		if check[3] == "connected" || check[3] == "disconnected" {
			return b.Workers().TaskRunByName(ctx, TaskNameInterfaceConnection)
		}

		return b.Workers().TaskRunByName(ctx, TaskNameInterfaceConnection)

	case b.config.SyslogTagL2TP:
		check := vpnClientRegexp.FindStringSubmatch(content.(string))
		if len(check) < 2 {
			return nil
		}

		if check[2] == "in" || check[2] == "out" {
			return b.Workers().TaskRunByName(ctx, TaskNameInterfaceConnection)
		}
	}

	return nil
}
