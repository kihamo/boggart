package listeners

import (
	"context"
	"fmt"
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/go-workers"
	"github.com/kihamo/go-workers/listener"
	"github.com/kihamo/shadow"
	"github.com/kihamo/shadow/components/messengers/platforms/telegram"
)

type TelegramListener struct {
	listener.BaseListener

	application    shadow.Application
	devicesManager boggart.DevicesManager
	messenger      *telegram.Telegram
	chats          []string
}

func NewTelegramListener(messenger *telegram.Telegram, devicesManager boggart.DevicesManager, application shadow.Application, chats []string) *TelegramListener {
	t := &TelegramListener{
		application:    application,
		devicesManager: devicesManager,
		messenger:      messenger,
		chats:          chats,
	}
	t.Init()

	return t
}

func (l *TelegramListener) Events() []workers.Event {
	return []workers.Event{
		boggart.DeviceEventDeviceDisabledAfterCheck,
		boggart.DeviceEventDeviceEnabledAfterCheck,
		boggart.DeviceEventDevicesManagerReady,
	}
}

func (l *TelegramListener) Run(_ context.Context, event workers.Event, t time.Time, args ...interface{}) {
	switch event {
	case boggart.DeviceEventDeviceDisabledAfterCheck:
		if !l.devicesManager.IsReady() {
			return
		}

		device := args[0].(boggart.Device)
		err := args[2]

		message := fmt.Sprintf("Device is down %s #%s (%s)", args[1], device.ID(), device.Description())
		if err == nil {
			l.sendMessage(message)
		} else {
			l.sendMessage(message + ". Reason: " + err.(error).Error())
		}

	case boggart.DeviceEventDeviceEnabledAfterCheck:
		if !l.devicesManager.IsReady() {
			return
		}

		device := args[0].(boggart.Device)
		l.sendMessage("Device is up %s #%s (%s)", args[1], device.ID(), device.Description())

	case boggart.DeviceEventDevicesManagerReady:
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
