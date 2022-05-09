package apcupsd

import (
	"github.com/kihamo/snitch"
)

var (
	metricLoad           = snitch.NewGauge("load_percent", "Load on UPS (percent of full)")
	metricInputVoltage   = snitch.NewGauge("input_voltage_volts", "Input voltage volts")
	metricBatteryCharge  = snitch.NewGauge("battery_charge_percent", "Battery charge (percent of full)")
	metricBatteryRuntime = snitch.NewGauge("battery_runtime_seconds", "Battery runtime seconds")
	metricBatteryVoltage = snitch.NewGauge("battery_voltage_volts", "Battery voltage volts")
)

func (b *Bind) Describe(ch chan<- *snitch.Description) {
	sn := b.Meta().SerialNumber()
	if sn == "" {
		return
	}

	metricLoad.With("serial_number", sn).Describe(ch)
	metricInputVoltage.With("serial_number", sn).Describe(ch)
	metricBatteryCharge.With("serial_number", sn).Describe(ch)
	metricBatteryRuntime.With("serial_number", sn).Describe(ch)
	metricBatteryVoltage.With("serial_number", sn).Describe(ch)
}

func (b *Bind) Collect(ch chan<- snitch.Metric) {
	sn := b.Meta().SerialNumber()
	if sn == "" {
		return
	}

	metricLoad.With("serial_number", sn).Collect(ch)
	metricInputVoltage.With("serial_number", sn).Collect(ch)
	metricBatteryCharge.With("serial_number", sn).Collect(ch)
	metricBatteryRuntime.With("serial_number", sn).Collect(ch)
	metricBatteryVoltage.With("serial_number", sn).Collect(ch)
}
