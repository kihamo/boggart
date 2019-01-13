package mikrotik

import (
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/snitch"
)

var (
	metricTrafficReceivedBytes = snitch.NewGauge(boggart.ComponentName+"_bind_mikrotik_traffic_received_bytes", "Bind traffic received bytes")
	metricTrafficSentBytes     = snitch.NewGauge(boggart.ComponentName+"_bind_mikrotik_traffic_sent_bytes", "Bind traffic sent bytes")
	metricWifiClients          = snitch.NewGauge(boggart.ComponentName+"_bind_mikrotik_wifi_clients_total", "Bind wifi clients online")
	metricCPULoad              = snitch.NewGauge(boggart.ComponentName+"_bind_mikrotik_cpu_load_percent", "CPU load in percents")
	metricMemoryUsage          = snitch.NewGauge(boggart.ComponentName+"_bind_mikrotik_memory_usage_bytes", "Memory usage in Bind router")
	metricMemoryAvailable      = snitch.NewGauge(boggart.ComponentName+"_bind_mikrotik_memory_available_bytes", "Memory available in Bind router")
	metricStorageUsage         = snitch.NewGauge(boggart.ComponentName+"_bind_mikrotik_storage_usage_bytes", "Storage usage in Bind router")
	metricStorageAvailable     = snitch.NewGauge(boggart.ComponentName+"_bind_mikrotik_storage_available_bytes", "Storage available in Bind router")
	metricDiskUsage            = snitch.NewGauge(boggart.ComponentName+"_bind_mikrotik_disk_usage_bytes", "Disk usage in Bind router")
	metricDiskAvailable        = snitch.NewGauge(boggart.ComponentName+"_bind_mikrotik_disk_available_bytes", "Disk available in Bind router")
	metricVoltage              = snitch.NewGauge(boggart.ComponentName+"_bind_mikrotik_voltage_volt", "Voltage")
	metricTemperature          = snitch.NewGauge(boggart.ComponentName+"_bind_mikrotik_temperature_celsius", "Temperature")
)

func (b *Bind) Describe(ch chan<- *snitch.Description) {
	sn := b.SerialNumber()
	if sn == "" {
		return
	}

	metricTrafficReceivedBytes.With("serial_number", sn).Describe(ch)
	metricTrafficSentBytes.With("serial_number", sn).Describe(ch)
	metricWifiClients.With("serial_number", sn).Describe(ch)
	metricCPULoad.With("serial_number", sn).Describe(ch)
	metricMemoryUsage.With("serial_number", sn).Describe(ch)
	metricMemoryAvailable.With("serial_number", sn).Describe(ch)
	metricStorageUsage.With("serial_number", sn).Describe(ch)
	metricStorageAvailable.With("serial_number", sn).Describe(ch)
	metricDiskUsage.With("serial_number", sn).Describe(ch)
	metricDiskAvailable.With("serial_number", sn).Describe(ch)
	metricVoltage.With("serial_number", sn).Describe(ch)
	metricTemperature.With("serial_number", sn).Describe(ch)
}

func (b *Bind) Collect(ch chan<- snitch.Metric) {
	sn := b.SerialNumber()
	if sn == "" {
		return
	}

	metricTrafficReceivedBytes.With("serial_number", sn).Collect(ch)
	metricTrafficSentBytes.With("serial_number", sn).Collect(ch)
	metricWifiClients.With("serial_number", sn).Collect(ch)
	metricCPULoad.With("serial_number", sn).Collect(ch)
	metricMemoryUsage.With("serial_number", sn).Collect(ch)
	metricMemoryAvailable.With("serial_number", sn).Collect(ch)
	metricStorageUsage.With("serial_number", sn).Collect(ch)
	metricStorageAvailable.With("serial_number", sn).Collect(ch)
	metricDiskUsage.With("serial_number", sn).Collect(ch)
	metricDiskAvailable.With("serial_number", sn).Collect(ch)
	metricVoltage.With("serial_number", sn).Collect(ch)
	metricTemperature.With("serial_number", sn).Collect(ch)
}
