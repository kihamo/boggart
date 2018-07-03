package listeners

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/devices"
	"github.com/kihamo/boggart/components/boggart/protocols/apcupsd"
	"github.com/kihamo/boggart/components/boggart/providers/hikvision"
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
		boggart.DeviceEventHikvisionEventNotificationAlert,
		boggart.DeviceEventDeviceDisabledAfterCheck,
		boggart.DeviceEventDeviceEnabledAfterCheck,
		boggart.DeviceEventDevicesManagerReady,
		boggart.DeviceEventWifiClientConnected,
		boggart.DeviceEventWifiClientDisconnected,
		devices.EventDoorGPIOReedSwitchOpen,
		devices.EventDoorGPIOReedSwitchClose,
		devices.EventUPSApcupsdStatusChanged,
	}
}

func (l *TelegramListener) Run(_ context.Context, event workers.Event, t time.Time, args ...interface{}) {
	switch event {
	case devices.EventDoorGPIOReedSwitchOpen:
		l.sendMessage(args[0].(boggart.Device).Description() + " is opened")

		if camera := l.devicesManager.Device(boggart.DeviceIdCameraHall.String()); camera != nil {
			time.AfterFunc(time.Second, func() {
				func(camera boggart.Camera) {
					l.sendSnapshotCamera(camera.(boggart.Camera))
				}(camera.(boggart.Camera))
			})
		}

	case devices.EventDoorGPIOReedSwitchClose:
		l.sendMessage(args[0].(boggart.Device).Description() + " is closed")

		if camera := l.devicesManager.Device(boggart.DeviceIdCameraHall.String()); camera != nil {
			time.AfterFunc(time.Second, func() {
				func(camera boggart.Camera) {
					l.sendSnapshotCamera(camera.(boggart.Camera))
				}(camera.(boggart.Camera))
			})
		}

	case devices.EventUPSApcupsdStatusChanged:
		current := args[1]
		prev := args[2]

		if current != nil && prev != nil {
			status := current.(apcupsd.Status)
			var (
				message string
				reason  string
			)

			if status.Status != nil {
				if (*status.Status).IsOnBattery {
					message = "UPS switch to battery"

					if status.LastTransfer != nil {
						reason = *status.LastTransfer
					}
				} else if (*status.Status).IsOnline {
					message = "UPS switch to power"

					if status.LastTransfer != nil {
						reason = *status.LastTransfer
					}

					if status.XOnBattery != nil || status.XOffBattery != nil {
						diff := status.XOffBattery.Sub(*status.XOnBattery)

						message += fmt.Sprintf(". Offline for %.2f seconds", diff.Seconds())
					}
				}
			}

			if reason != "" {
				l.sendMessage(fmt.Sprintf("%s. Reason: %s", message, reason))
			} else {
				l.sendMessage(message)
			}
		}

	case boggart.DeviceEventDeviceDisabledAfterCheck:
		device := args[0].(boggart.Device)
		err := args[2]

		message := fmt.Sprintf("Device is down %s #%s (%s)", args[1], device.Id(), device.Description())
		if err == nil {
			l.sendMessage(message)
		} else {
			l.sendMessage(message + ". Reason: " + err.(error).Error())
		}

	case boggart.DeviceEventDeviceEnabledAfterCheck:
		device := args[0].(boggart.Device)
		l.sendMessage(fmt.Sprintf("Device is up %s #%s (%s)", args[1], device.Id(), device.Description()))

	case boggart.DeviceEventDevicesManagerReady:
		l.sendMessage("Hello. I'm online and ready")

	case boggart.DeviceEventWifiClientConnected:
		mac := args[1].(*devices.MikrotikRouterMac)

		l.sendMessage(fmt.Sprintf("%s with IP %s (%s, %s) connected to %s", mac.Address, mac.ARP.IP, mac.ARP.Comment, mac.DHCP.Hostname, args[2]))

	case boggart.DeviceEventWifiClientDisconnected:
		mac := args[1].(*devices.MikrotikRouterMac)

		l.sendMessage(fmt.Sprintf("%s with IP %s (%s, %s) disconnected to %s", mac.Address, mac.ARP.IP, mac.ARP.Comment, mac.DHCP.Hostname, args[2]))

	case boggart.DeviceEventHikvisionEventNotificationAlert:
		event := args[1].(*hikvision.EventNotificationAlertStreamResponse)

		// FIXME: ID прибит гвоздями, нужен реальный
		videoRecorderDevice := l.devicesManager.Device(boggart.DeviceIdVideoRecorder.String()).(boggart.VideoRecorder)

		switch event.EventType {
		case hikvision.EventTypeIO:
		case hikvision.EventTypeVMD:
			go func() {
				l.sendSnapshotFromVideoRecorder(videoRecorderDevice, event)
			}()

			l.sendMessage(fmt.Sprintf("Hikvision alert %s %s", event.EventType, event.EventDescription))

		case hikvision.EventTypeVideoLoss:
		case hikvision.EventTypeShelterAlarm:
		case hikvision.EventTypeFaceDetection:
		case hikvision.EventTypeDefocus:
		case hikvision.EventTypeAudioException:
		case hikvision.EventTypeSceneChangeDetection:
		case hikvision.EventTypeFieldDetection:
		case hikvision.EventTypeLineDetection:
		case hikvision.EventTypeRegionEntrance:
		case hikvision.EventTypeRegionExiting:
		case hikvision.EventTypeLoitering:
		case hikvision.EventTypeGroup:
		case hikvision.EventTypeRapidMove:
		case hikvision.EventTypeParking:
		case hikvision.EventTypeUnattendedBaggage:
		case hikvision.EventTypeAttendedBaggage:
		case hikvision.EventTypePIR:
		case hikvision.EventTypePeopleDetection:
		default:
			go func() {
				l.sendSnapshotFromVideoRecorder(videoRecorderDevice, event)
			}()

			l.sendMessage(fmt.Sprintf("Hikvision alert with unknown type %s", event.EventType))
		}
	}
}

func (l *TelegramListener) sendSnapshotFromVideoRecorder(videoRecorder boggart.VideoRecorder, event *hikvision.EventNotificationAlertStreamResponse) {
	if !videoRecorder.IsEnabled() {
		return
	}

	if image, err := videoRecorder.Snapshot(context.Background(), event.DynChannelID, 1); err == nil {
		for _, chatId := range l.chats {
			l.messenger.SendPhoto(chatId, videoRecorder.Description()+" snapshot", bytes.NewReader(image))
		}
	} else {
		fmt.Println("sendSnapshotFromVideoRecorder", err)
	}
}

func (l *TelegramListener) sendSnapshotCamera(camera boggart.Camera) {
	if !camera.IsEnabled() {
		return
	}

	if image, err := camera.Snapshot(context.Background()); err == nil {
		for _, chatId := range l.chats {
			l.messenger.SendPhoto(chatId, camera.Description()+" snapshot", bytes.NewReader(image))
		}
	}
}

func (l *TelegramListener) sendMessage(message string) {
	for _, chatId := range l.chats {
		l.messenger.SendMessage(chatId, message)
	}
}

func (l *TelegramListener) Name() string {
	return boggart.ComponentName + ".telegram"
}
