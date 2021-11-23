package openweathermap

import (
	"github.com/kihamo/snitch"
)

var (
	metricCurrent     = snitch.NewGauge("current", "Current weather")
	metricTempMin     = snitch.NewGauge("temp_min", "Min daily temperature")
	metricTempMax     = snitch.NewGauge("temp_max", "Max daily temperature")
	metricTempDay     = snitch.NewGauge("temp_day", "Day temperature")
	metricTempNight   = snitch.NewGauge("temp_night", "Night temperature")
	metricTempMorning = snitch.NewGauge("temp_morning", "Morning temperature")
	metricWindSpeed   = snitch.NewGauge("wind_speed", "Wind speed")
)

func (b *Bind) Describe(ch chan<- *snitch.Description) {
	id := b.Meta().ID()
	if id == "" {
		return
	}

	metricCurrent.With("id", id).Describe(ch)
	metricTempMin.With("id", id).Describe(ch)
	metricTempMax.With("id", id).Describe(ch)
	metricTempDay.With("id", id).Describe(ch)
	metricTempNight.With("id", id).Describe(ch)
	metricTempMorning.With("id", id).Describe(ch)
	metricWindSpeed.With("id", id).Describe(ch)
}

func (b *Bind) Collect(ch chan<- snitch.Metric) {
	id := b.Meta().ID()
	if id == "" {
		return
	}

	metricCurrent.With("id", id).Collect(ch)
	metricTempMin.With("id", id).Collect(ch)
	metricTempMax.With("id", id).Collect(ch)
	metricTempDay.With("id", id).Collect(ch)
	metricTempNight.With("id", id).Collect(ch)
	metricTempMorning.With("id", id).Collect(ch)
	metricWindSpeed.With("id", id).Collect(ch)
}
