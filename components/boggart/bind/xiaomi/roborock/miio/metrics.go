package miio

import (
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/snitch"
)

var (
	metricBattery   = snitch.NewGauge(boggart.ComponentName+"_bind_xiaomi_roborock_miio_battery_percent", "Roborock battery in percents")
	metricCleanArea = snitch.NewGauge(boggart.ComponentName+"_bind_xiaomi_roborock_miio_clean_area_millimeters", "Roborock clean area in millimeters")
	metricCleanTime = snitch.NewGauge(boggart.ComponentName+"_bind_xiaomi_roborock_miio_clean_time_seconds", "Roborock clean time in seconds")
)

func (b *Bind) Describe(ch chan<- *snitch.Description) {
	sn := b.SerialNumber()
	if sn == "" {
		return
	}

	metricBattery.With("serial_number", sn).Describe(ch)
	metricCleanArea.With("serial_number", sn).Describe(ch)
	metricCleanTime.With("serial_number", sn).Describe(ch)
}

func (b *Bind) Collect(ch chan<- snitch.Metric) {
	sn := b.SerialNumber()
	if sn == "" {
		return
	}

	metricBattery.With("serial_number", sn).Collect(ch)
	metricCleanArea.With("serial_number", sn).Collect(ch)
	metricCleanTime.With("serial_number", sn).Collect(ch)
}
