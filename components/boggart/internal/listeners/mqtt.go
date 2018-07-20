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
		boggart.DeviceEventWifiClientConnected,
		boggart.DeviceEventWifiClientDisconnected,
		boggart.DeviceEventVPNClientConnected,
		boggart.DeviceEventVPNClientDisconnected,
		boggart.DeviceEventSoftVideoBalanceChanged,
		boggart.DeviceEventMegafonBalanceChanged,
		boggart.DeviceEventPulsarChanged,
		boggart.DeviceEventMercury200Changed,
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

	case boggart.DeviceEventPulsarChanged:
		values := args[1].(devices.PulsarHeadMeterChange)

		l.mqtt.Publish(fmt.Sprintf("boggart/meter/pulsar/%s/temperature_in", args[2]), 0, true, l.float64(values.TemperatureIn))
		l.mqtt.Publish(fmt.Sprintf("boggart/meter/pulsar/%s/temperature_out", args[2]), 0, true, l.float64(values.TemperatureOut))
		l.mqtt.Publish(fmt.Sprintf("boggart/meter/pulsar/%s/temperature_delta", args[2]), 0, true, l.float64(values.TemperatureDelta))
		l.mqtt.Publish(fmt.Sprintf("boggart/meter/pulsar/%s/energy", args[2]), 0, true, l.float64(values.Energy))
		l.mqtt.Publish(fmt.Sprintf("boggart/meter/pulsar/%s/consumption", args[2]), 0, true, l.float64(values.Consumption))

	case boggart.DeviceEventMercury200Changed:
		values := args[1].(devices.Mercury200ElectricityMeterChange)

		l.mqtt.Publish(fmt.Sprintf("boggart/meter/mercury200/%s/tariff_1", args[2]), 0, true, l.float64(values.Tariff1))
		l.mqtt.Publish(fmt.Sprintf("boggart/meter/mercury200/%s/tariff_2", args[2]), 0, true, l.float64(values.Tariff2))
		l.mqtt.Publish(fmt.Sprintf("boggart/meter/mercury200/%s/tariff_3", args[2]), 0, true, l.float64(values.Tariff3))
		l.mqtt.Publish(fmt.Sprintf("boggart/meter/mercury200/%s/tariff_4", args[2]), 0, true, l.float64(values.Tariff4))
		l.mqtt.Publish(fmt.Sprintf("boggart/meter/mercury200/%s/voltage", args[2]), 0, true, l.float64(values.Voltage))
		l.mqtt.Publish(fmt.Sprintf("boggart/meter/mercury200/%s/amperage", args[2]), 0, true, l.float64(values.Amperage))
		l.mqtt.Publish(fmt.Sprintf("boggart/meter/mercury200/%s/power", args[2]), 0, true, l.float64(values.Power))
		l.mqtt.Publish(fmt.Sprintf("boggart/meter/mercury200/%s/battery_voltage", args[2]), 0, true, l.float64(values.BatteryVoltage))
	}
}

func (l *MQTTListener) Name() string {
	return boggart.ComponentName + ".mqtt"
}

func (l *MQTTListener) float64(value float64) string {
	return fmt.Sprintf("%.2f", value)
}
