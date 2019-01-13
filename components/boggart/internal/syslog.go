package internal

import (
	"context"

	"github.com/kihamo/boggart/components/boggart"
)

func (c *Component) SyslogHandler(message map[string]interface{}) {
	c.listenersManager.AsyncTrigger(context.TODO(), boggart.BindEventSyslogReceive, message)
}
