package devices

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"sync"
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
	metricRouterMikrotikMemoryUsage          = snitch.NewGauge(boggart.ComponentName+"_device_router_mikrotik_memory_usage_bytes", "Memory usage in Mikrotik router")
	metricRouterMikrotikMemoryAvailable      = snitch.NewGauge(boggart.ComponentName+"_device_router_mikrotik_memory_available_bytes", "Memory available in Mikrotik router")
	metricRouterMikrotikStorageUsage         = snitch.NewGauge(boggart.ComponentName+"_device_router_mikrotik_storage_usage_bytes", "Storage usage in Mikrotik router")
	metricRouterMikrotikStorageAvailable     = snitch.NewGauge(boggart.ComponentName+"_device_router_mikrotik_storage_available_bytes", "Storage available in Mikrotik router")
)

type MikrotikRouter struct {
	boggart.DeviceBase

	mutex        sync.RWMutex
	provider     *mikrotik.Client
	serialNumber string
	interval     time.Duration
}

func NewMikrotikRouter(provider *mikrotik.Client, interval time.Duration) *MikrotikRouter {
	device := &MikrotikRouter{
		provider: provider,
		interval: interval,
	}
	device.Init()
	device.SetDescription("Mikrotik router")

	return device
}

func (d *MikrotikRouter) Types() []boggart.DeviceType {
	return []boggart.DeviceType{
		boggart.DeviceTypeRouter,
	}
}

func (d *MikrotikRouter) Describe(ch chan<- *snitch.Description) {
	serialNumber := d.SerialNumber()
	if serialNumber == "" {
		return
	}

	metricRouterMikrotikTrafficReceivedBytes.With("serial_number", serialNumber).Describe(ch)
	metricRouterMikrotikTrafficSentBytes.With("serial_number", serialNumber).Describe(ch)
	metricRouterMikrotikWifiClients.With("serial_number", serialNumber).Describe(ch)
	metricRouterMikrotikMemoryUsage.With("serial_number", serialNumber).Describe(ch)
	metricRouterMikrotikMemoryAvailable.With("serial_number", serialNumber).Describe(ch)
	metricRouterMikrotikStorageUsage.With("serial_number", serialNumber).Describe(ch)
	metricRouterMikrotikStorageAvailable.With("serial_number", serialNumber).Describe(ch)
}

func (d *MikrotikRouter) Collect(ch chan<- snitch.Metric) {
	serialNumber := d.SerialNumber()
	if serialNumber == "" {
		return
	}

	metricRouterMikrotikTrafficReceivedBytes.With("serial_number", serialNumber).Collect(ch)
	metricRouterMikrotikTrafficSentBytes.With("serial_number", serialNumber).Collect(ch)
	metricRouterMikrotikWifiClients.With("serial_number", serialNumber).Collect(ch)
	metricRouterMikrotikMemoryUsage.With("serial_number", serialNumber).Collect(ch)
	metricRouterMikrotikMemoryAvailable.With("serial_number", serialNumber).Collect(ch)
	metricRouterMikrotikStorageUsage.With("serial_number", serialNumber).Collect(ch)
	metricRouterMikrotikStorageAvailable.With("serial_number", serialNumber).Collect(ch)
}

func (d *MikrotikRouter) Ping(_ context.Context) bool {
	_, err := d.provider.SystemResource()
	return err == nil
}

func (d *MikrotikRouter) SerialNumber() string {
	d.mutex.RLock()
	defer d.mutex.RUnlock()

	return d.serialNumber
}

func (d *MikrotikRouter) Tasks() []workers.Task {
	taskSerialNumber := task.NewFunctionTillStopTask(d.taskSerialNumber)
	taskSerialNumber.SetTimeout(time.Second * 5)
	taskSerialNumber.SetRepeats(-1)
	taskSerialNumber.SetRepeatInterval(time.Minute)
	taskSerialNumber.SetName("device-router-mikrotik-serial-number")

	taskUpdater := task.NewFunctionTask(d.taskUpdater)
	taskUpdater.SetRepeats(-1)
	taskUpdater.SetRepeatInterval(d.interval)
	taskUpdater.SetName("device-router-mikrotik-updater-" + d.serialNumber)

	return []workers.Task{
		taskSerialNumber,
		taskUpdater,
	}
}

