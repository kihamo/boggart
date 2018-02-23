package listeners

import (
	"context"
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/go-workers"
	"github.com/kihamo/go-workers/listener"
	"github.com/kihamo/shadow/components/logger"
)

type LoggingListener struct {
	listener.BaseListener

	logger logger.Logger
}

func NewLoggingListener(logger logger.Logger) *LoggingListener {
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
	case boggart.DeviceEventDeviceRegister:
		l.logger.Debug("Register device", map[string]interface{}{
			"device.id":         args[0].(boggart.Device).Id(),
			"device_manager.id": args[1],
		})

	case boggart.DeviceEventDeviceDisabledAfterCheck:
		err := args[2]

		if err == nil {
			l.logger.Warn("Device has been disabled because the ping returns false", map[string]interface{}{
				"device.id":  args[0].(boggart.Device).Id(),
				"device.key": args[1],
			})
		} else {
			l.logger.Warn("Device has been disabled because the ping failed", map[string]interface{}{
				"device.id":  args[0].(boggart.Device).Id(),
				"device.key": args[1],
				"error":      err.(error).Error(),
			})
		}

	case boggart.DeviceEventDeviceEnabledAfterCheck:
		l.logger.Info("Device has been enabled because the ping returns true", map[string]interface{}{
			"device.id":  args[0].(boggart.Device).Id(),
			"device.key": args[1],
		})

	case boggart.DeviceEventSyslogReceive:
		l.logger.Debug("Syslog message receive", args[0].(map[string]interface{}))

	default:
		l.logger.Debug("Fire unknown event", map[string]interface{}{
			"event.id":   event.Id(),
			"event.name": event.Name(),
			"args":       args,
		})
	}
}

func (l *LoggingListener) Name() string {
	return boggart.ComponentName + ".logging"
}
