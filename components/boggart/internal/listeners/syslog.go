package listeners

import (
	"context"
	"fmt"
	"net"
	"regexp"
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/go-workers"
	"github.com/kihamo/go-workers/listener"
	"github.com/kihamo/go-workers/manager"
)

const (
	SyslogFieldTag     = "tag"
	SyslogFieldContent = "content"
)

var (
	wifiClientRegexp = regexp.MustCompile(`^([^@]+)@([^:\s]+):\s+([^\s,]+)`)
)

type SyslogListener struct {
	listener.BaseListener

	listenersManager *manager.ListenersManager
}

func NewSyslogListener(listenersManager *manager.ListenersManager) *SyslogListener {
	l := &SyslogListener{
		listenersManager: listenersManager,
	}
	l.Init()

	return l
}

func (l *SyslogListener) Events() []workers.Event {
	return []workers.Event{
		boggart.DeviceEventSyslogReceive,
	}
}

func (l *SyslogListener) Run(_ context.Context, event workers.Event, t time.Time, args ...interface{}) {
	switch event {
	case boggart.DeviceEventSyslogReceive:
		message := args[0].(map[string]interface{})

		tag, ok := message[SyslogFieldTag]
		if !ok {
			return
		}

		content, ok := message[SyslogFieldContent]
		if !ok {
			return
		}

		check := wifiClientRegexp.FindStringSubmatch(fmt.Sprintf("%s:%s", tag, content))
		if len(check) >= 4 {
			if _, err := net.ParseMAC(check[1]); err == nil {
				switch check[3] {
				case "connected":
					l.listenersManager.AsyncTrigger(boggart.DeviceEventWifiClientConnected, check[1], check[2])
				case "disconnected":
					l.listenersManager.AsyncTrigger(boggart.DeviceEventWifiClientDisconnected, check[1], check[2])
				}
			}
		}
	}
}

func (l *SyslogListener) Name() string {
	return boggart.ComponentName + ".syslog"
}
