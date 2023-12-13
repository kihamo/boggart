package myheat

import (
	"github.com/kihamo/snitch"
)

var (
	metricSensorValue = snitch.NewGauge("sensor_value", "Sensor value")
)

func (b *Bind) Describe(ch chan<- *snitch.Description) {
	sn := b.Meta().SerialNumber()
	if sn == "" {
		return
	}

	metricSensorValue.With("serial_number", sn).Describe(ch)
}

func (b *Bind) Collect(ch chan<- snitch.Metric) {
	sn := b.Meta().SerialNumber()
	if sn == "" {
		return
	}

	metricSensorValue.With("serial_number", sn).Collect(ch)
}
