package softvideo

import (
	"github.com/kihamo/snitch"
)

var (
	metricBalance = snitch.NewGauge("balance_rubles", "SoftVideo balance in rubles")
)

func (b *Bind) Describe(ch chan<- *snitch.Description) {
	metricBalance.With("account", b.config().Login).Describe(ch)
}

func (b *Bind) Collect(ch chan<- snitch.Metric) {
	metricBalance.With("account", b.config().Login).Collect(ch)
}
