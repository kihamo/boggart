package devices

import (
	"context"
	"errors"
	"net"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/providers/mikrotik"
	"github.com/kihamo/go-workers"
	"github.com/kihamo/go-workers/listener"
	"github.com/kihamo/go-workers/task"
	"github.com/kihamo/snitch"
)

var (
	metricRouterMikrotikTrafficReceivedBytes = snitch.NewGauge(boggart.ComponentName+"_device_router_mikrotik_traffic_received_bytes", "Mikrotik traffic received bytes")
	metricRouterMikrotikTrafficSentBytes     = snitch.NewGauge(boggart.ComponentName+"_device_router_mikrotik_traffic_sent_bytes", "Mikrotik traffic sent bytes")
	metricRouterMikrotikWifiClients          = snitch.NewGauge(boggart.ComponentName+"_device_router_mikrotik_wifi_clients_total", "Mikrotik wifi clients online")
	metricRouterMikrotikCPULoad              = snitch.NewGauge(boggart.ComponentName+"_device_router_mikrotik_cpu_load_percent", "CPU load in percents")
	metricRouterMikrotikMemoryUsage          = snitch.NewGauge(boggart.ComponentName+"_device_router_mikrotik_memory_usage_bytes", "Memory usage in Mikrotik router")
	metricRouterMikrotikMemoryAvailable      = snitch.NewGauge(boggart.ComponentName+"_device_router_mikrotik_memory_available_bytes", "Memory available in Mikrotik router")
	metricRouterMikrotikStorageUsage         = snitch.NewGauge(boggart.ComponentName+"_device_router_mikrotik_storage_usage_bytes", "Storage usage in Mikrotik router")
	metricRouterMikrotikStorageAvailable     = snitch.NewGauge(boggart.ComponentName+"_device_router_mikrotik_storage_available_bytes", "Storage available in Mikrotik router")
	metricRouterMikrotikDiskUsage            = snitch.NewGauge(boggart.ComponentName+"_device_router_mikrotik_disk_usage_bytes", "Disk usage in Mikrotik router")
	metricRouterMikrotikDiskAvailable        = snitch.NewGauge(boggart.ComponentName+"_device_router_mikrotik_disk_available_bytes", "Disk available in Mikrotik router")
	metricRouterMikrotikVoltage              = snitch.NewGauge(boggart.ComponentName+"_device_router_mikrotik_voltage_volt", "Voltage")
	metricRouterMikrotikTemperature          = snitch.NewGauge(boggart.ComponentName+"_device_router_mikrotik_temperature_celsius", "Temperature")

	wifiClientRegexp = regexp.MustCompile(`^([^@]+)@([^:\s]+):\s+([^\s,]+)`)
	vpnClientRegexp  = regexp.MustCompile(`^(\S+) logged (in|out), (.+?)$`)
)

type MikrotikRouter struct {
	boggart.DeviceBase
	boggart.DeviceSerialNumber

	provider *mikrotik.Client
	interval time.Duration
}

type MikrotikRouterListener struct {
	listener.BaseListener

	router *MikrotikRouter
}

