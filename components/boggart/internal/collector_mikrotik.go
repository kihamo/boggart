package internal

import (
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/providers/mikrotik"
	"github.com/kihamo/snitch"
)

const (
	MetricMikrotikWifiClients = boggart.ComponentName + "_mikrotik_wifi_clients_total"
)

var (
	metricMikrotikWifiClients = snitch.NewGauge(MetricMikrotikWifiClients, "Mikrotik wifi clients online")
)

func (c *MetricsCollector) UpdaterMikrotik() error {
	client, err := mikrotik.NewClient(
		c.component.config.GetString(boggart.ConfigMikrotikAddress),
		c.component.config.GetString(boggart.ConfigMikrotikUsername),
		c.component.config.GetString(boggart.ConfigMikrotikPassword),
		c.component.config.GetDuration(boggart.ConfigMikrotikTimeout))

	if err != nil {
		return err
	}

	clients, err := client.WifiClients()
	if err != nil {
		return err
	}

	metricMikrotikWifiClients.Set(float64(len(clients)))

	return nil
}

func (c *MetricsCollector) DescribeMikrotik(ch chan<- *snitch.Description) {
	metricMikrotikWifiClients.Describe(ch)
}

func (c *MetricsCollector) CollectMikrotik(ch chan<- snitch.Metric) {
	metricMikrotikWifiClients.Collect(ch)
}
