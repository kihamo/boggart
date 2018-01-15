package internal

import (
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/providers/doors"
	"github.com/kihamo/snitch"
)

const (
	MetricDoorsEntranceStatus = boggart.ComponentName + "_doors_entrance_status"
)

var (
	metricDoorsEntranceStatus = snitch.NewGauge(MetricDoorsEntranceStatus, "Entrance door status")
)

func (c *MetricsCollector) UpdaterEntranceDoorsCallback(status bool) {
	if status == doors.OPEN {
		metricDoorsEntranceStatus.Set(1)
	} else {
		metricDoorsEntranceStatus.Set(0)
	}
}

func (c *MetricsCollector) UpdaterDoors() error {
	if !c.component.config.GetBool(boggart.ConfigDoorsEnabled) {
		return nil
	}

	door, err := doors.NewDoor(c.component.config.GetInt(boggart.ConfigDoorsEntrancePin))
	if err != nil {
		return err
	}

	if door.IsOpen() {
		metricDoorsEntranceStatus.Set(1)
	} else {
		metricDoorsEntranceStatus.Set(0)
	}

	return nil
}

func (c *MetricsCollector) DescribeDoors(ch chan<- *snitch.Description) {
	metricDoorsEntranceStatus.Describe(ch)
}

func (c *MetricsCollector) CollectDoors(ch chan<- snitch.Metric) {
	metricDoorsEntranceStatus.Collect(ch)
}
