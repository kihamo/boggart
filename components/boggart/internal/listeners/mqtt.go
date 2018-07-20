package listeners

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"time"

	m "github.com/eclipse/paho.mqtt.golang"
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/devices"
	"github.com/kihamo/go-workers"
	"github.com/kihamo/go-workers/listener"
)

type MQTTListener struct {
	listener.BaseListener

	mqtt m.Client
}

func NewMQTTListener(servers []*url.URL, username, password string) *MQTTListener {
	l := &MQTTListener{}

	opts := &m.ClientOptions{
		Servers:  servers,
		ClientID: l.Name(),
		Username: username,
		Password: password,
	}

	l.Init()
	l.mqtt = m.NewClient(opts)

	l.mqtt.Connect().WaitTimeout(time.Second * 2)

	return l
}

func (l *MQTTListener) Events() []workers.Event {
	return []workers.Event{
		boggart.SecurityOpen,
		boggart.SecurityClosed,
		boggart.DeviceEventHikvisionEventNotificationAlert,
		boggart.DeviceEventDeviceDisabledAfterCheck,
		boggart.DeviceEventDeviceEnabledAfterCheck,
		boggart.DeviceEventDevicesManagerReady,
		boggart.DeviceEventWifiClientConnected,
		boggart.DeviceEventWifiClientDisconnected,
		boggart.DeviceEventVPNClientConnected,
		boggart.DeviceEventVPNClientDisconnected,
		boggart.DeviceEventSoftVideoBalanceChanged,
		boggart.DeviceEventMegafonBalanceChanged,
		devices.EventDoorGPIOReedSwitchOpen,
		devices.EventDoorGPIOReedSwitchClose,
	}
}

func (l *MQTTListener) Run(_ context.Context, event workers.Event, t time.Time, args ...interface{}) {
	switch event {
	case boggart.DeviceEventWifiClientConnected:
		mac := args[1].(*devices.MikrotikRouterMac)
		macAddress := strings.Replace(mac.Address, ":", "-", -1)

		l.mqtt.Publish(fmt.Sprintf("boggart/wifi/clients/%s/ip", macAddress), 0, false, mac.ARP.IP)
		l.mqtt.Publish(fmt.Sprintf("boggart/wifi/clients/%s/comment", macAddress), 0, false, mac.ARP.Comment)
		l.mqtt.Publish(fmt.Sprintf("boggart/wifi/clients/%s/hostname", macAddress), 0, false, mac.DHCP.Hostname)
		l.mqtt.Publish(fmt.Sprintf("boggart/wifi/clients/%s/interface", macAddress), 0, false, args[2])
		l.mqtt.Publish(fmt.Sprintf("boggart/wifi/clients/%s/state", macAddress), 0, true, "ON")

	case boggart.DeviceEventWifiClientDisconnected:
		mac := args[1].(*devices.MikrotikRouterMac)
		macAddress := strings.Replace(mac.Address, ":", "-", -1)

		l.mqtt.Publish(fmt.Sprintf("boggart/wifi/clients/%s/state", macAddress), 0, true, "OFF")

	case boggart.DeviceEventVPNClientConnected:
		l.mqtt.Publish(fmt.Sprintf("boggart/vpn/clients/%s/ip", args[1]), 0, false, args[2])
		l.mqtt.Publish(fmt.Sprintf("boggart/vpn/clients/%s/state", args[1]), 0, true, "ON")

	case boggart.DeviceEventVPNClientDisconnected:
		l.mqtt.Publish(fmt.Sprintf("boggart/vpn/clients/%s/state", args[1]), 0, true, "OFF")

	case boggart.DeviceEventSoftVideoBalanceChanged:
		l.mqtt.Publish(fmt.Sprintf("boggart/service/softvideo/%s/balance", args[2]), 0, true, l.float64(args[1].(float64)))

	case boggart.DeviceEventMegafonBalanceChanged:
		l.mqtt.Publish(fmt.Sprintf("boggart/service/megafon/%s/balance", args[2]), 0, true, l.float64(args[1].(float64)))
	}
}

func (l *MQTTListener) Name() string {
	return boggart.ComponentName + ".mqtt"
}

func (l *MQTTListener) float64(value float64) string {
	return fmt.Sprintf("%.2f", value)
}
