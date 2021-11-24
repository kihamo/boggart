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
	metricCurrentPm25Value   = snitch.NewGauge("current_pm25_value", "Current PM25")
	metricCurrentPm10Value   = snitch.NewGauge("current_pm10_value", "Current PM10")
	metricCurrentO3Value     = snitch.NewGauge("current_o3_value", "Current O3")
	metricCurrentNO2Value    = snitch.NewGauge("current_no2_value", "Current NO2")
	metricCurrentCOValue     = snitch.NewGauge("current_co_value", "Current CO")
	metricCurrentSO2Value    = snitch.NewGauge("current_so2_value", "Current SO2")
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
	metricCurrentPm25Value.With("id", id).Describe(ch)
	metricCurrentPm10Value.With("id", id).Describe(ch)
	metricCurrentO3Value.With("id", id).Describe(ch)
	metricCurrentNO2Value.With("id", id).Describe(ch)
	metricCurrentCOValue.With("id", id).Describe(ch)
	metricCurrentSO2Value.With("id", id).Describe(ch)
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
	metricCurrentPm25Value.With("id", id).Collect(ch)
	metricCurrentPm10Value.With("id", id).Collect(ch)
	metricCurrentO3Value.With("id", id).Collect(ch)
	metricCurrentNO2Value.With("id", id).Collect(ch)
	metricCurrentCOValue.With("id", id).Collect(ch)
	metricCurrentSO2Value.With("id", id).Collect(ch)
}
