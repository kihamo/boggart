package v1

import (
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/snitch"
)

var (
	metricTariff         = snitch.NewGauge(boggart.ComponentName+"_bind_mercury_tariff_watts", "Mercury 200 tariff in watts")
	metricVoltage        = snitch.NewGauge(boggart.ComponentName+"_bind_mercury_voltage_volts", "Mercury 200 voltage in volts")
	metricAmperage       = snitch.NewGauge(boggart.ComponentName+"_bind_mercury_amperage_amperes", "Mercury 200 amperage in amperes")
	metricPower          = snitch.NewGauge(boggart.ComponentName+"_bind_mercury_power_watts", "Mercury 200 current power in watts")
	metricBatteryVoltage = snitch.NewGauge(boggart.ComponentName+"_bind_mercury_battery_voltage_volts", "Mercury 200 battery voltage in volts")
)

func (b *Bind) Describe(ch chan<- *snitch.Description) {
	sn := b.Meta().SerialNumber()

	metricTariff.With("serial_number", sn).Describe(ch)
	metricVoltage.With("serial_number", sn).Describe(ch)
	metricAmperage.With("serial_number", sn).Describe(ch)
	metricPower.With("serial_number", sn).Describe(ch)
	metricBatteryVoltage.With("serial_number", sn).Describe(ch)
}

func (b *Bind) Collect(ch chan<- snitch.Metric) {
	sn := b.Meta().SerialNumber()

	metricTariff.With("serial_number", sn).Collect(ch)
	metricVoltage.With("serial_number", sn).Collect(ch)
	metricAmperage.With("serial_number", sn).Collect(ch)
	metricPower.With("serial_number", sn).Collect(ch)
	metricBatteryVoltage.With("serial_number", sn).Collect(ch)
}
