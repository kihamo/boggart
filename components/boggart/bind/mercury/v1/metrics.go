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
	metricTariff.With("serial_number", b.config.Address).Describe(ch)
	metricVoltage.With("serial_number", b.config.Address).Describe(ch)
	metricAmperage.With("serial_number", b.config.Address).Describe(ch)
	metricPower.With("serial_number", b.config.Address).Describe(ch)
	metricBatteryVoltage.With("serial_number", b.config.Address).Describe(ch)
}

func (b *Bind) Collect(ch chan<- snitch.Metric) {
	metricTariff.With("serial_number", b.config.Address).Collect(ch)
	metricVoltage.With("serial_number", b.config.Address).Collect(ch)
	metricAmperage.With("serial_number", b.config.Address).Collect(ch)
	metricPower.With("serial_number", b.config.Address).Collect(ch)
	metricBatteryVoltage.With("serial_number", b.config.Address).Collect(ch)
}
