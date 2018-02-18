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
	to             string
}

func NewTelegramListener(messenger *telegram.Telegram, devicesManager boggart.DevicesManager) *TelegramListener {
	t := &TelegramListener{
		devicesManager: devicesManager,
		messenger:      messenger,
		// FIXME:
		to: "238815343",
	}
	t.BaseListener.Init()

	return t
}

func (l *TelegramListener) Events() []workers.Event {
	return []workers.Event{
		boggart.DeviceEventDeviceDisabledAfterCheck,
		boggart.DeviceEventDeviceEnabledAfterCheck,
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

		message := fmt.Sprintf("Device %s #%s (%s) is DOWN", args[1], device.Id(), device.Description())
		if err == nil {
			l.messenger.SendMessage(l.to, message)
		} else {
			l.messenger.SendMessage(l.to, message+". Reason: "+err.(error).Error())
		}

	case boggart.DeviceEventDeviceEnabledAfterCheck:
		device := args[0].(boggart.Device)
		l.messenger.SendMessage(
			l.to,
			fmt.Sprintf("Device %s #%s (%s) is UP", args[1], device.Id(), device.Description()))
	}
}

func (l *TelegramListener) eventDoor(open bool, device boggart.Device, changed *time.Time) {
	if open {
		l.messenger.SendMessage(l.to, device.Description()+" is opened")
	} else {
		l.messenger.SendMessage(l.to, device.Description()+" is closed")
	}

	deviceCamera := l.devicesManager.Device(boggart.DeviceIdCameraHall.String())
	if deviceCamera != nil && deviceCamera.IsEnabled() {
		time.AfterFunc(time.Second, func() {
			func(camera boggart.Camera) {
				if image, err := camera.Snapshot(context.Background()); err == nil {
					l.messenger.SendPhoto(l.to, "Hall snapshot", bytes.NewReader(image))
				}
			}(deviceCamera.(boggart.Camera))
		})
	}
}

func (l *TelegramListener) Name() string {
	return boggart.ComponentName + ".telegram"
}
