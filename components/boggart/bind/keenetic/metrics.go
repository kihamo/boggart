package keenetic

import (
	"github.com/kihamo/snitch"
)

var (
	metricUpTime          = snitch.NewGauge("uptime_seconds", "Uptime in seconds")
	metricCPULoad         = snitch.NewGauge("cpu_load_percent", "CPU load in percents")
	metricMemoryUsage     = snitch.NewGauge("memory_usage_bytes", "Memory usage in bytes")
	metricMemoryAvailable = snitch.NewGauge("memory_available_bytes", "Memory available in bytes")
)

func (b *Bind) Describe(ch chan<- *snitch.Description) {
	sn := b.Meta().SerialNumber()
	if sn == "" {
		return
	}

	metricUpTime.With("serial_number", sn).Describe(ch)
	metricCPULoad.With("serial_number", sn).Describe(ch)
	metricMemoryUsage.With("serial_number", sn).Describe(ch)
	metricMemoryAvailable.With("serial_number", sn).Describe(ch)
}

func (b *Bind) Collect(ch chan<- snitch.Metric) {
	sn := b.Meta().SerialNumber()
	if sn == "" {
		return
	}

	metricUpTime.With("serial_number", sn).Collect(ch)
	metricCPULoad.With("serial_number", sn).Collect(ch)
	metricMemoryUsage.With("serial_number", sn).Collect(ch)
	metricMemoryAvailable.With("serial_number", sn).Collect(ch)
}
