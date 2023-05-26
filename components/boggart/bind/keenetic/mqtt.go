package keenetic

import (
	"context"
	"errors"
	"regexp"
	"time"

	"github.com/kihamo/boggart/components/boggart/tasks"
	"github.com/kihamo/boggart/components/mqtt"
)

var (
	// Отключение:  Network::Interface::Rtx::WifiMonitor: "WifiMaster1/AccessPoint0": STA(XX:XX:XX:XX:XX:XX) had deauthenticated by STA (reason: STA is leaving or has left BSS).
	// Подключение: Network::Interface::Rtx::WifiMonitor: "WifiMaster1/AccessPoint0": STA(XX:XX:XX:XX:XX:XX) had associated successfully (FT mode).

	wifiMonitorRegexp = regexp.MustCompile(`STA\(([^\)].+?)\) had (associated|deauthenticated|disassociated)`)
)

func (b *Bind) callbackMQTTHotspotSearchZombies(ctx context.Context, _ mqtt.Component, message mqtt.Message) error {
	if message.IsFalse() || !message.Retained() {
		return nil
	}

	topic := message.Topic()
	parts := topic.Split()
	if len(parts) < 2 {
		return nil
	}

	if _, ok := b.hotspotConnections.Load(parts[len(parts)-2]); !ok {
		return b.MQTT().PublishAsyncRawWithoutCache(ctx, topic, message.Qos(), true, nil)
	}

	return nil
}

// Для ускорения проверки изменения в подключениях, чекаем syslog сообщения и парсим оттуда необходиму инфу
func (b *Bind) callbackMQTTSyslog(ctx context.Context, _ mqtt.Component, message mqtt.Message) error {
	const tagNDM = "ndm"

	var request map[string]interface{}

	if err := message.JSONUnmarshal(&request); err != nil {
		return err
	}

	tag, ok := request["tag"]
	if !ok || tag != tagNDM {
		return nil
	}

	content, ok := request["content"]
	if !ok {
		return nil
	}

	if !wifiMonitorRegexp.MatchString(content.(string)) {
		return nil
	}

	// принудительно запускаем синхронизацию, чтобы два раза логику не писать
	// у роутера есть задержка, сразу после лога подключение не видно через api
	// в новом статусе, поэтому искусственно делаем задержку (на секунде успевал
	// кеш обновиться, но лучше вынести в параметр)
	time.AfterFunc(b.config().SyncAfterSyslogDelay, func() {
		err := b.Workers().TaskRunByName(ctx, TaskNameHotspotSync)
		if err != nil && !errors.Is(err, tasks.ErrAlreadyRunning) {
			b.Logger().Error("Run task hotspot sync error", "error", err)
		}
	})

	return nil
}
