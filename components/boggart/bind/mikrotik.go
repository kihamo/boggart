package bind

import (
	"context"
	"errors"
	"net"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/providers/mikrotik"
	"github.com/kihamo/boggart/components/mqtt"
	"github.com/kihamo/go-workers"
	"github.com/kihamo/go-workers/listener"
	"github.com/kihamo/go-workers/task"
	"github.com/kihamo/snitch"
)

const (
	MikrotikMQTTTopicWiFiMACState         mqtt.Topic = boggart.ComponentName + "/router/+/wifi/clients/+/state"
	MikrotikMQTTTopicWiFiConnectedMAC     mqtt.Topic = boggart.ComponentName + "/router/+/wifi/clients/last/on/mac"
	MikrotikMQTTTopicWiFiDisconnectedMAC  mqtt.Topic = boggart.ComponentName + "/router/+/wifi/clients/last/on/mac"
	MikrotikMQTTTopicVPNLoginState        mqtt.Topic = boggart.ComponentName + "/router/+/vpn/clients/+/state"
	MikrotikMQTTTopicVPNConnectedLogin    mqtt.Topic = boggart.ComponentName + "/router/+/vpn/clients/last/on/login"
	MikrotikMQTTTopicVPNDisconnectedLogin mqtt.Topic = boggart.ComponentName + "/router/+/vpn/clients/last/off/login"
)

var (
	metricMikrotikTrafficReceivedBytes = snitch.NewGauge(boggart.ComponentName+"_device_router_mikrotik_traffic_received_bytes", "Mikrotik traffic received bytes")
	metricMikrotikTrafficSentBytes     = snitch.NewGauge(boggart.ComponentName+"_device_router_mikrotik_traffic_sent_bytes", "Mikrotik traffic sent bytes")
	metricMikrotikWifiClients          = snitch.NewGauge(boggart.ComponentName+"_device_router_mikrotik_wifi_clients_total", "Mikrotik wifi clients online")
	metricMikrotikCPULoad              = snitch.NewGauge(boggart.ComponentName+"_device_router_mikrotik_cpu_load_percent", "CPU load in percents")
	metricMikrotikMemoryUsage          = snitch.NewGauge(boggart.ComponentName+"_device_router_mikrotik_memory_usage_bytes", "Memory usage in Mikrotik router")
	metricMikrotikMemoryAvailable      = snitch.NewGauge(boggart.ComponentName+"_device_router_mikrotik_memory_available_bytes", "Memory available in Mikrotik router")
	metricMikrotikStorageUsage         = snitch.NewGauge(boggart.ComponentName+"_device_router_mikrotik_storage_usage_bytes", "Storage usage in Mikrotik router")
	metricMikrotikStorageAvailable     = snitch.NewGauge(boggart.ComponentName+"_device_router_mikrotik_storage_available_bytes", "Storage available in Mikrotik router")
	metricMikrotikDiskUsage            = snitch.NewGauge(boggart.ComponentName+"_device_router_mikrotik_disk_usage_bytes", "Disk usage in Mikrotik router")
	metricMikrotikDiskAvailable        = snitch.NewGauge(boggart.ComponentName+"_device_router_mikrotik_disk_available_bytes", "Disk available in Mikrotik router")
	metricMikrotikVoltage              = snitch.NewGauge(boggart.ComponentName+"_device_router_mikrotik_voltage_volt", "Voltage")
	metricMikrotikTemperature          = snitch.NewGauge(boggart.ComponentName+"_device_router_mikrotik_temperature_celsius", "Temperature")

	wifiClientRegexp = regexp.MustCompile(`^([^@]+)@([^:\s]+):\s+([^\s,]+)`)
	vpnClientRegexp  = regexp.MustCompile(`^(\S+) logged (in|out), (.+?)$`)
)

type Mikrotik struct {
	boggart.DeviceBindBase
	boggart.DeviceBindMQTT

	provider     *mikrotik.Client
	host         string
	syslogClient string
}

type MikrotikListener struct {
	listener.BaseListener

	router *Mikrotik
}

