package internal

import (
	"github.com/kihamo/boggart/components/syslog"
	"github.com/kihamo/snitch"
)

const (
	MetricHandled = syslog.ComponentName + "_handled_total"
)

var (
	metricHandled = snitch.NewCounter(MetricHandled, "Total handled")
)

func (c *Component) Metrics() snitch.Collector {
	return metricHandled
}
