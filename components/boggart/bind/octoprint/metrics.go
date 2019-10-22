package octoprint

import (
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/snitch"
)

var (
	metricDeviceTemperatureActual = snitch.NewGauge(boggart.ComponentName+"_bind_octoprint_temperature_actual", "Current temperature")
	metricDeviceTemperatureTarget = snitch.NewGauge(boggart.ComponentName+"_bind_octoprint_temperature_target", "Target temperature")
)

func (b *Bind) Describe(ch chan<- *snitch.Description) {
	metricDeviceTemperatureActual.With("address", b.config.Address.Host).Describe(ch)
	metricDeviceTemperatureTarget.With("address", b.config.Address.Host).Describe(ch)
}

func (b *Bind) Collect(ch chan<- snitch.Metric) {
	metricDeviceTemperatureActual.With("address", b.config.Address.Host).Collect(ch)
	metricDeviceTemperatureTarget.With("address", b.config.Address.Host).Collect(ch)
}
