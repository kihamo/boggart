package internal

import (
	"github.com/kihamo/snitch"
)

func (c *Component) Metrics() snitch.Collector {
	<-c.application.ReadyComponent(c.Name())
	return c.manager
}
