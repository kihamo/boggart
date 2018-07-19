package internal

import (
	"net/url"
	"strings"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/internal/listeners"
	"github.com/kihamo/shadow/components/annotations"
	"github.com/kihamo/shadow/components/messengers"
	"github.com/kihamo/shadow/components/messengers/platforms/telegram"
)

func (c *Component) initListeners() {
	c.listenersManager.AddListener(listeners.NewLoggingListener(c.logger))

	if c.config.Bool(boggart.ConfigMQTTEnabled) {
		servers := make([]*url.URL, 0)
		for _, u := range strings.Split(c.config.String(boggart.ConfigMQTTServers), ";") {
			if p, err := url.Parse(u); err == nil {
				servers = append(servers, p)
			}
		}

		c.listenersManager.AddListener(listeners.NewMQTTListener(
			servers,
			c.config.String(boggart.ConfigMQTTUsername),
			c.config.String(boggart.ConfigMQTTPassword)))
	}

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
					c.securityManager,
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
