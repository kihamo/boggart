package dom24

import (
	"github.com/kihamo/snitch"
)

var (
	metricAccountBalance = snitch.NewGauge("account_balance_rubles", "Dom24 account balance in rubles")
	metricMeterLastValue = snitch.NewGauge("meter", "Dom24 meter last value")
)

func (b *Bind) Describe(ch chan<- *snitch.Description) {
	metricAccountBalance.Describe(ch)
	metricMeterLastValue.Describe(ch)
}

func (b *Bind) Collect(ch chan<- snitch.Metric) {
	metricAccountBalance.Collect(ch)
	metricMeterLastValue.Collect(ch)
}
