package internal

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/providers/mikrotik"
	"github.com/kihamo/snitch"
)

const (
	MetricMikrotikWifiClients             = boggart.ComponentName + "_mikrotik_wifi_clients_total"
	MetricMikrotikWifiClientReceivedBytes = boggart.ComponentName + "_mikrotik_wifi_client_received_bytes"
	MetricMikrotikWifiClientSentBytes     = boggart.ComponentName + "_mikrotik_wifi_client_sent_bytes"
)

var (
	metricMikrotikWifiClients             = snitch.NewGauge(MetricMikrotikWifiClients, "Mikrotik wifi clients online")
	metricMikrotikWifiClientReceivedBytes = snitch.NewGauge(MetricMikrotikWifiClientReceivedBytes, "Mikrotik wifi client reveived bytes")
	metricMikrotikWifiClientSentBytes     = snitch.NewGauge(MetricMikrotikWifiClientSentBytes, "Mikrotik wifi client sent bytes")
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

	for _, client := range clients {
		bytes := strings.Split(client["bytes"], ",")
		if len(bytes) != 2 {
			return fmt.Errorf("Wrong value of bytes field %s", client["bytes"])
		}

		sent, err := strconv.ParseFloat(bytes[0], 64)
		if err != nil {
			return fmt.Errorf("Failed convert sent bytes to float64 %s", bytes[0])
		}

		received, err := strconv.ParseFloat(bytes[1], 64)
		if err != nil {
			return fmt.Errorf("Failed convert received bytes to float64 %s", bytes[1])
		}

		metricMikrotikWifiClientReceivedBytes.With("mac", client["mac-address"]).Set(received)
		metricMikrotikWifiClientSentBytes.With("mac", client["mac-address"]).Set(sent)
	}

	return nil
}

func (c *MetricsCollector) DescribeMikrotik(ch chan<- *snitch.Description) {
	metricMikrotikWifiClients.Describe(ch)
	metricMikrotikWifiClientReceivedBytes.Describe(ch)
	metricMikrotikWifiClientSentBytes.Describe(ch)
}

func (c *MetricsCollector) CollectMikrotik(ch chan<- snitch.Metric) {
	metricMikrotikWifiClients.Collect(ch)
	metricMikrotikWifiClientReceivedBytes.Collect(ch)
	metricMikrotikWifiClientSentBytes.Collect(ch)
}
