package cloud

import (
	"github.com/kihamo/snitch"
)

var (
	metricWeatherTemperature                    = snitch.NewGauge("weather_temperature_celsius", "Weather temperature in celsius")
	metricEnvironmentStateTemperatureCelsius    = snitch.NewGauge("environment_state_temperature_celsius", "Environment current temperature in celsius")
	metricHeaterHeatingFeedTemperatureCelsius   = snitch.NewGauge("heater_heating_feed_temperature_celsius", "Heating feed temperature in celsius")
	metricHeaterHeatingReturnTemperatureCelsius = snitch.NewGauge("heater_heating_return_temperature_celsius", "Heating return temperature in celsius")
	metricHeaterHeatingCircuitPressureBar       = snitch.NewGauge("heater_heating_circuit_pressure_bar", "Heater heating circuit pressure in bar")
	metricHeaterModulationPercent               = snitch.NewGauge("heater_heating_modulation_percent", "Heater modulation in percent")
)

func (b *Bind) Describe(ch chan<- *snitch.Description) {
	metricWeatherTemperature.Describe(ch)
	metricEnvironmentStateTemperatureCelsius.Describe(ch)
	metricHeaterHeatingFeedTemperatureCelsius.Describe(ch)
	metricHeaterHeatingReturnTemperatureCelsius.Describe(ch)
	metricHeaterHeatingCircuitPressureBar.Describe(ch)
	metricHeaterModulationPercent.Describe(ch)
}

func (b *Bind) Collect(ch chan<- snitch.Metric) {
	metricWeatherTemperature.Collect(ch)
	metricEnvironmentStateTemperatureCelsius.Collect(ch)
	metricHeaterHeatingFeedTemperatureCelsius.Collect(ch)
	metricHeaterHeatingReturnTemperatureCelsius.Collect(ch)
	metricHeaterHeatingCircuitPressureBar.Collect(ch)
	metricHeaterModulationPercent.Collect(ch)
}