type MikrotikMac struct {
	Address string
	ARP     struct {
		IP      string
		Comment string
	}
	DHCP struct {
		Hostname string
	}
}

type MikrotikConfig struct {
	Address      string `valid:"url,required"`
	SyslogClient string `valid:"host"`
}

func (d Mikrotik) Config() interface{} {
	return &MikrotikConfig{}
}

func (d Mikrotik) CreateBind(c interface{}) (boggart.DeviceBind, error) {
	config := c.(*MikrotikConfig)

	u, _ := url.Parse(config.Address)
	username := u.User.Username()
	password, _ := u.User.Password()

	device := &Mikrotik{
		provider:     mikrotik.NewClient(u.Host, username, password, time.Second*10),
		host:         u.Host + "-" + u.Port(),
		syslogClient: config.SyslogClient,
	}

	if device.syslogClient == "" {
		device.syslogClient = u.Hostname() + ":514"
	}

	device.Init()

	return device, nil
}

func NewMikrotikRouterListener(router *Mikrotik) *MikrotikListener {
	l := &MikrotikListener{
		router: router,
	}
	l.Init()

	return l
}

func (d *Mikrotik) Describe(ch chan<- *snitch.Description) {
	serialNumber := d.SerialNumber()
	if serialNumber == "" {
		return
	}

	metricMikrotikTrafficReceivedBytes.With("serial_number", serialNumber).Describe(ch)
	metricMikrotikTrafficSentBytes.With("serial_number", serialNumber).Describe(ch)
	metricMikrotikWifiClients.With("serial_number", serialNumber).Describe(ch)
	metricMikrotikCPULoad.With("serial_number", serialNumber).Describe(ch)
	metricMikrotikMemoryUsage.With("serial_number", serialNumber).Describe(ch)
	metricMikrotikMemoryAvailable.With("serial_number", serialNumber).Describe(ch)
	metricMikrotikStorageUsage.With("serial_number", serialNumber).Describe(ch)
	metricMikrotikStorageAvailable.With("serial_number", serialNumber).Describe(ch)
	metricMikrotikDiskUsage.With("serial_number", serialNumber).Describe(ch)
	metricMikrotikDiskAvailable.With("serial_number", serialNumber).Describe(ch)
	metricMikrotikVoltage.With("serial_number", serialNumber).Describe(ch)
	metricMikrotikTemperature.With("serial_number", serialNumber).Describe(ch)
}

func (d *Mikrotik) Collect(ch chan<- snitch.Metric) {
	serialNumber := d.SerialNumber()
	if serialNumber == "" {
		return
	}

	metricMikrotikTrafficReceivedBytes.With("serial_number", serialNumber).Collect(ch)
	metricMikrotikTrafficSentBytes.With("serial_number", serialNumber).Collect(ch)
	metricMikrotikWifiClients.With("serial_number", serialNumber).Collect(ch)
	metricMikrotikCPULoad.With("serial_number", serialNumber).Collect(ch)
	metricMikrotikMemoryUsage.With("serial_number", serialNumber).Collect(ch)
	metricMikrotikMemoryAvailable.With("serial_number", serialNumber).Collect(ch)
	metricMikrotikStorageUsage.With("serial_number", serialNumber).Collect(ch)
	metricMikrotikStorageAvailable.With("serial_number", serialNumber).Collect(ch)
	metricMikrotikDiskUsage.With("serial_number", serialNumber).Collect(ch)
	metricMikrotikDiskAvailable.With("serial_number", serialNumber).Collect(ch)
	metricMikrotikVoltage.With("serial_number", serialNumber).Collect(ch)
	metricMikrotikTemperature.With("serial_number", serialNumber).Collect(ch)
}

