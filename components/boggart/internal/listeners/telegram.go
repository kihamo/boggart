package listeners

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/devices"
	"github.com/kihamo/go-workers"
	"github.com/kihamo/go-workers/listener"
	"github.com/kihamo/shadow/components/messengers/platforms/telegram"
)

type TelegramListener struct {
	listener.BaseListener

	devicesManager boggart.DevicesManager
	messenger      *telegram.Telegram
	chats          []string
}

func NewTelegramListener(messenger *telegram.Telegram, devicesManager boggart.DevicesManager, chats []string) *TelegramListener {
	t := &TelegramListener{
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
		boggart.DeviceEventWifiClientConnected,
		boggart.DeviceEventWifiClientDisconnected,
		devices.EventDoorGPIOReedSwitchOpen,
		devices.EventDoorGPIOReedSwitchClose,
	}
}

func (l *TelegramListener) Run(_ context.Context, event workers.Event, t time.Time, args ...interface{}) {
	switch event {
	case devices.EventDoorGPIOReedSwitchOpen:
		l.eventDoor(true, args[0].(boggart.Device), args[1].(*time.Time))

	case devices.EventDoorGPIOReedSwitchClose:
		l.eventDoor(false, args[0].(boggart.Device), args[1].(*time.Time))

	case boggart.DeviceEventDeviceDisabledAfterCheck:
		device := args[0].(boggart.Device)
		err := args[2]

		message := fmt.Sprintf("Device is down %s #%s (%s)", args[1], device.Id(), device.Description())
		if err == nil {
			l.send(message)
		} else {
			l.send(message + ". Reason: " + err.(error).Error())
		}

	case boggart.DeviceEventDeviceEnabledAfterCheck:
		device := args[0].(boggart.Device)
		l.send(fmt.Sprintf("Device is up %s #%s (%s)", args[1], device.Id(), device.Description()))

	case boggart.DeviceEventDevicesManagerReady:
		l.send("Hello. I'm online and ready")

	case boggart.DeviceEventWifiClientConnected:
		l.send(fmt.Sprintf("%s connected to %s", args[0], args[1]))

	case boggart.DeviceEventWifiClientDisconnected:
		l.send(fmt.Sprintf("%s disconnected to %s", args[0], args[1]))
	}
}

func (l *TelegramListener) eventDoor(open bool, device boggart.Device, changed *time.Time) {
	if open {
		l.send(device.Description() + " is opened")
	} else {
		l.send(device.Description() + " is closed")
	}

	deviceCamera := l.devicesManager.Device(boggart.DeviceIdCameraHall.String())
	if deviceCamera != nil && deviceCamera.IsEnabled() {
		time.AfterFunc(time.Second, func() {
			func(camera boggart.Camera) {
				if image, err := camera.Snapshot(context.Background()); err == nil {
					for _, chatId := range l.chats {
						l.messenger.SendPhoto(chatId, "Hall snapshot", bytes.NewReader(image))
					}
				}
			}(deviceCamera.(boggart.Camera))
		})
	}
}

func (l *TelegramListener) send(message string) {
	for _, chatId := range l.chats {
		l.messenger.SendMessage(chatId, message)
	}
}

func (l *TelegramListener) Name() string {
	return boggart.ComponentName + ".telegram"
}
