package softvideo

import (
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/snitch"
)

var (
	metricBalance = snitch.NewGauge(boggart.ComponentName+"_bind_softvideo_balance_rubles", "SoftVideo balance in rubles")
)

func (b *Bind) Describe(ch chan<- *snitch.Description) {
	metricBalance.With("account", b.SerialNumber()).Describe(ch)
}

func (b *Bind) Collect(ch chan<- snitch.Metric) {
	metricBalance.With("account", b.SerialNumber()).Collect(ch)
}
