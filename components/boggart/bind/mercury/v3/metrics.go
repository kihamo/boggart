package v3

import (
	"github.com/kihamo/snitch"
)

var (
	metricTariff   = snitch.NewGauge("tariff_watts", "Tariff in watts")
	metricVoltage  = snitch.NewGauge("voltage_volts", "Voltage in volts")
	metricAmperage = snitch.NewGauge("amperage_amperes", "Amperage in amperes")
	metricPower    = snitch.NewGauge("power_watts", "Current power in watts")
)

func (b *Bind) Describe(ch chan<- *snitch.Description) {
	sn := b.Meta().SerialNumber()
	if sn == "" {
		return
	}

	metricTariff.With("serial_number", sn).Describe(ch)
	metricVoltage.With("serial_number", sn).Describe(ch)
	metricAmperage.With("serial_number", sn).Describe(ch)
	metricPower.With("serial_number", sn).Describe(ch)
}

func (b *Bind) Collect(ch chan<- snitch.Metric) {
	sn := b.Meta().SerialNumber()
	if sn == "" {
		return
	}

	metricTariff.With("serial_number", sn).Collect(ch)
	metricVoltage.With("serial_number", sn).Collect(ch)
	metricAmperage.With("serial_number", sn).Collect(ch)
	metricPower.With("serial_number", sn).Collect(ch)
}
