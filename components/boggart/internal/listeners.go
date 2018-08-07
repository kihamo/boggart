package internal

import (
	"strings"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/internal/listeners"
	"github.com/kihamo/boggart/components/mqtt"
	"github.com/kihamo/shadow/components/annotations"
	"github.com/kihamo/shadow/components/messengers"
	"github.com/kihamo/shadow/components/messengers/platforms/telegram"
)

func (c *Component) initListeners() {
	c.listenersManager.AddListener(listeners.NewLoggingListener(c.logger))

	c.listenersManager.AddListener(listeners.NewMQTTListener(c.application.GetComponent(mqtt.ComponentName).(mqtt.Component)))

	if c.application.HasComponent(messengers.ComponentName) {
		messenger := c.application.GetComponent(messengers.ComponentName).(messengers.Component).
			Messenger(messengers.MessengerTelegram)

		if messenger != nil {
			var chats []string
			chatsFromConfig := strings.FieldsFunc(c.config.String(boggart.ConfigListenerTelegramChats), func(c rune) bool {
				return c == ','
			})

			for _, id := range chatsFromConfig {
				chats = append(chats, id)
			}

			if len(chats) > 0 {
				c.listenersManager.AddListener(listeners.NewTelegramListener(
					messenger.(*telegram.Telegram),
					c.devicesManager,
					c.application,
					chats))
			}
		}
	}

	if c.application.HasComponent(annotations.ComponentName) {
		c.listenersManager.AddListener(listeners.NewAnnotationsListener(
			c.application.GetComponent(annotations.ComponentName).(annotations.Component),
			c.application.StartDate(),
			c.devicesManager))
	}
}
