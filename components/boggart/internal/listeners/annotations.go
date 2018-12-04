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

	annotations     annotations.Component
	startDate       *time.Time
	applicationName string
	devicesManager  boggart.DevicesManager
}

func NewAnnotationsListener(annotations annotations.Component, applicationName string, startDate *time.Time, devicesManager boggart.DevicesManager) *AnnotationsListener {
	t := &AnnotationsListener{
		annotations:     annotations,
		applicationName: applicationName,
		startDate:       startDate,
		devicesManager:  devicesManager,
	}
	t.Init()

	return t
}

func (l *AnnotationsListener) Events() []workers.Event {
	return []workers.Event{
		boggart.DeviceEventDeviceDisabledAfterCheck,
		boggart.DeviceEventDeviceDisabled,
		boggart.DeviceEventDevicesManagerReady,
	}
}

func (l *AnnotationsListener) Run(_ context.Context, event workers.Event, t time.Time, args ...interface{}) {
	switch event {
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
		tags = append(tags, l.applicationName)

		l.annotations.CreateInStorages(
			annotations.NewAnnotation("Device is disabled", device.Description(), tags, &t, nil),
			[]string{annotations.StorageGrafana})

	case boggart.DeviceEventDevicesManagerReady:
		l.annotations.CreateInStorages(
			annotations.NewAnnotation("System is ready", "", []string{"system", l.applicationName}, &t, nil),
			[]string{annotations.StorageGrafana})
	}
}

func (l *AnnotationsListener) Name() string {
	return boggart.ComponentName + ".annotations"
}
