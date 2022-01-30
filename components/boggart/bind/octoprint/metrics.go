package octoprint

import (
	"github.com/kihamo/snitch"
)

var (
	metricDeviceTemperatureActual = snitch.NewGauge("temperature_actual", "Current temperature")
	metricDeviceTemperatureTarget = snitch.NewGauge("temperature_target", "Target temperature")
)

func (b *Bind) Describe(ch chan<- *snitch.Description) {
	id := b.Meta().ID()

	metricDeviceTemperatureActual.With("id", id).Describe(ch)
	metricDeviceTemperatureTarget.With("id", id).Describe(ch)
}

func (b *Bind) Collect(ch chan<- snitch.Metric) {
	id := b.Meta().ID()

	metricDeviceTemperatureActual.With("id", id).Collect(ch)
	metricDeviceTemperatureTarget.With("id", id).Collect(ch)
}
