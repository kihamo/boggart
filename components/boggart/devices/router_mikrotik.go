package devices

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/providers/mikrotik"
	"github.com/kihamo/go-workers"
	"github.com/kihamo/go-workers/task"
	"github.com/kihamo/snitch"
)

var (
	metricRouterMikrotikTrafficReceivedBytes = snitch.NewGauge(boggart.ComponentName+"_device_router_mikrotik_traffic_received_bytes", "Mikrotik traffic received bytes")
	metricRouterMikrotikTrafficSentBytes     = snitch.NewGauge(boggart.ComponentName+"_device_router_mikrotik_traffic_sent_bytes", "Mikrotik traffic sent bytes")
	metricRouterMikrotikWifiClients          = snitch.NewGauge(boggart.ComponentName+"_device_router_mikrotik_wifi_clients_total", "Mikrotik wifi clients online")
)

type MikrotikRouter struct {
	boggart.DeviceBase

	provider     *mikrotik.Client
	serialNumber string
	interval     time.Duration
}

func NewMikrotikRouter(provider *mikrotik.Client, interval time.Duration) (*MikrotikRouter, error) {
	device := &MikrotikRouter{
		provider: provider,
		interval: interval,
	}
	device.Init()

	system, err := provider.SystemRouterboard()
	if err != nil {
		return nil, err
	}

	var ok bool
	device.serialNumber, ok = system["serial-number"]
	if !ok {
		return nil, errors.New("Serial number not found")
	}

	device.SetDescription("Mikrotik router with serial number " + device.serialNumber)

	return device, nil
}

func (d *MikrotikRouter) Types() []boggart.DeviceType {
	return []boggart.DeviceType{
		boggart.DeviceTypeRouter,
	}
}

func (d *MikrotikRouter) Describe(ch chan<- *snitch.Description) {
	metricRouterMikrotikTrafficReceivedBytes.With("serial_number", d.serialNumber).Describe(ch)
	metricRouterMikrotikTrafficSentBytes.With("serial_number", d.serialNumber).Describe(ch)
	metricRouterMikrotikWifiClients.With("serial_number", d.serialNumber).Describe(ch)
}

func (d *MikrotikRouter) Collect(ch chan<- snitch.Metric) {
	metricRouterMikrotikTrafficReceivedBytes.With("serial_number", d.serialNumber).Collect(ch)
	metricRouterMikrotikTrafficSentBytes.With("serial_number", d.serialNumber).Collect(ch)
	metricRouterMikrotikWifiClients.With("serial_number", d.serialNumber).Collect(ch)
}

func (d *MikrotikRouter) Tasks() []workers.Task {
	taskUpdater := task.NewFunctionTask(d.updater)
	taskUpdater.SetRepeats(-1)
	taskUpdater.SetRepeatInterval(d.interval)
	taskUpdater.SetName("device-router-mikrotik-updater-" + d.serialNumber)

	return []workers.Task{
		taskUpdater,
	}
}

func (d *MikrotikRouter) updater(ctx context.Context) (interface{}, error) {
	if !d.IsEnabled() {
		return nil, nil
	}

	arp, err := d.provider.ARPTable()
	if err != nil {
		return nil, err
	}

	dns, err := d.provider.DNSStatic()
	if err != nil {
		return nil, err
	}

	dhcp, err := d.provider.DHCPLease()
	if err != nil {
		return nil, err
	}

	// Wifi clients
	clients, err := d.provider.WifiClients()
	if err != nil {
		return nil, err
	}

	metricRouterMikrotikWifiClients.With("serial_number", d.serialNumber).Set(float64(len(clients)))

	for _, client := range clients {
		bytes := strings.Split(client["bytes"], ",")
		if len(bytes) != 2 {
			fmt.Println("mikrotik", err)

			return nil, err
		}

		name := mikrotik.GetNameByMac(client["mac-address"], arp, dns, dhcp)

		sent, err := strconv.ParseFloat(bytes[0], 64)
		if err != nil {
			fmt.Println("mikrotik", err)
			return nil, err
		}

		received, err := strconv.ParseFloat(bytes[1], 64)
		if err != nil {
			fmt.Println("mikrotik", err)
			return nil, err
		}

		metricRouterMikrotikTrafficReceivedBytes.With("serial_number", d.serialNumber).With(
			"interface", client["interface"],
			"mac", client["mac-address"],
			"name", name).Set(received)
		metricRouterMikrotikTrafficSentBytes.With("serial_number", d.serialNumber).With(
			"interface", client["interface"],
			"mac", client["mac-address"],
			"name", name).Set(sent)
	}

	// Ports on mikrotik
	stats, err := d.provider.EthernetStats()
	if err != nil {
		fmt.Println("mikrotik", err)
		return nil, err
	}

	for _, stat := range stats {
		sent, err := strconv.ParseFloat(stat["tx-byte"], 64)
		if err != nil {
			fmt.Println("mikrotik", err)
			return nil, err
		}

		received, err := strconv.ParseFloat(stat["rx-byte"], 64)
		if err != nil {
			fmt.Println("mikrotik", err)
			return nil, err
		}

		metricRouterMikrotikTrafficReceivedBytes.With("serial_number", d.serialNumber).With(
			"interface", stat["name"],
			"mac", stat["mac-address"]).Set(received)
		metricRouterMikrotikTrafficSentBytes.With("serial_number", d.serialNumber).With(
			"interface", stat["name"],
			"mac", stat["mac-address"]).Set(sent)
	}

	return nil, nil
}
