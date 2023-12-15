package myheat

import (
	"github.com/kihamo/snitch"
)

var (
	metricSensorValue                         = snitch.NewGauge("sensor_value", "Sensor value")
	metricEnvironmentStateTemperatureCelsius  = snitch.NewGauge("environment_state_temperature_celsius", "Environment current temperature in celsius")
	metricHeaterHeatingFlowTemperatureCelsius = snitch.NewGauge("heater_heating_flow_temperature_celsius", "Heating flow temperature in celsius")
	metricHeaterHeatingCircuitPressureBar     = snitch.NewGauge("heater_heating_circuit_pressure_bar", "Heater heating circuit pressure in bar")
)

func (b *Bind) Describe(ch chan<- *snitch.Description) {
	sn := b.Meta().SerialNumber()
	if sn == "" {
		return
	}

	metricSensorValue.With("serial_number", sn).Describe(ch)
	metricEnvironmentStateTemperatureCelsius.With("serial_number", sn).Describe(ch)
	metricHeaterHeatingFlowTemperatureCelsius.With("serial_number", sn).Describe(ch)
	metricHeaterHeatingCircuitPressureBar.With("serial_number", sn).Describe(ch)
}

func (b *Bind) Collect(ch chan<- snitch.Metric) {
	sn := b.Meta().SerialNumber()
	if sn == "" {
		return
	}

	metricSensorValue.With("serial_number", sn).Collect(ch)
	metricEnvironmentStateTemperatureCelsius.With("serial_number", sn).Collect(ch)
	metricHeaterHeatingFlowTemperatureCelsius.With("serial_number", sn).Collect(ch)
	metricHeaterHeatingCircuitPressureBar.With("serial_number", sn).Collect(ch)
}
