package owntracks

import (
	"github.com/kihamo/snitch"
)

var (
	metricBatteryLevel = snitch.NewGauge("battery_level_percentage", "Device battery level")
	metricVelocity     = snitch.NewGauge("velocity_kmh", "Velocity in kmh")
)

func (b *Bind) Describe(ch chan<- *snitch.Description) {
	metricBatteryLevel.Describe(ch)
	metricVelocity.Describe(ch)
}

func (b *Bind) Collect(ch chan<- snitch.Metric) {
	metricBatteryLevel.Collect(ch)
	metricVelocity.Collect(ch)
}
