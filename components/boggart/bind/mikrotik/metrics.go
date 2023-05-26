package mikrotik

import (
	"github.com/kihamo/snitch"
)

var (
	metricTrafficReceivedBytes = snitch.NewGauge("traffic_received_bytes", "Traffic received bytes")
	metricTrafficSentBytes     = snitch.NewGauge("traffic_sent_bytes", "Traffic sent bytes")
	metricWifiClients          = snitch.NewGauge("wifi_clients_total", "Wifi clients online")
	metricCPULoad              = snitch.NewGauge("cpu_load_percent", "CPU load in percents")
	metricMemoryUsage          = snitch.NewGauge("memory_usage_bytes", "Memory usage in bytes")
	metricMemoryAvailable      = snitch.NewGauge("memory_available_bytes", "Memory available in bytes")
	metricStorageUsage         = snitch.NewGauge("storage_usage_bytes", "Storage usage in bytes")
	metricStorageAvailable     = snitch.NewGauge("storage_available_bytes", "Storage available in bytes")
	metricDiskUsage            = snitch.NewGauge("disk_usage_bytes", "Disk usage in bytes")
	metricDiskAvailable        = snitch.NewGauge("disk_available_bytes", "Disk available in bytes")
	metricVoltage              = snitch.NewGauge("voltage_volt", "Voltage")
	metricTemperature          = snitch.NewGauge("temperature_celsius", "Temperature")
	metricUPSLoad              = snitch.NewGauge("ups_load_percent", "Load on UPS (percent of full)")
	metricUPSInputVoltage      = snitch.NewGauge("ups_input_voltage_volts", "Input voltage volts")
	metricUPSBatteryCharge     = snitch.NewGauge("ups_battery_charge_percent", "Battery charge (percent of full)")
	metricUPSBatteryRuntime    = snitch.NewGauge("ups_battery_runtime_seconds", "Battery runtime seconds")
	metricUPSBatteryVoltage    = snitch.NewGauge("ups_battery_voltage_volts", "Battery voltage volts")
)

func (b *Bind) Describe(ch chan<- *snitch.Description) {
	sn := b.Meta().SerialNumber()
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
	metricUPSLoad.With("serial_number", sn).Describe(ch)
	metricUPSInputVoltage.With("serial_number", sn).Describe(ch)
	metricUPSBatteryCharge.With("serial_number", sn).Describe(ch)
	metricUPSBatteryRuntime.With("serial_number", sn).Describe(ch)
	metricUPSBatteryVoltage.With("serial_number", sn).Describe(ch)
}

func (b *Bind) Collect(ch chan<- snitch.Metric) {
	sn := b.Meta().SerialNumber()
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
	metricUPSLoad.With("serial_number", sn).Collect(ch)
	metricUPSInputVoltage.With("serial_number", sn).Collect(ch)
	metricUPSBatteryCharge.With("serial_number", sn).Collect(ch)
	metricUPSBatteryRuntime.With("serial_number", sn).Collect(ch)
	metricUPSBatteryVoltage.With("serial_number", sn).Collect(ch)
}
