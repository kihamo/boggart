package internal

import (
	"github.com/kihamo/boggart/components/boggart"
)

func (c *Component) SyslogHandler(message map[string]interface{}) {
	c.listenersManager.AsyncTrigger(boggart.DeviceEventSyslogReceive, message)
}
