package modbus

import (
	"github.com/kihamo/snitch"
)

var (
	metricRoomTemperature  = snitch.NewGauge("room_temperature", "Current room temperature")
	metricFlourTemperature = snitch.NewGauge("flour_temperature", "Current floor temperature")
	metricHumidity         = snitch.NewGauge("humidity", "Current humidity")
)

func (b *Bind) Describe(ch chan<- *snitch.Description) {
	id := b.Meta().ID()
	if id == "" {
		return
	}

	metricRoomTemperature.With("id", id).Describe(ch)
	metricFlourTemperature.With("id", id).Describe(ch)
	metricHumidity.With("id", id).Describe(ch)
}

func (b *Bind) Collect(ch chan<- snitch.Metric) {
	id := b.Meta().ID()
	if id == "" {
		return
	}

	metricRoomTemperature.With("id", id).Collect(ch)
	metricFlourTemperature.With("id", id).Collect(ch)
	metricHumidity.With("id", id).Collect(ch)
}
