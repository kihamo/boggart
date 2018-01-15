package internal

import (
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/snitch"
)

const (
	MetricDoorsEntranceStatus = boggart.ComponentName + "_doors_entrance_status"
)

var (
	metricDoorsEntranceStatus = snitch.NewGauge(MetricDoorsEntranceStatus, "Entrance door status")
)

func (c *MetricsCollector) DescribeDoors(ch chan<- *snitch.Description) {
	metricDoorsEntranceStatus.Describe(ch)
}

func (c *MetricsCollector) CollectDoors(ch chan<- snitch.Metric) {
	if c.component.DoorEntrance().IsOpen() {
		metricDoorsEntranceStatus.Set(1)
	} else {
		metricDoorsEntranceStatus.Set(0)
	}

	metricDoorsEntranceStatus.Collect(ch)
}
