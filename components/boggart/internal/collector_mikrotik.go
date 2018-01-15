package internal

import (
	"strconv"
	"strings"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/providers/mikrotik"
	"github.com/kihamo/snitch"
)

const (
	MetricMikrotikTrafficReceivedBytes = boggart.ComponentName + "_mikrotik_traffic_received_bytes"
	MetricMikrotikTrafficSentBytes     = boggart.ComponentName + "_mikrotik_traffic_sent_bytes"
	MetricMikrotikWifiClients          = boggart.ComponentName + "_mikrotik_wifi_clients_total"
)

var (
	metricMikrotikTrafficReceivedBytes = snitch.NewGauge(MetricMikrotikTrafficReceivedBytes, "Mikrotik traffic received bytes")
	metricMikrotikTrafficSentBytes     = snitch.NewGauge(MetricMikrotikTrafficSentBytes, "Mikrotik traffic sent bytes")
	metricMikrotikWifiClients          = snitch.NewGauge(MetricMikrotikWifiClients, "Mikrotik wifi clients online")
)

func (c *MetricsCollector) DescribeMikrotik(ch chan<- *snitch.Description) {
	metricMikrotikTrafficReceivedBytes.Describe(ch)
	metricMikrotikTrafficSentBytes.Describe(ch)
	metricMikrotikWifiClients.Describe(ch)
}

func (c *MetricsCollector) CollectMikrotik(ch chan<- snitch.Metric) {
	client, err := mikrotik.NewClient(
		c.component.config.GetString(boggart.ConfigMikrotikAddress),
		c.component.config.GetString(boggart.ConfigMikrotikUsername),
		c.component.config.GetString(boggart.ConfigMikrotikPassword),
		c.component.config.GetDuration(boggart.ConfigMikrotikTimeout))

	if err != nil {
		c.component.logger.Error("Create client failed", map[string]interface{}{
			"provider": "mikrotik",
			"error":    err.Error(),
		})
		return
	}

	clients, err := client.WifiClients()
	if err != nil {
		c.component.logger.Error("Get wifi clients failed", map[string]interface{}{
			"provider": "mikrotik",
			"error":    err.Error(),
		})
		return
	}

	metricMikrotikWifiClients.Set(float64(len(clients)))
	metricMikrotikWifiClients.Collect(ch)

	for _, client := range clients {
		bytes := strings.Split(client["bytes"], ",")
		if len(bytes) != 2 {
			c.component.logger.Errorf("Wrong value of bytes field %s", client["bytes"], map[string]interface{}{
				"provider": "mikrotik",
				"error":    err.Error(),
			})
			return
		}

		sent, err := strconv.ParseFloat(bytes[0], 64)
		if err != nil {
			c.component.logger.Errorf("Failed convert sent bytes to float64 %s", bytes[0], map[string]interface{}{
				"interface": client["interface"],
				"mac":       client["mac-address"],
				"provider":  "mikrotik",
				"error":     err.Error(),
			})
			return
		}

		received, err := strconv.ParseFloat(bytes[1], 64)
		if err != nil {
			c.component.logger.Errorf("Failed convert received bytes to float64 %s", bytes[1], map[string]interface{}{
				"interface": client["interface"],
				"mac":       client["mac-address"],
				"provider":  "mikrotik",
				"error":     err.Error(),
			})
			return
		}

		metricMikrotikTrafficReceivedBytes.With(
			"interface", client["interface"],
			"mac", client["mac-address"]).Set(received)
		metricMikrotikTrafficSentBytes.With(
			"interface", client["interface"],
			"mac", client["mac-address"]).Set(sent)
	}

	stats, err := client.EthernetStats()
	if err != nil {
		c.component.logger.Error("Get ethernet stats failed", map[string]interface{}{
			"provider": "mikrotik",
			"error":    err.Error(),
		})
		return
	}

	for _, stat := range stats {
		sent, err := strconv.ParseFloat(stat["tx-byte"], 64)
		if err != nil {
			c.component.logger.Errorf("Failed convert sent bytes to float64 %s", stat["tx-byte"], map[string]interface{}{
				"interface": stat["name"],
				"mac":       stat["mac-address"],
				"provider":  "mikrotik",
				"error":     err.Error(),
			})
			return
		}

		received, err := strconv.ParseFloat(stat["rx-byte"], 64)
		if err != nil {
			c.component.logger.Errorf("Failed convert received bytes to float64 %s", stat["rx-byte"], map[string]interface{}{
				"interface": stat["name"],
				"mac":       stat["mac-address"],
				"provider":  "mikrotik",
				"error":     err.Error(),
			})
			return
		}

		metricMikrotikTrafficReceivedBytes.With(
			"interface", stat["name"],
			"mac", stat["mac-address"]).Set(received)
		metricMikrotikTrafficSentBytes.With(
			"interface", stat["name"],
			"mac", stat["mac-address"]).Set(sent)
	}

	metricMikrotikTrafficReceivedBytes.Collect(ch)
	metricMikrotikTrafficSentBytes.Collect(ch)
}
