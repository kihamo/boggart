package internal

import (
	"github.com/kihamo/snitch"
)

func (c *Component) Metrics() snitch.Collector {
	return c.devicesManager
}
