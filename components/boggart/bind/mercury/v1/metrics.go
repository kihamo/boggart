package v1

import (
	"github.com/kihamo/snitch"
)

var (
	metricTariff         = snitch.NewGauge("tariff_watts", "Tariff in watts")
	metricVoltage        = snitch.NewGauge("voltage_volts", "Voltage in volts")
	metricAmperage       = snitch.NewGauge("amperage_amperes", "Amperage in amperes")
	metricPower          = snitch.NewGauge("power_watts", "Current power in watts")
	metricBatteryVoltage = snitch.NewGauge("battery_voltage_volts", "Battery voltage in volts")
)

func (b *Bind) Describe(ch chan<- *snitch.Description) {
	address := b.config().Address

	metricTariff.With("serial_number", address).Describe(ch)
	metricVoltage.With("serial_number", address).Describe(ch)
	metricAmperage.With("serial_number", address).Describe(ch)
	metricPower.With("serial_number", address).Describe(ch)
	metricBatteryVoltage.With("serial_number", address).Describe(ch)
}

func (b *Bind) Collect(ch chan<- snitch.Metric) {
	address := b.config().Address

	metricTariff.With("serial_number", address).Collect(ch)
	metricVoltage.With("serial_number", address).Collect(ch)
	metricAmperage.With("serial_number", address).Collect(ch)
	metricPower.With("serial_number", address).Collect(ch)
	metricBatteryVoltage.With("serial_number", address).Collect(ch)
}
