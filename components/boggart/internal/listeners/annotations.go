package listeners

import (
	"context"
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/go-workers"
	"github.com/kihamo/go-workers/listener"
	"github.com/kihamo/shadow/components/annotations"
)

type AnnotationsListener struct {
	listener.BaseListener

	annotations    annotations.Component
	startDate      *time.Time
	devicesManager boggart.DevicesManager
}

func NewAnnotationsListener(annotations annotations.Component, startDate *time.Time, devicesManager boggart.DevicesManager) *AnnotationsListener {
	t := &AnnotationsListener{
		annotations:    annotations,
		startDate:      startDate,
		devicesManager: devicesManager,
	}
	t.Init()

	return t
}

func (l *AnnotationsListener) Events() []workers.Event {
	return []workers.Event{
		//devices.EventDoorGPIOReedSwitchClose,
		// devices.EventUPSApcupsdStatusChanged,
		boggart.DeviceEventDeviceDisabledAfterCheck,
		boggart.DeviceEventDeviceDisabled,
		boggart.DeviceEventDevicesManagerReady,
	}
}

func (l *AnnotationsListener) Run(_ context.Context, event workers.Event, t time.Time, args ...interface{}) {
	switch event {
	/*
		case devices.EventDoorGPIOReedSwitchClose:
			var changed *time.Time

			if args[1] == nil {
				changed = l.startDate
			} else {
				changed = args[1].(*time.Time)
			}

			timeEnd := time.Now()
			diff := timeEnd.Sub(*changed)

			device := args[0].(boggart.Device)

			tags := make([]string, 0, len(device.Types()))
			for _, deviceType := range device.Types() {
				tags = append(tags, deviceType.String())
			}
			tags = append(tags, "door closed")

			l.annotations.Create(annotations.NewAnnotation(
				"Door is closed",
				fmt.Sprintf("Door was open for %.2f seconds", diff.Seconds()),
				tags,
				changed,
				&timeEnd))
			/*
				case devices.EventUPSApcupsdStatusChanged:
					if args[1] == nil || args[2] == nil {
						break
					}

					statusCurrent := args[1].(apcupsd.Status)
					if statusCurrent.Status == nil || statusCurrent.XOnBattery == nil || statusCurrent.XOffBattery == nil {
						break
					}

					statusPrev := args[2].(apcupsd.Status)
					if statusPrev.Status == nil {
						break
					}

					if (*statusCurrent.Status).IsOnline && !(*statusPrev.Status).IsOnline {
						diff := statusCurrent.XOffBattery.Sub(*statusCurrent.XOnBattery)

						var reason string
						if statusCurrent.LastTransfer != nil {
							reason = *statusCurrent.LastTransfer
						}

						l.annotations.CreateInStorages(
							annotations.NewAnnotation(
								"UPS switched",
								fmt.Sprintf("UPS switched to power. Offline for %.2f seconds with reason %s", diff.Seconds(), reason),
								[]string{"UPS"},
								statusCurrent.XOnBattery,
								statusCurrent.XOffBattery),
							[]string{annotations.StorageGrafana})
					}
	*/
	case boggart.DeviceEventDeviceDisabledAfterCheck, boggart.DeviceEventDeviceDisabled:
		if !l.devicesManager.IsReady() {
			return
		}

		device := args[0].(boggart.Device)

		tags := make([]string, 0, len(device.Types()))
		for _, deviceType := range device.Types() {
			tags = append(tags, deviceType.String())
		}
		tags = append(tags, "device disabled")

		if event == boggart.DeviceEventDeviceDisabled {
			tags = append(tags, "manually")
		}

		l.annotations.CreateInStorages(
			annotations.NewAnnotation("Device is disabled", device.Description(), tags, &t, nil),
			[]string{annotations.StorageGrafana})

	case boggart.DeviceEventDevicesManagerReady:
		l.annotations.CreateInStorages(
			annotations.NewAnnotation("System is ready", "", []string{"system"}, &t, nil),
			[]string{annotations.StorageGrafana})
	}
}

func (l *AnnotationsListener) Name() string {
	return boggart.ComponentName + ".annotations"
}
