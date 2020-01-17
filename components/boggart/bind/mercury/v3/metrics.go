package v3

import (
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/snitch"
)

var (
	metricTariff   = snitch.NewGauge(boggart.ComponentName+"_bind_mercury_tariff_watts", "Mercury 230 tariff in watts")
	metricVoltage  = snitch.NewGauge(boggart.ComponentName+"_bind_mercury_voltage_volts", "Mercury 230 voltage in volts")
	metricAmperage = snitch.NewGauge(boggart.ComponentName+"_bind_mercury_amperage_amperes", "Mercury 230 amperage in amperes")
	metricPower    = snitch.NewGauge(boggart.ComponentName+"_bind_mercury_power_watts", "Mercury 230 current power in watts")
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
