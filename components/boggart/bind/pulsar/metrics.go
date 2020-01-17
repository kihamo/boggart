package pulsar

import (
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/snitch"
)

var (
	metricTemperatureIn    = snitch.NewGauge(boggart.ComponentName+"_bind_pulsar_temperature_in_celsius", "Pulsar temperature in in celsius")
	metricTemperatureOut   = snitch.NewGauge(boggart.ComponentName+"_bind_pulsar_temperature_out_celsius", "Pulsar temperature out in celsius")
	metricTemperatureDelta = snitch.NewGauge(boggart.ComponentName+"_bind_pulsar_temperature_delta_celsius", "Pulsar temperature delta in celsius")
	metricEnergy           = snitch.NewGauge(boggart.ComponentName+"_bind_pulsar_energy_gigacalories", "Pulsar energy in gigacalories")
	metricConsumption      = snitch.NewGauge(boggart.ComponentName+"_bind_pulsar_consumption_cubic_meters_per_hour", "Pulsar consumption in cubic meters per hour")
	metricCapacity         = snitch.NewGauge(boggart.ComponentName+"_bind_pulsar_capacity_cubic_meters", "Pulsar capacity in cubic meters")
	metricPower            = snitch.NewGauge(boggart.ComponentName+"_bind_pulsar_power_gigacalories_per_hour", "Pulsar power in gigacalories per hour")
	metricInputPulses      = snitch.NewGauge(boggart.ComponentName+"_bind_pulsar_input_pulses_count", "Pulsar input in pulses")
	metricInputVolume      = snitch.NewGauge(boggart.ComponentName+"_bind_pulsar_input_volume_cubic_meters", "Pulsar input volume in cubic meters")
)

func (b *Bind) Describe(ch chan<- *snitch.Description) {
	sn := b.address

	metricTemperatureIn.With("serial_number", sn).Describe(ch)
	metricTemperatureOut.With("serial_number", sn).Describe(ch)
	metricTemperatureDelta.With("serial_number", sn).Describe(ch)
	metricEnergy.With("serial_number", sn).Describe(ch)
	metricConsumption.With("serial_number", sn).Describe(ch)
	metricCapacity.With("serial_number", sn).Describe(ch)
	metricPower.With("serial_number", sn).Describe(ch)
	metricInputPulses.With("serial_number", sn).Describe(ch)
	metricInputVolume.With("serial_number", sn).Describe(ch)
}

func (b *Bind) Collect(ch chan<- snitch.Metric) {
	sn := b.address

	metricTemperatureIn.With("serial_number", sn).Collect(ch)
	metricTemperatureOut.With("serial_number", sn).Collect(ch)
	metricTemperatureDelta.With("serial_number", sn).Collect(ch)
	metricEnergy.With("serial_number", sn).Collect(ch)
	metricConsumption.With("serial_number", sn).Collect(ch)
	metricCapacity.With("serial_number", sn).Collect(ch)
	metricPower.With("serial_number", sn).Collect(ch)
	metricInputPulses.With("serial_number", sn).Collect(ch)
	metricInputVolume.With("serial_number", sn).Collect(ch)
}