func (d *Mikrotik) Tasks() []workers.Task {
	taskLiveness := task.NewFunctionTask(d.taskLiveness)
	taskLiveness.SetTimeout(time.Second * 5)
	taskLiveness.SetRepeats(-1)
	taskLiveness.SetRepeatInterval(time.Minute)
	taskLiveness.SetName("bind-mikrotik-liveness")

	taskStateUpdater := task.NewFunctionTask(d.taskStateUpdater)
	taskStateUpdater.SetRepeats(-1)
	taskStateUpdater.SetRepeatInterval(time.Minute * 5)
	taskStateUpdater.SetName("bind-mikrotik-updater-" + d.host)

	return []workers.Task{
		taskLiveness,
		taskStateUpdater,
	}
}

func (d *Mikrotik) Listeners() []workers.ListenerWithEvents {
	return []workers.ListenerWithEvents{
		NewMikrotikRouterListener(d),
	}
}

func (d *Mikrotik) Mac(ctx context.Context, mac string) (*MikrotikMac, error) {
	if d.SerialNumber() == "" {
		return nil, errors.New("serial number is empty")
	}

	info := &MikrotikMac{
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

func (d *Mikrotik) taskLiveness(ctx context.Context) (interface{}, error) {
	system, err := d.provider.SystemRouterboard(ctx)
	if err != nil {
		d.UpdateStatus(boggart.DeviceStatusOffline)
		return nil, err
	}

	if system.SerialNumber == "" {
		d.UpdateStatus(boggart.DeviceStatusOffline)
		return nil, errors.New("serial number is empty")
	}

	d.UpdateStatus(boggart.DeviceStatusOnline)
	if d.SerialNumber() != "" {
		return nil, nil
	}

	d.SetSerialNumber(system.SerialNumber)
	sn := system.SerialNumber

	// wifi clients
	clients, err := d.provider.InterfaceWirelessRegistrationTable(ctx)
	if err != nil {
		return nil, err
	}

	for _, connection := range clients {
		mac, err := d.Mac(ctx, connection.MacAddress)
		if err != nil {
			return nil, err
		}

		login := mqtt.NameReplace(mac.Address)

		d.MQTTPublishAsync(ctx, MikrotikMQTTTopicWiFiConnectedMAC.Format(sn), 0, false, login)
		d.MQTTPublishAsync(ctx, MikrotikMQTTTopicWiFiMACState.Format(sn, login), 0, false, true)
	}

	// vpn clients
	connections, err := d.provider.PPPActive(ctx)
	if err != nil {
		return nil, err
	}

	for _, connection := range connections {
		login := mqtt.NameReplace(connection.Name)

		d.MQTTPublishAsync(ctx, MikrotikMQTTTopicVPNConnectedLogin.Format(sn), 0, false, login)
		d.MQTTPublishAsync(ctx, MikrotikMQTTTopicVPNLoginState.Format(sn, login), 0, false, true)
	}

	return nil, nil
}

func (d *Mikrotik) taskStateUpdater(ctx context.Context) (interface{}, error) {
	if d.Status() != boggart.DeviceStatusOnline {
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
		metricMikrotikWifiClients.With("serial_number", serialNumber).Set(float64(len(clients)))

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
				metricMikrotikTrafficReceivedBytes.With("serial_number", serialNumber).With(
					"interface", client.Interface,
					"mac", client.MacAddress,
					"name", name).Set(received)
				metricMikrotikTrafficSentBytes.With("serial_number", serialNumber).With(
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
			metricMikrotikTrafficReceivedBytes.With("serial_number", serialNumber).With(
				"interface", stat.Name,
				"mac", stat.MacAddress).Set(float64(stat.RXByte))
			metricMikrotikTrafficSentBytes.With("serial_number", serialNumber).With(
				"interface", stat.Name,
				"mac", stat.MacAddress).Set(float64(stat.TXByte))
		}
	} else if !mikrotik.IsEmptyResponse(err) {
		return nil, err
	}

	resource, err := d.provider.SystemResource(ctx)
	if err == nil {
		metricMikrotikCPULoad.With("serial_number", serialNumber).Set(float64(resource.CPULoad))
		metricMikrotikMemoryAvailable.With("serial_number", serialNumber).Set(float64(resource.FreeMemory))
		metricMikrotikMemoryUsage.With("serial_number", serialNumber).Set(float64(resource.TotalMemory - resource.FreeMemory))
		metricMikrotikStorageAvailable.With("serial_number", serialNumber).Set(float64(resource.FreeHDDSpace))
		metricMikrotikStorageUsage.With("serial_number", serialNumber).Set(float64(resource.TotalHDDSpace - resource.FreeHDDSpace))
	} else if !mikrotik.IsEmptyResponse(err) {
		return nil, err
	}

	disks, err := d.provider.SystemDisk(ctx)
	if err == nil {
		for _, disk := range disks {
			metricMikrotikDiskUsage.With("serial_number", serialNumber).With(
				"name", disk.Name,
				"label", disk.Label,
			).Set(float64(disk.Size - disk.Free))
			metricMikrotikDiskAvailable.With("serial_number", serialNumber).With(
				"name", disk.Name,
				"label", disk.Label,
			).Set(float64(disk.Free))
		}
	} else if !mikrotik.IsEmptyResponse(err) {
		return nil, err
	}

	health, err := d.provider.SystemHealth(ctx)
	if err == nil {
		metricMikrotikVoltage.With("serial_number", serialNumber).Set(health.Voltage)
		metricMikrotikTemperature.With("serial_number", serialNumber).Set(float64(health.Temperature))
	} else if !mikrotik.IsEmptyResponse(err) {
		return nil, err
	}

	return nil, nil
}

func (d *Mikrotik) MQTTTopics() []mqtt.Topic {
	return []mqtt.Topic{
		MikrotikMQTTTopicWiFiMACState,
		MikrotikMQTTTopicWiFiConnectedMAC,
		MikrotikMQTTTopicWiFiDisconnectedMAC,
		MikrotikMQTTTopicVPNLoginState,
		MikrotikMQTTTopicVPNConnectedLogin,
		MikrotikMQTTTopicVPNDisconnectedLogin,
	}
}

func (l *MikrotikListener) Events() []workers.Event {
	return []workers.Event{
		boggart.DeviceEventSyslogReceive,
	}
}

func (l *MikrotikListener) Name() string {
	return boggart.ComponentName + ".device.router.mikrotik." + l.router.SerialNumber()
}

func (l *MikrotikListener) Run(ctx context.Context, event workers.Event, t time.Time, args ...interface{}) {
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

			sn := l.router.SerialNumber()
			login := mqtt.NameReplace(mac.Address)

			switch check[3] {
			case "connected":
				l.router.MQTTPublishAsync(ctx, MikrotikMQTTTopicWiFiConnectedMAC.Format(sn), 0, false, login)
				l.router.MQTTPublishAsync(ctx, MikrotikMQTTTopicWiFiMACState.Format(sn, login), 0, false, []byte(`1`))
			case "disconnected":
				l.router.MQTTPublishAsync(ctx, MikrotikMQTTTopicWiFiDisconnectedMAC.Format(sn), 0, false, login)
				l.router.MQTTPublishAsync(ctx, MikrotikMQTTTopicWiFiMACState.Format(sn, login), 0, false, []byte(`0`))
			}

		case "vpn":
			check := vpnClientRegexp.FindStringSubmatch(content.(string))
			if len(check) < 2 {
				return
			}

			sn := l.router.SerialNumber()
			login := mqtt.NameReplace(check[1])

			switch check[2] {
			case "in":
				l.router.MQTTPublishAsync(ctx, MikrotikMQTTTopicVPNConnectedLogin.Format(sn), 0, false, login)
				l.router.MQTTPublishAsync(ctx, MikrotikMQTTTopicVPNLoginState.Format(sn, login), 0, false, []byte(`1`))
			case "out":
				l.router.MQTTPublishAsync(ctx, MikrotikMQTTTopicVPNDisconnectedLogin.Format(sn), 0, false, login)
				l.router.MQTTPublishAsync(ctx, MikrotikMQTTTopicVPNLoginState.Format(sn, login), 0, false, []byte(`0`))
			}
		}
	}
}
