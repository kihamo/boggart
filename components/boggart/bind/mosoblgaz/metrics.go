package mosoblgaz

import (
	"github.com/kihamo/snitch"
)

var (
	metricBalance = snitch.NewGauge("balance_rubles", "MosOblGaz balance in rubles")
)

func (b *Bind) Describe(ch chan<- *snitch.Description) {
	metricBalance.Describe(ch)
}

func (b *Bind) Collect(ch chan<- snitch.Metric) {
	metricBalance.Collect(ch)
}
