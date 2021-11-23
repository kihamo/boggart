package aqicn

import (
	"github.com/kihamo/snitch"
)

var (
	metricCurrentTemperature = snitch.NewGauge("current_temperature", "Current temperature")
	metricCurrentPressure    = snitch.NewGauge("current_pressure", "Current pressure")
	metricCurrentHumidity    = snitch.NewGauge("current_humidity", "Current humidity")
	metricCurrentDewPoint    = snitch.NewGauge("current_dew_point", "Current dew point")
	metricCurrentWindSpeed   = snitch.NewGauge("current_wind_speed", "Current wind speed")
)

func (b *Bind) Describe(ch chan<- *snitch.Description) {
	id := b.Meta().ID()
	if id == "" {
		return
	}

	metricCurrentTemperature.With("id", id).Describe(ch)
	metricCurrentPressure.With("id", id).Describe(ch)
	metricCurrentHumidity.With("id", id).Describe(ch)
	metricCurrentDewPoint.With("id", id).Describe(ch)
	metricCurrentWindSpeed.With("id", id).Describe(ch)
}

func (b *Bind) Collect(ch chan<- snitch.Metric) {
	id := b.Meta().ID()
	if id == "" {
		return
	}

	metricCurrentTemperature.With("id", id).Collect(ch)
	metricCurrentPressure.With("id", id).Collect(ch)
	metricCurrentHumidity.With("id", id).Collect(ch)
	metricCurrentDewPoint.With("id", id).Collect(ch)
	metricCurrentWindSpeed.With("id", id).Collect(ch)
}
