package rkcm

import (
	"github.com/kihamo/snitch"
)

var (
	metricBalance    = snitch.NewGauge("balance_rubles", "RKCM balance in rubles")
	metricMeterValue = snitch.NewGauge("meter_value", "RKCM meter value")
)

func (b *Bind) Describe(ch chan<- *snitch.Description) {
	metricBalance.Describe(ch)
	metricMeterValue.Describe(ch)
}

func (b *Bind) Collect(ch chan<- snitch.Metric) {
	metricBalance.Collect(ch)
	metricMeterValue.Collect(ch)
}
