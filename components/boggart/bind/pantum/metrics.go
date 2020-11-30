package pantum

import (
	"github.com/kihamo/snitch"
)

var (
	metricTonerRemain = snitch.NewGauge("toner_remain_percentage", "Remaining Amount of Toner in percentage")
)

func (b *Bind) Describe(ch chan<- *snitch.Description) {
	if sn := b.Meta().SerialNumber(); sn != "" {
		metricTonerRemain.With("serial_number", sn).Describe(ch)
	}
}

func (b *Bind) Collect(ch chan<- snitch.Metric) {
	if sn := b.Meta().SerialNumber(); sn != "" {
		metricTonerRemain.With("serial_number", sn).Collect(ch)
	}
}
