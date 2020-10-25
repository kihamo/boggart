package mosenergosbyt

import (
	"github.com/kihamo/snitch"
)

var (
	metricBalance        = snitch.NewGauge("balance_rubles", "MosEnergoSbyt balance in rubles")
	metricServiceBalance = snitch.NewGauge("service_balance_rubles", "MosEnergoSbyt service balance in rubles")
)

func (b *Bind) Describe(ch chan<- *snitch.Description) {
	metricBalance.Describe(ch)
	metricServiceBalance.Describe(ch)
}

func (b *Bind) Collect(ch chan<- snitch.Metric) {
	metricBalance.Collect(ch)
	metricServiceBalance.Collect(ch)
}
