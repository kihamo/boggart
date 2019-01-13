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
	case boggart.BindEventDeviceRegister:
		l.logger.Debug("Register device",
			"device.id", args[0].(boggart.Device).ID(),
			"device_manager.id", args[1],
		)

	case boggart.BindEventDeviceEnabled:
		l.logger.Info("Device has been enabled manually",
			"device.id", args[0].(boggart.Device).ID(),
		)

	case boggart.BindEventDeviceDisabled:
		l.logger.Info("Device has been disabled manually",
			"device.id", args[0].(boggart.Device).ID(),
		)

	case boggart.BindEventDeviceDisabledAfterCheck:
		err := args[2]

		if err == nil {
			l.logger.Warn("Device has been disabled because the ping returns false",
				"device.id", args[0].(boggart.Device).ID(),
				"device.key", args[1],
			)
		} else {
			l.logger.Warn("Device has been disabled because the ping failed",
				"device.id", args[0].(boggart.Device).ID(),
				"device.key", args[1],
				"error", err.(error).Error(),
			)
		}

	case boggart.BindEventDeviceEnabledAfterCheck:
		l.logger.Info("Device has been enabled because the ping returns true",
			"device.id", args[0].(boggart.Device).ID(),
			"device.key", args[1],
		)

	case boggart.BindEventSyslogReceive:
		fields := make([]interface{}, 0)
		record := args[0].(map[string]interface{})

		for k, v := range record {
			fields = append(fields, k, v)
		}

		l.logger.Debug("Syslog message receive", fields...)

	case boggart.BindEventDevicesManagerReady:
		l.logger.Debug("Device manager is ready")

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
