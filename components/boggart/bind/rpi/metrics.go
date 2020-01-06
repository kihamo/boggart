package rpi

import (
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/snitch"
)

var (
	metricCPUFrequentie = snitch.NewGauge(boggart.ComponentName+"_bind_pri_cpu_frequentie_hz", "CPU frequentie in Hz")
	metricTemperature   = snitch.NewGauge(boggart.ComponentName+"_bind_pri_temperature_celsius", "Core temperature of BCM2835 SoC")
	metricVoltage       = snitch.NewGauge(boggart.ComponentName+"_bind_rpi_voltage_volts", "Voltage in volts")
)

func (b *Bind) Describe(ch chan<- *snitch.Description) {
	sn := b.SerialNumber()

	metricCPUFrequentie.With("serial_number", sn).Describe(ch)
	metricTemperature.With("serial_number", sn).Describe(ch)
	metricVoltage.With("serial_number", sn).Describe(ch)
}

func (b *Bind) Collect(ch chan<- snitch.Metric) {
	sn := b.SerialNumber()

	metricCPUFrequentie.With("serial_number", sn).Collect(ch)
	metricTemperature.With("serial_number", sn).Collect(ch)
	metricVoltage.With("serial_number", sn).Collect(ch)
}
