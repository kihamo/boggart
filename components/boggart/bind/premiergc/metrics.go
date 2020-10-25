package premiergc

import (
	"github.com/kihamo/snitch"
)

var (
	metricBalance = snitch.NewGauge("balance_rubles", "Premier GC balance in rubles")
)

func (b *Bind) Describe(ch chan<- *snitch.Description) {
	sn := b.Meta().SerialNumber()
	if sn == "" {
		return
	}

	metricBalance.With("contract", sn).Describe(ch)
}

func (b *Bind) Collect(ch chan<- snitch.Metric) {
	sn := b.Meta().SerialNumber()
	if sn == "" {
		return
	}

	metricBalance.With("contract", sn).Collect(ch)
}
