package rkcm

import (
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/snitch"
)

var (
	metricBalance = snitch.NewGauge(boggart.ComponentName+"_bind_rkcm_balance_rubles", "RKCM balance in rubles")
)

func (b *Bind) Describe(ch chan<- *snitch.Description) {
	metricBalance.Describe(ch)
}

func (b *Bind) Collect(ch chan<- snitch.Metric) {
	metricBalance.Collect(ch)
}
