package smcenter

import (
	"github.com/kihamo/snitch"
)

var (
	metricAccountBalance = snitch.NewGauge("account_balance_rubles", "Account balance in rubles")
	metricMeterLastValue = snitch.NewGauge("meter", "Meter last value")
)

func (b *Bind) Describe(ch chan<- *snitch.Description) {
	metricAccountBalance.Describe(ch)
	metricMeterLastValue.Describe(ch)
}

func (b *Bind) Collect(ch chan<- snitch.Metric) {
	metricAccountBalance.Collect(ch)
	metricMeterLastValue.Collect(ch)
}
