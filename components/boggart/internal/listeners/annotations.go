package listeners

import (
	"context"
	"fmt"
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/devices"
	"github.com/kihamo/go-workers"
	"github.com/kihamo/go-workers/listener"
	"github.com/kihamo/shadow/components/annotations"
)

type AnnotationsListener struct {
	listener.BaseListener

	annotations annotations.Component
	startDate   *time.Time
}

func NewAnnotationsListener(annotations annotations.Component, startDate *time.Time) *AnnotationsListener {
	t := &AnnotationsListener{
		annotations: annotations,
		startDate:   startDate,
	}
	t.Init()

	return t
}

func (l *AnnotationsListener) Events() []workers.Event {
	return []workers.Event{
		devices.EventDoorGPIOReedSwitchClose,
		boggart.DeviceEventDeviceDisabledAfterCheck,
		boggart.DeviceEventDevicesManagerReady,
	}
}

func (l *AnnotationsListener) Run(_ context.Context, event workers.Event, t time.Time, args ...interface{}) {
	switch event {
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

	case boggart.DeviceEventDeviceDisabledAfterCheck:
		device := args[0].(boggart.Device)

		tags := make([]string, 0, len(device.Types()))
		for _, deviceType := range device.Types() {
			tags = append(tags, deviceType.String())
		}
		tags = append(tags, "device disabled")

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
