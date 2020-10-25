package hikvision

import (
	"github.com/kihamo/snitch"
)

var (
	metricUpTime           = snitch.NewGauge("uptime_seconds", "Uptime in seconds")
	metricMemoryUsage      = snitch.NewGauge("memory_usage_bytes", "Memory usage in bytes")
	metricMemoryAvailable  = snitch.NewGauge("memory_available_bytes", "Memory available in bytes")
	metricStorageUsage     = snitch.NewGauge("storage_usage_bytes", "Storage usage in bytes")
	metricStorageAvailable = snitch.NewGauge("storage_available_bytes", "Storage available in bytes")
)

func (b *Bind) Describe(ch chan<- *snitch.Description) {
	sn := b.Meta().SerialNumber()
	if sn == "" {
		return
	}

	metricUpTime.With("serial_number", sn).Describe(ch)
	metricMemoryUsage.With("serial_number", sn).Describe(ch)
	metricMemoryAvailable.With("serial_number", sn).Describe(ch)
	metricStorageUsage.With("serial_number", sn).Describe(ch)
	metricStorageAvailable.With("serial_number", sn).Describe(ch)
}

func (b *Bind) Collect(ch chan<- snitch.Metric) {
	sn := b.Meta().SerialNumber()
	if sn == "" {
		return
	}

	metricUpTime.With("serial_number", sn).Collect(ch)
	metricMemoryUsage.With("serial_number", sn).Collect(ch)
	metricMemoryAvailable.With("serial_number", sn).Collect(ch)
	metricStorageUsage.With("serial_number", sn).Collect(ch)
	metricStorageAvailable.With("serial_number", sn).Collect(ch)
}
