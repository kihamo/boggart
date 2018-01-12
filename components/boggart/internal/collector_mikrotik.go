package internal

import (
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/providers/mikrotik"
)

func (c *MetricsCollector) CollectMikrotik() error {
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
