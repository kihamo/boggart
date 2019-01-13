package listeners

import (
	"context"
	"fmt"
	"time"

	"github.com/kihamo/boggart/components/boggart/internal/manager"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/go-workers"
	"github.com/kihamo/go-workers/listener"
	"github.com/kihamo/shadow"
	"github.com/kihamo/shadow/components/messengers/platforms/telegram"
)

type TelegramListener struct {
	listener.BaseListener

	application shadow.Application
	manager     *manager.Manager
	messenger   *telegram.Telegram
	chats       []string
}

func NewTelegramListener(messenger *telegram.Telegram, manager *manager.Manager, application shadow.Application, chats []string) *TelegramListener {
	t := &TelegramListener{
		application: application,
		manager:     manager,
		messenger:   messenger,
		chats:       chats,
	}
	t.Init()

	return t
}

func (l *TelegramListener) Events() []workers.Event {
	return []workers.Event{
		boggart.BindEventManagerReady,
	}
}

func (l *TelegramListener) Run(_ context.Context, event workers.Event, t time.Time, args ...interface{}) {
	switch event {

	case boggart.BindEventManagerReady:
		l.sendMessage("Hello. I'm %s and I'm online and ready", l.application.Name())
	}
}

func (l *TelegramListener) sendMessage(message string, a ...interface{}) {
	for _, chatId := range l.chats {
		l.messenger.SendMessage(chatId, fmt.Sprintf(message, a...))
	}
}

func (l *TelegramListener) Name() string {
	return boggart.ComponentName + ".telegram"
}
