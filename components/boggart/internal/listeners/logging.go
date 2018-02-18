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
	t.BaseListener.Init()

	return t
}

func (l *LoggingListener) Run(_ context.Context, event workers.Event, t time.Time, args ...interface{}) {
	switch event {
	case boggart.DeviceEventDeviceRegister:
		l.logger.Debug("Register device", map[string]interface{}{
			"device.id":         args[0].(boggart.Device).Id(),
			"device_manager.id": args[1],
		})

	default:
		l.logger.Info("Fire unknown event", map[string]interface{}{
			"event.id":   event.Id(),
			"event.name": event.Name(),
			"args":       args,
		})
	}
}

func (l *LoggingListener) Name() string {
	return boggart.ComponentName + ".logging"
}
