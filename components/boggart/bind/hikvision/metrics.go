package hikvision

import (
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/snitch"
)

var (
	metricUpTime           = snitch.NewGauge(boggart.ComponentName+"_bind_hikvision_uptime_seconds", "HikVision uptime in seconds")
	metricMemoryUsage      = snitch.NewGauge(boggart.ComponentName+"_bind_hikvision_memory_usage_bytes", "HikVision memory usage in bytes")
	metricMemoryAvailable  = snitch.NewGauge(boggart.ComponentName+"_bind_hikvision_memory_available_bytes", "HikVision memory available in bytes")
	metricStorageUsage     = snitch.NewGauge(boggart.ComponentName+"_bind_hikvision_storage_usage_bytes", "HikVision storage usage in bytes")
	metricStorageAvailable = snitch.NewGauge(boggart.ComponentName+"_bind_hikvision_storage_available_bytes", "HikVision storage available in bytes")
)

func (b *Bind) Describe(ch chan<- *snitch.Description) {
	sn := b.SerialNumber()
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
	sn := b.SerialNumber()
	if sn == "" {
		return
	}

	metricUpTime.With("serial_number", sn).Collect(ch)
	metricMemoryUsage.With("serial_number", sn).Collect(ch)
	metricMemoryAvailable.With("serial_number", sn).Collect(ch)
	metricStorageUsage.With("serial_number", sn).Collect(ch)
	metricStorageAvailable.With("serial_number", sn).Collect(ch)
}
