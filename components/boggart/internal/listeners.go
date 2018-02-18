package internal

import (
	"strings"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/internal/listeners"
	w "github.com/kihamo/go-workers"
	"github.com/kihamo/shadow/components/annotations"
	"github.com/kihamo/shadow/components/messengers"
	"github.com/kihamo/shadow/components/messengers/platforms/telegram"
)

func (c *Component) initListeners() {
	c.devicesManager.Attach(w.EventAll, listeners.NewLoggingListener(c.logger))

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
				listenerTelegram := listeners.NewTelegramListener(messenger.(*telegram.Telegram), c.devicesManager, chats)

				for _, event := range listenerTelegram.Events() {
					c.devicesManager.Attach(event, listenerTelegram)
				}
			}
		}
	}

	if c.application.HasComponent(annotations.ComponentName) {
		listenerAnnotations := listeners.NewAnnotationsListener(
			c.application.GetComponent(annotations.ComponentName).(annotations.Component),
			c.application.StartDate())

		for _, event := range listenerAnnotations.Events() {
			c.devicesManager.Attach(event, listenerAnnotations)
		}
	}
}
