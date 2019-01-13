package listeners

import (
	"context"
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/go-workers"
	"github.com/kihamo/go-workers/listener"
	"github.com/kihamo/shadow/components/logging"
)

type LoggingListener struct {
	listener.BaseListener

	logger logging.Logger
}

func NewLoggingListener(logger logging.Logger) *LoggingListener {
	t := &LoggingListener{
		logger: logger,
	}
	t.Init()

	return t
}

func (l *LoggingListener) Events() []workers.Event {
	return []workers.Event{
		workers.EventAll,
	}
}

func (l *LoggingListener) Run(_ context.Context, event workers.Event, t time.Time, args ...interface{}) {
	switch event {
	case boggart.BindEventSyslogReceive:
		fields := make([]interface{}, 0)
		record := args[0].(map[string]interface{})

		for k, v := range record {
			fields = append(fields, k, v)
		}

		l.logger.Debug("Syslog message receive", fields...)

	case boggart.BindEventManagerReady:
		l.logger.Debug("Manager is ready")

	default:
		l.logger.Warn("Fire unknown event",
			"event.id", event.Id(),
			"event.name", event.Name(),
			"args", args,
		)
	}
}

func (l *LoggingListener) Name() string {
	return boggart.ComponentName + ".logging"
}
