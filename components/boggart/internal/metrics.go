package internal

import (
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/snitch"
)

const (
	MetricPulsarTemperatureIn     = boggart.ComponentName + "_pulsar_temperature_in_celsius"
	MetricPulsarTemperatureOut    = boggart.ComponentName + "_pulsar_temperature_out_celsius"
	MetricPulsarTemperatureDelta  = boggart.ComponentName + "_pulsar_temperature_delta_celsius"
	MetricPulsarEnergy            = boggart.ComponentName + "_pulsar_energy_gigacolories"
	MetricPulsarConsumption       = boggart.ComponentName + "_pulsar_consumption_cubic_metres_per_hour"
	MetricPulsarColdWaterCapacity = boggart.ComponentName + "_pulsar_cold_water_capacity_cubic_metres"
	MetricPulsarHotWaterCapacity  = boggart.ComponentName + "_pulsar_hot_water_capacity_cubic_metres"
	MetricSoftVideoBalance        = boggart.ComponentName + "_softvideo_balance_rubles_total"
	MetricMikrotikWifiClients     = boggart.ComponentName + "_mikrotik_wifi_clients_total"
)

var (
	metricPulsarTemperatureIn     = snitch.NewGauge(MetricPulsarTemperatureIn, "Pulsar temperature in")
	metricPulsarTemperatureOut    = snitch.NewGauge(MetricPulsarTemperatureOut, "Pulsar temperature out")
	metricPulsarTemperatureDelta  = snitch.NewGauge(MetricPulsarTemperatureDelta, "Pulsar temperature delta")
	metricPulsarEnergy            = snitch.NewGauge(MetricPulsarEnergy, "Pulsar energy")
	metricPulsarConsumption       = snitch.NewGauge(MetricPulsarConsumption, "Pulsar consumption")
	metricPulsarColdWaterCapacity = snitch.NewGauge(MetricPulsarColdWaterCapacity, "Pulsar capacity of cold water")
	metricPulsarHotWaterCapacity  = snitch.NewGauge(MetricPulsarHotWaterCapacity, "Pulsar capacity of hot water")
	metricSoftVideoBalance        = snitch.NewGauge(MetricSoftVideoBalance, "SoftVideo balance in rubles")
	metricMikrotikWifiClients     = snitch.NewGauge(MetricMikrotikWifiClients, "Mikrotik wifi clients online")
)

type MetricsCollector struct {
	component *Component
}

func NewMetricsCollector(component *Component) *MetricsCollector {
	return &MetricsCollector{
		component: component,
	}
}

func (c *MetricsCollector) Describe(ch chan<- *snitch.Description) {
	metricPulsarTemperatureIn.Describe(ch)
	metricPulsarTemperatureOut.Describe(ch)
	metricPulsarTemperatureDelta.Describe(ch)
	metricPulsarEnergy.Describe(ch)
	metricPulsarConsumption.Describe(ch)
	metricPulsarColdWaterCapacity.Describe(ch)
	metricPulsarHotWaterCapacity.Describe(ch)
	metricSoftVideoBalance.Describe(ch)
	metricMikrotikWifiClients.Describe(ch)
}

func (c *MetricsCollector) Collect(ch chan<- snitch.Metric) {
	metricPulsarTemperatureIn.Collect(ch)
	metricPulsarTemperatureOut.Collect(ch)
	metricPulsarTemperatureDelta.Collect(ch)
	metricPulsarEnergy.Collect(ch)
	metricPulsarConsumption.Collect(ch)
	metricPulsarColdWaterCapacity.Collect(ch)
	metricPulsarHotWaterCapacity.Collect(ch)
	metricSoftVideoBalance.Collect(ch)
	metricMikrotikWifiClients.Collect(ch)
}

func (c *Component) Metrics() snitch.Collector {
	return c.collector
}
