package octoprint

import (
	"github.com/kihamo/snitch"
)

var (
	metricDeviceTemperatureActual = snitch.NewGauge("temperature_actual", "Current temperature")
	metricDeviceTemperatureTarget = snitch.NewGauge("temperature_target", "Target temperature")
)

func (b *Bind) Describe(ch chan<- *snitch.Description) {
	address := b.config().Address.Host

	metricDeviceTemperatureActual.With("address", address).Describe(ch)
	metricDeviceTemperatureTarget.With("address", address).Describe(ch)
}

func (b *Bind) Collect(ch chan<- snitch.Metric) {
	address := b.config().Address.Host

	metricDeviceTemperatureActual.With("address", address).Collect(ch)
	metricDeviceTemperatureTarget.With("address", address).Collect(ch)
}
