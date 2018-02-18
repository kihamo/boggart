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
	t.BaseListener.Init()

	return t
}

func (l *AnnotationsListener) Events() []workers.Event {
	return []workers.Event{
		devices.EventDoorGPIOReedSwitchClose,
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

		annotation := annotations.NewAnnotation(
			"Door is closed",
			fmt.Sprintf("Door was open for %.2f seconds", diff.Seconds()),
			[]string{"door", "close"},
			changed,
			&timeEnd)

		// TODO: err log
		l.annotations.Create(annotation)
	}
}

func (l *AnnotationsListener) Name() string {
	return boggart.ComponentName + ".annotations"
}
