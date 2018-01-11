package internal

import (
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/snitch"
)

const (
	MetricPulsarTemperatureIn    = boggart.ComponentName + "_pulsar_temperature_in_celsius"
	MetricPulsarTemperatureOut   = boggart.ComponentName + "_pulsar_temperature_out_celsius"
	MetricPulsarTemperatureDelta = boggart.ComponentName + "_pulsar_temperature_delta_celsius"
	MetricSoftVideoBalance       = boggart.ComponentName + "_softvideo_balance_rubles_total"
	MetricMikrotikWifiClients    = boggart.ComponentName + "_mikrotik_wifi_clients_total"
)

var (
	metricPulsarTemperatureIn    = snitch.NewGauge(MetricPulsarTemperatureIn, "Pulsar temperature in")
	metricPulsarTemperatureOut   = snitch.NewGauge(MetricPulsarTemperatureOut, "Pulsar temperature out")
	metricPulsarTemperatureDelta = snitch.NewGauge(MetricPulsarTemperatureDelta, "Pulsar temperature delta")
	metricSoftVideoBalance       = snitch.NewGauge(MetricSoftVideoBalance, "SoftVideo balance in rubles")
	metricMikrotikWifiClients    = snitch.NewGauge(MetricMikrotikWifiClients, "Mikrotik wifi clients online")
)

func (c *Component) Metrics() snitch.Collector {
	return c.collector
}
