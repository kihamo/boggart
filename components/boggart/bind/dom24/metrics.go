package dom24

import (
	"github.com/kihamo/snitch"
)

var (
	metricAccountBalance = snitch.NewGauge("account_balance_rubles", "Dom24 account balance in rubles")
)

func (b *Bind) Describe(ch chan<- *snitch.Description) {
	metricAccountBalance.Describe(ch)
}

func (b *Bind) Collect(ch chan<- snitch.Metric) {
	metricAccountBalance.Collect(ch)
}
