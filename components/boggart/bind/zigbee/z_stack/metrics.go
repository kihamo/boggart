package zstack

import (
	"github.com/kihamo/snitch"
)

var (
	metricLinkQuality    = snitch.NewGauge("link_quality", "Indicates the link quality measured during reception")
	metricBatteryPercent = snitch.NewGauge("battery_percent", "Battery voltage in percent")
	metricBatteryVoltage = snitch.NewGauge("battery_voltage", "Battery voltage")
)

func (b *Bind) Describe(ch chan<- *snitch.Description) {
	sn := b.Meta().SerialNumber()
	if sn == "" {
		return
	}

	metricLinkQuality.With("serial_number", sn).Describe(ch)
	metricBatteryPercent.With("serial_number", sn).Describe(ch)
	metricBatteryVoltage.With("serial_number", sn).Describe(ch)
}

func (b *Bind) Collect(ch chan<- snitch.Metric) {
	sn := b.Meta().SerialNumber()
	if sn == "" {
		return
	}

	metricLinkQuality.With("serial_number", sn).Collect(ch)
	metricBatteryPercent.With("serial_number", sn).Collect(ch)
	metricBatteryVoltage.With("serial_number", sn).Collect(ch)
}
