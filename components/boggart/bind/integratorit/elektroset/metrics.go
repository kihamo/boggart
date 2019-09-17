package elektroset

import (
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/snitch"
)

var (
	metricBalance        = snitch.NewGauge(boggart.ComponentName+"_bind_elektroset_balance_rubles", "Elektroset balance in rubles")
	metricServiceBalance = snitch.NewGauge(boggart.ComponentName+"_bind_elektroset_service_balance_rubles", "Elektroset service balance in rubles")
	metricMeterValue     = snitch.NewGauge(boggart.ComponentName+"_bind_elektroset_meter_value", "Elektroset meter value")
)

func (b *Bind) Describe(ch chan<- *snitch.Description) {
	metricBalance.Describe(ch)
	metricServiceBalance.Describe(ch)
	metricMeterValue.Describe(ch)
}

func (b *Bind) Collect(ch chan<- snitch.Metric) {
	metricBalance.Collect(ch)
	metricServiceBalance.Collect(ch)
	metricMeterValue.Collect(ch)
}
