package ds18b20

import (
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/snitch"
)

var (
	metricValue = snitch.NewGauge(boggart.ComponentName+"_bind_ds18b20_value_celsius", "DS18B20 sensor value in celsius")
)

func (b *Bind) Describe(ch chan<- *snitch.Description) {
	metricValue.With("serial_number", b.SerialNumber()).Describe(ch)
}

func (b *Bind) Collect(ch chan<- snitch.Metric) {
	metricValue.With("serial_number", b.SerialNumber()).Collect(ch)
}
