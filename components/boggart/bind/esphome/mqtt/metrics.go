package mqtt

import (
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/snitch"
)

var (
	metricState = snitch.NewGauge(boggart.ComponentName+"_bind_esphome_state", "ESPHome component state")
)

func (b *Bind) Describe(ch chan<- *snitch.Description) {
	mac := b.Meta().MAC()
	if mac == nil {
		return
	}

	metricState.With("mac", mac.String()).Describe(ch)
}

func (b *Bind) Collect(ch chan<- snitch.Metric) {
	mac := b.Meta().MAC()
	if mac == nil {
		return
	}

	metricState.With("mac", mac.String()).Collect(ch)
}
