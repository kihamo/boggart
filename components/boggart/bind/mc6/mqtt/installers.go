package mqtt

import (
	"context"

	"github.com/kihamo/boggart/components/boggart/installer"
)

func (b *Bind) InstallersSupport() []installer.System {
	return []installer.System{
		installer.SystemDevice,
	}
}

func (b *Bind) InstallerSteps(_ context.Context, system installer.System) ([]installer.Step, error) {
	cfg := b.config()

	if system == installer.SystemDevice {
		return []installer.Step{{
			Description: "Чтобы подключить собственный MQTT сервер необходимо",
			Content: `1. Через команду 3036 выбрать режим Wifi Cloud
2. Через команду 1976 установить адрес своего mqtt сервера в виде AES256 ключа ` + cfg.AES256Key + `
3. По умолчанию устройство использует логин/пароль admin/admin для доступа к mqtt серверу, его сменить можно только через прошивку
4. Все сообщения в mqtt закодированы и приходят от устройства в топик {mac} и подписка самим устройством слушается в топике temprx/{mac}`,
		}, {
			Description: "Как задать корректно свой ключ шифрования",
			Content: `1. Он должен быть строго 32 байта
2. Иметь постфикс aes256ecb7, все что до этого постфикса - считается доменом. В ключе ` + cfg.AES256Key + ` доменом будет считаться ` + cfg.AES256Key[0:22] + `
3. Если домен не будет резолвится в локальный сервер mqtt девайс автоматически перейдет на китайское облако
4. Ключ шифрования для данных, которые отправляются в китайское облако, зашит в прошивке и равен www.brande.commc6mqtt_aes256ecb7`,
		}}, nil
	}

	return nil, nil
}
