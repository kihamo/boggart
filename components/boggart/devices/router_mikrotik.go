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

	provider     *mikrotik.Client
	syslogClient string
	interval     time.Duration
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

func NewMikrotikRouter(provider *mikrotik.Client, syslogHostname string, interval time.Duration) *MikrotikRouter {
	device := &MikrotikRouter{
		provider:     provider,
		syslogClient: syslogHostname,
		interval:     interval,
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

func (d *MikrotikRouter) Ping(ctx context.Context) bool {
	_, err := d.provider.SystemResource(ctx)
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

func (d *MikrotikRouter) Mac(ctx context.Context, mac string) (*MikrotikRouterMac, error) {
	if !d.IsEnabled() {
		return nil, errors.New("Device is disabled")
	}

	if d.SerialNumber() == "" {
		return nil, errors.New("Serial number is empty")
	}

	info := &MikrotikRouterMac{
		Address: mac,
	}

	if table, err := d.provider.IPARP(ctx); err == nil {
		for _, row := range table {
			if row.MacAddress == mac {
				info.ARP.IP = row.Address
				info.ARP.Comment = row.Comment
				break
			}
		}
	} else {
		return nil, err
	}

	if leases, err := d.provider.IPDHCPServerLease(ctx); err == nil {
		for _, lease := range leases {
			if lease.MacAddress == mac {
				info.DHCP.Hostname = lease.MacAddress
				break
			}
		}
	} else {
		return nil, err
	}

	return info, nil
}

func (d *MikrotikRouter) taskSerialNumber(ctx context.Context) (interface{}, error, bool) {
	if !d.IsEnabled() {
		return nil, nil, false
	}

	system, err := d.provider.SystemRouterboard(ctx)
	if err != nil {
		return nil, err, false
	}

	if system.SerialNumber == "" {
		return nil, errors.New("Serial number is empty"), false
	}

	d.SetSerialNumber(system.SerialNumber)

	// wifi clients
	clients, err := d.provider.InterfaceWirelessRegistrationTable(ctx)
	if err != nil {
		return nil, err, false
	}

	for _, connection := range clients {
		mac, err := d.Mac(ctx, connection.MacAddress)
		if err != nil {
			return nil, err, false
		}

		d.TriggerEvent(ctx, boggart.DeviceEventWifiClientConnected, mac, connection.Interface, system.SerialNumber)
	}

	// vpn clients
	connections, err := d.provider.PPPActive(ctx)
	if err != nil {
		return nil, err, false
	}

	for _, connection := range connections {
		d.TriggerEvent(ctx, boggart.DeviceEventVPNClientConnected, connection.Name, connection.Address, system.SerialNumber)
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

	arp, err := d.provider.IPARP(ctx)
	if err != nil && !mikrotik.IsEmptyResponse(err) {
		return nil, err
	}

	dns, err := d.provider.IPDNSStatic(ctx)
	if err != nil && !mikrotik.IsEmptyResponse(err) {
		return nil, err
	}

	leases, err := d.provider.IPDHCPServerLease(ctx)
	if err != nil && !mikrotik.IsEmptyResponse(err) {
		return nil, err
	}

	// Wifi clients
	clients, err := d.provider.InterfaceWirelessRegistrationTable(ctx)
	if err == nil {
		metricRouterMikrotikWifiClients.With("serial_number", serialNumber).Set(float64(len(clients)))

		for _, client := range clients {
			bytes := strings.Split(client.Bytes, ",")
			if len(bytes) != 2 {
				return nil, err
			}

			name := mikrotik.GetNameByMac(client.MacAddress, arp, dns, leases)

			sent, err := strconv.ParseFloat(bytes[0], 64)
			if err != nil {
				return nil, err
			}

			received, err := strconv.ParseFloat(bytes[1], 64)
			if err == nil {
				metricRouterMikrotikTrafficReceivedBytes.With("serial_number", serialNumber).With(
					"interface", client.Interface,
					"mac", client.MacAddress,
					"name", name).Set(received)
				metricRouterMikrotikTrafficSentBytes.With("serial_number", serialNumber).With(
					"interface", client.Interface,
					"mac", client.MacAddress,
					"name", name).Set(sent)
			} else if !mikrotik.IsEmptyResponse(err) {
				return nil, err
			}
		}
	} else if !mikrotik.IsEmptyResponse(err) {
		return nil, err
	}

	// Ports on mikrotik
	stats, err := d.provider.InterfaceStats(ctx)
	if err == nil {
		for _, stat := range stats {
			metricRouterMikrotikTrafficReceivedBytes.With("serial_number", serialNumber).With(
				"interface", stat.Name,
				"mac", stat.MacAddress).Set(float64(stat.RXByte))
			metricRouterMikrotikTrafficSentBytes.With("serial_number", serialNumber).With(
				"interface", stat.Name,
				"mac", stat.MacAddress).Set(float64(stat.TXByte))
		}
	} else if !mikrotik.IsEmptyResponse(err) {
		return nil, err
	}

	resource, err := d.provider.SystemResource(ctx)
	if err == nil {
		metricRouterMikrotikCPULoad.With("serial_number", serialNumber).Set(float64(resource.CPULoad))
		metricRouterMikrotikMemoryAvailable.With("serial_number", serialNumber).Set(float64(resource.FreeMemory))
		metricRouterMikrotikMemoryUsage.With("serial_number", serialNumber).Set(float64(resource.TotalMemory - resource.FreeMemory))
		metricRouterMikrotikStorageAvailable.With("serial_number", serialNumber).Set(float64(resource.FreeHDDSpace))
		metricRouterMikrotikStorageUsage.With("serial_number", serialNumber).Set(float64(resource.TotalHDDSpace - resource.FreeHDDSpace))
	} else if !mikrotik.IsEmptyResponse(err) {
		return nil, err
	}

	disks, err := d.provider.SystemDisk(ctx)
	if err == nil {
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
	} else if !mikrotik.IsEmptyResponse(err) {
		return nil, err
	}

	health, err := d.provider.SystemHealth(ctx)
	if err == nil {
		metricRouterMikrotikVoltage.With("serial_number", serialNumber).Set(health.Voltage)
		metricRouterMikrotikTemperature.With("serial_number", serialNumber).Set(float64(health.Temperature))
	} else if !mikrotik.IsEmptyResponse(err) {
		return nil, err
	}

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

func (l *MikrotikRouterListener) Run(ctx context.Context, event workers.Event, t time.Time, args ...interface{}) {
	switch event {
	case boggart.DeviceEventSyslogReceive:
		message := args[0].(map[string]interface{})

		client, ok := message["client"]
		if !ok || client != l.router.syslogClient {
			return
		}

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
			check := wifiClientRegexp.FindStringSubmatch(content.(string))
			if len(check) < 4 {
				return
			}

			if _, err := net.ParseMAC(check[1]); err != nil {
				return
			}

			mac, err := l.router.Mac(ctx, check[1])
			if err != nil {
				return
			}

			switch check[3] {
			case "connected":
				l.router.TriggerEvent(ctx, boggart.DeviceEventWifiClientConnected, mac, check[2], l.router.SerialNumber())
			case "disconnected":
				l.router.TriggerEvent(ctx, boggart.DeviceEventWifiClientDisconnected, mac, check[2], l.router.SerialNumber())
			}

		case "vpn":
			check := vpnClientRegexp.FindStringSubmatch(content.(string))
			if len(check) < 2 {
				return
			}

			switch check[2] {
			case "in":
				l.router.TriggerEvent(ctx, boggart.DeviceEventVPNClientConnected, check[1], check[3], l.router.SerialNumber())
			case "out":
				l.router.TriggerEvent(ctx, boggart.DeviceEventVPNClientDisconnected, check[1], l.router.SerialNumber())
			}
		}
	}
}
