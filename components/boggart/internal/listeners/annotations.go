package listeners

import (
	"context"
	"time"

	"github.com/kihamo/boggart/components/boggart/internal/manager"

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
	manager         *manager.Manager
}

func NewAnnotationsListener(annotations annotations.Component, applicationName string, startDate *time.Time, manager *manager.Manager) *AnnotationsListener {
	t := &AnnotationsListener{
		annotations:     annotations,
		applicationName: applicationName,
		startDate:       startDate,
		manager:         manager,
	}
	t.Init()

	return t
}

func (l *AnnotationsListener) Events() []workers.Event {
	return []workers.Event{
		boggart.BindEventDeviceDisabledAfterCheck,
		boggart.BindEventDeviceDisabled,
		boggart.BindEventDevicesManagerReady,
	}
}

func (l *AnnotationsListener) Run(_ context.Context, event workers.Event, t time.Time, args ...interface{}) {
	switch event {
	case boggart.BindEventDeviceDisabledAfterCheck, boggart.BindEventDeviceDisabled:
		if !l.manager.IsReady() {
			return
		}

		device := args[0].(boggart.Device)
		tags := append(device.Tags(), "device disabled")

		if event == boggart.BindEventDeviceDisabled {
			tags = append(tags, "manually")
		}
		tags = append(tags, l.applicationName)

		l.annotations.CreateInStorages(
			annotations.NewAnnotation("Device is disabled", device.Description(), tags, &t, nil),
			[]string{annotations.StorageGrafana})

	case boggart.BindEventDevicesManagerReady:
		l.annotations.CreateInStorages(
			annotations.NewAnnotation("System is ready", "", []string{"system", l.applicationName}, &t, nil),
			[]string{annotations.StorageGrafana})
	}
}

func (l *AnnotationsListener) Name() string {
	return boggart.ComponentName + ".annotations"
}