func (d *MikrotikRouter) taskSerialNumber(ctx context.Context) (interface{}, error, bool) {
	if !d.IsEnabled() {
		return nil, nil, false
	}

	system, err := d.provider.SystemRouterboard()
	if err != nil {
		return nil, err, false
	}

	serialNumber, ok := system["serial-number"]
	if !ok {
		return nil, errors.New("Serial number not found"), false
	}

	d.mutex.Lock()
	d.serialNumber = serialNumber
	d.mutex.Unlock()

	d.SetDescription("Mikrotik router with serial number " + serialNumber)

	return nil, nil, true
}

func (d *MikrotikRouter) taskUpdater(ctx context.Context) (interface{}, error) {
	if !d.IsEnabled() {
		return nil, nil
	}

	serialNumber := d.SerialNumber()
	if serialNumber == "" {
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

	metricRouterMikrotikWifiClients.With("serial_number", serialNumber).Set(float64(len(clients)))

	for _, client := range clients {
		bytes := strings.Split(client["bytes"], ",")
		if len(bytes) != 2 {
			return nil, err
		}

		name := mikrotik.GetNameByMac(client["mac-address"], arp, dns, dhcp)

		sent, err := strconv.ParseFloat(bytes[0], 64)
		if err != nil {
			return nil, err
		}

		received, err := strconv.ParseFloat(bytes[1], 64)
		if err != nil {
			return nil, err
		}

		metricRouterMikrotikTrafficReceivedBytes.With("serial_number", serialNumber).With(
			"interface", client["interface"],
			"mac", client["mac-address"],
			"name", name).Set(received)
		metricRouterMikrotikTrafficSentBytes.With("serial_number", serialNumber).With(
			"interface", client["interface"],
			"mac", client["mac-address"],
			"name", name).Set(sent)
	}

	// Ports on mikrotik
	stats, err := d.provider.EthernetStats()
	if err != nil {
		return nil, err
	}

	for _, stat := range stats {
		sent, err := strconv.ParseFloat(stat["tx-byte"], 64)
		if err != nil {
			return nil, err
		}

		received, err := strconv.ParseFloat(stat["rx-byte"], 64)
		if err != nil {
			return nil, err
		}

		metricRouterMikrotikTrafficReceivedBytes.With("serial_number", serialNumber).With(
			"interface", stat["name"],
			"mac", stat["mac-address"]).Set(received)
		metricRouterMikrotikTrafficSentBytes.With("serial_number", serialNumber).With(
			"interface", stat["name"],
			"mac", stat["mac-address"]).Set(sent)
	}

	resource, err := d.provider.SystemResource()
	if err != nil {
		return nil, err
	}

	memoryFree, err := strconv.ParseFloat(resource["free-memory"], 64)
	if err != nil {
		return nil, err
	}
	metricRouterMikrotikMemoryAvailable.With("serial_number", serialNumber).Set(memoryFree)

	memoryTotal, err := strconv.ParseFloat(resource["total-memory"], 64)
	if err != nil {
		return nil, err
	}
	metricRouterMikrotikMemoryUsage.With("serial_number", serialNumber).Set(memoryTotal - memoryFree)

	storageFree, err := strconv.ParseFloat(resource["free-hdd-space"], 64)
	if err != nil {
		return nil, err
	}
	metricRouterMikrotikStorageAvailable.With("serial_number", serialNumber).Set(storageFree)

	storageSpace, err := strconv.ParseFloat(resource["total-hdd-space"], 64)
	if err != nil {
		return nil, err
	}
	metricRouterMikrotikStorageUsage.With("serial_number", serialNumber).Set(storageSpace - storageFree)

	return nil, nil
}
