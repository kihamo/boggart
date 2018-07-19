package listeners

import (
	"context"
	"fmt"
	"net/url"
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
		devices.EventDoorGPIOReedSwitchOpen,
		devices.EventDoorGPIOReedSwitchClose,
	}
}

func (l *MQTTListener) Run(_ context.Context, event workers.Event, t time.Time, args ...interface{}) {
	switch event {
	case boggart.DeviceEventWifiClientConnected:
		mac := args[1].(*devices.MikrotikRouterMac)

		l.mqtt.Publish(fmt.Sprintf("/boggart/wifi/clients/%s/ip", mac.Address), 0, false, mac.ARP.IP)
		l.mqtt.Publish(fmt.Sprintf("/boggart/wifi/clients/%s/comment", mac.Address), 0, false, mac.ARP.Comment)
		l.mqtt.Publish(fmt.Sprintf("/boggart/wifi/clients/%s/hostname", mac.Address), 0, false, mac.DHCP.Hostname)
		l.mqtt.Publish(fmt.Sprintf("/boggart/wifi/clients/%s/interface", mac.Address), 0, false, args[2])
		l.mqtt.Publish(fmt.Sprintf("/boggart/wifi/clients/%s/status", mac.Address), 0, false, "ONLINE")

	case boggart.DeviceEventWifiClientDisconnected:
		mac := args[1].(*devices.MikrotikRouterMac)

		l.mqtt.Publish(fmt.Sprintf("/boggart/wifi/clients/%s/status", mac.Address), 0, false, "OFFLINE")

	case boggart.DeviceEventVPNClientConnected:
		l.mqtt.Publish(fmt.Sprintf("/boggart/vpn/clients/%s/ip", args[1]), 0, false, args[2])
		l.mqtt.Publish(fmt.Sprintf("/boggart/vpn/clients/%s/status", args[1]), 0, false, "ONLINE")

	case boggart.DeviceEventVPNClientDisconnected:
		l.mqtt.Publish(fmt.Sprintf("/boggart/vpn/clients/%s/status", args[1]), 0, false, "OFFLINE")
	}
}

func (l *MQTTListener) Name() string {
	return boggart.ComponentName + ".mqtt"
}
