package internal

import (
	"github.com/kihamo/boggart/components/boggart"
)

func (c *Component) SyslogHandler(message map[string]interface{}) {
	c.devicesManager.ListenersManager().AsyncTrigger(boggart.DeviceEventSyslogReceive, message)
}
