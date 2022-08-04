package mqtt

import (
	"github.com/kihamo/snitch"
)

var (
	metricTemperature    = snitch.NewGauge("temperature", "Current temperature")
	metricSetTemperature = snitch.NewGauge("set_temperature", "Current setting temperature")
	metricHumidity       = snitch.NewGauge("humidity", "Current humidity")
)

func (b *Bind) Describe(ch chan<- *snitch.Description) {
	id := b.Meta().ID()
	if id == "" {
		return
	}

	metricTemperature.With("id", id).Describe(ch)
	metricSetTemperature.With("id", id).Describe(ch)
	metricHumidity.With("id", id).Describe(ch)
}

func (b *Bind) Collect(ch chan<- snitch.Metric) {
	id := b.Meta().ID()
	if id == "" {
		return
	}

	metricTemperature.With("id", id).Collect(ch)
	metricSetTemperature.With("id", id).Collect(ch)
	metricHumidity.With("id", id).Collect(ch)
}