type MikrotikRouterMac struct {
	Address string
	ARP     struct {
		IP      string
		Comment string
	}
	DHCP struct {
		Hostname string
	}
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

func NewMikrotikRouterListener(router *MikrotikRouter) *MikrotikRouterListener {
	l := &MikrotikRouterListener{
		router: router,
	}
	l.Init()

	return l
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
	metricRouterMikrotikCPULoad.With("serial_number", serialNumber).Describe(ch)
	metricRouterMikrotikMemoryUsage.With("serial_number", serialNumber).Describe(ch)
	metricRouterMikrotikMemoryAvailable.With("serial_number", serialNumber).Describe(ch)
	metricRouterMikrotikStorageUsage.With("serial_number", serialNumber).Describe(ch)
	metricRouterMikrotikStorageAvailable.With("serial_number", serialNumber).Describe(ch)
	metricRouterMikrotikDiskUsage.With("serial_number", serialNumber).Describe(ch)
	metricRouterMikrotikDiskAvailable.With("serial_number", serialNumber).Describe(ch)
	metricRouterMikrotikVoltage.With("serial_number", serialNumber).Describe(ch)
	metricRouterMikrotikTemperature.With("serial_number", serialNumber).Describe(ch)
}

func (d *MikrotikRouter) Collect(ch chan<- snitch.Metric) {
	serialNumber := d.SerialNumber()
	if serialNumber == "" {
		return
	}

	metricRouterMikrotikTrafficReceivedBytes.With("serial_number", serialNumber).Collect(ch)
	metricRouterMikrotikTrafficSentBytes.With("serial_number", serialNumber).Collect(ch)
	metricRouterMikrotikWifiClients.With("serial_number", serialNumber).Collect(ch)
	metricRouterMikrotikCPULoad.With("serial_number", serialNumber).Collect(ch)
	metricRouterMikrotikMemoryUsage.With("serial_number", serialNumber).Collect(ch)
	metricRouterMikrotikMemoryAvailable.With("serial_number", serialNumber).Collect(ch)
	metricRouterMikrotikStorageUsage.With("serial_number", serialNumber).Collect(ch)
	metricRouterMikrotikStorageAvailable.With("serial_number", serialNumber).Collect(ch)
	metricRouterMikrotikDiskUsage.With("serial_number", serialNumber).Collect(ch)
	metricRouterMikrotikDiskAvailable.With("serial_number", serialNumber).Collect(ch)
	metricRouterMikrotikVoltage.With("serial_number", serialNumber).Collect(ch)
	metricRouterMikrotikTemperature.With("serial_number", serialNumber).Collect(ch)
}

func (d *MikrotikRouter) Ping(_ context.Context) bool {
	_, err := d.provider.SystemResource()
	return err == nil
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
	taskUpdater.SetName("device-router-mikrotik-updater-" + d.Id())

	return []workers.Task{
		taskSerialNumber,
		taskUpdater,
	}
}

func (d *MikrotikRouter) Listeners() []workers.ListenerWithEvents {
	return []workers.ListenerWithEvents{
		NewMikrotikRouterListener(d),
	}
}

func (d *MikrotikRouter) Mac(mac string) *MikrotikRouterMac {
	info := &MikrotikRouterMac{
		Address: mac,
	}

	if !d.IsEnabled() {
		return info
	}

	if d.SerialNumber() == "" {
		return info
	}

	if arp, err := d.provider.ARPTable(); err == nil {
		if record, ok := arp[mac]; ok {
			info.ARP.IP = record["address"]
			info.ARP.Comment = record["comment"]
		}
	}

	if dhcp, err := d.provider.DHCPLease(); err == nil {
		if record, ok := dhcp[mac]; ok {
			info.DHCP.Hostname = record
		}
	}

	return info
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

	d.SetSerialNumber(serialNumber)
	d.SetDescription("Mikrotik router with serial number " + serialNumber)

	// wifi clients
	clients, err := d.provider.WifiClients()
	if err != nil {
		return nil, err, false
	}

	for _, connection := range clients {
		d.TriggerEvent(boggart.DeviceEventWifiClientConnected, d.Mac(connection["mac-address"]), connection["interface"])
	}

	// vpn clients
	connections, err := d.provider.PPPActiveConnections()
	if err != nil {
		return nil, err, false
	}

	for _, connection := range connections {
		d.TriggerEvent(boggart.DeviceEventVPNClientConnected, connection["name"], connection["address"])
	}

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

	cpuLoad, err := strconv.ParseFloat(resource["cpu-load"], 64)
	if err != nil {
		return nil, err
	}
	metricRouterMikrotikCPULoad.With("serial_number", serialNumber).Set(cpuLoad)

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

	disks, err := d.provider.SystemDisk(ctx)
	if err != nil {
		return nil, err
	}

	for _, disk := range disks {
		metricRouterMikrotikDiskUsage.With("serial_number", serialNumber).With(
			"name", disk.Name,
			"label", disk.Label,
		).Set(float64(disk.Size - disk.Free))
		metricRouterMikrotikDiskAvailable.With("serial_number", serialNumber).With(
			"name", disk.Name,
			"label", disk.Label,
		).Set(float64(disk.Free))
	}

	health, err := d.provider.SystemHealth()
	if err != nil {
		return nil, err
	}

	metricRouterMikrotikVoltage.With("serial_number", serialNumber).Set(health.Voltage)
	metricRouterMikrotikTemperature.With("serial_number", serialNumber).Set(float64(health.Temperature))

	return nil, nil
}

func (l *MikrotikRouterListener) Events() []workers.Event {
	return []workers.Event{
		boggart.DeviceEventSyslogReceive,
	}
}

func (l *MikrotikRouterListener) Name() string {
	return boggart.ComponentName + ".device.router.mikrotik." + l.router.SerialNumber()
}

func (l *MikrotikRouterListener) Run(_ context.Context, event workers.Event, t time.Time, args ...interface{}) {
	switch event {
	case boggart.DeviceEventSyslogReceive:
		message := args[0].(map[string]interface{})

		tag, ok := message["tag"]
		if !ok {
			return
		}

		content, ok := message["content"]
		if !ok {
			return
		}

		switch tag {
		case "wifi":
			// TODO: check hostname

			check := wifiClientRegexp.FindStringSubmatch(content.(string))
			if len(check) < 4 {
				return
			}

			if _, err := net.ParseMAC(check[1]); err != nil {
				return
			}

			mac := l.router.Mac(check[1])

			switch check[3] {
			case "connected":
				l.router.TriggerEvent(boggart.DeviceEventWifiClientConnected, mac, check[2])
			case "disconnected":
				l.router.TriggerEvent(boggart.DeviceEventWifiClientDisconnected, mac, check[2])
			}

		case "vpn":
			check := vpnClientRegexp.FindStringSubmatch(content.(string))
			if len(check) < 2 {
				return
			}

			switch check[2] {
			case "in":
				l.router.TriggerEvent(boggart.DeviceEventVPNClientConnected, check[1], check[3])
			case "out":
				l.router.TriggerEvent(boggart.DeviceEventVPNClientDisconnected, check[1])
			}
		}
	}
}
