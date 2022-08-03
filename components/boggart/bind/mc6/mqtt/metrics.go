package mqtt

import (
	"github.com/kihamo/snitch"
)

var (
	metricTemperature     = snitch.NewGauge("temperature", "Current temperature")
	metricHoldTemperature = snitch.NewGauge("hold_temperature", "Current hold temperature")
	metricHumidity        = snitch.NewGauge("humidity", "Current humidity")
)

func (b *Bind) Describe(ch chan<- *snitch.Description) {
	id := b.Meta().ID()
	if id == "" {
		return
	}

	metricTemperature.With("id", id).Describe(ch)
	metricHoldTemperature.With("id", id).Describe(ch)
	metricHumidity.With("id", id).Describe(ch)
}

func (b *Bind) Collect(ch chan<- snitch.Metric) {
	id := b.Meta().ID()
	if id == "" {
		return
	}

	metricTemperature.With("id", id).Collect(ch)
	metricHoldTemperature.With("id", id).Collect(ch)
	metricHumidity.With("id", id).Collect(ch)
}
