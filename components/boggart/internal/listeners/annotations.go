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
		boggart.BindEventManagerReady,
	}
}

func (l *AnnotationsListener) Run(_ context.Context, event workers.Event, t time.Time, args ...interface{}) {
	switch event {
	case boggart.BindEventManagerReady:
		l.annotations.CreateInStorages(
			annotations.NewAnnotation("System is ready", "", []string{"system", l.applicationName}, &t, nil),
			[]string{annotations.StorageGrafana})
	}
}

func (l *AnnotationsListener) Name() string {
	return boggart.ComponentName + ".annotations"
}
