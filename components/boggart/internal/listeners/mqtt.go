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
	"github.com/kihamo/boggart/components/boggart/providers/hikvision"
	"github.com/kihamo/go-workers"
	"github.com/kihamo/go-workers/listener"
)

const (
	TopicPrefix = boggart.ComponentName + "/"

	ValueOff = 0
	ValueOn  = 1
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
		boggart.DeviceEventHikvisionEventNotificationAlert,
		boggart.DeviceEventWifiClientConnected,
		boggart.DeviceEventWifiClientDisconnected,
		boggart.DeviceEventVPNClientConnected,
		boggart.DeviceEventVPNClientDisconnected,
		boggart.DeviceEventSoftVideoBalanceChanged,
		boggart.DeviceEventMegafonBalanceChanged,
		boggart.DeviceEventPulsarChanged,
		boggart.DeviceEventPulsarPulsedChanged,
		boggart.DeviceEventMercury200Changed,
		boggart.DeviceEventBME280Changed,
		boggart.DeviceEventGPIOPinChanged,
	}
}

func (l *MQTTListener) Run(_ context.Context, event workers.Event, t time.Time, args ...interface{}) {
	switch event {
	case boggart.DeviceEventHikvisionEventNotificationAlert:
		event := args[1].(*hikvision.EventNotificationAlertStreamResponse)
		id := strings.Replace(args[2].(string), "/", "-", -1)

		l.publish(fmt.Sprintf("cctv/%s/%d/%s", id, event.DynChannelID, event.EventType), false, event.EventDescription)

	case boggart.DeviceEventWifiClientConnected:
		mac := args[1].(*devices.MikrotikRouterMac)
		macAddress := l.macAddress(mac.Address)

		l.publish(fmt.Sprintf("wifi/clients/%s/ip", macAddress), false, mac.ARP.IP)
		l.publish(fmt.Sprintf("wifi/clients/%s/comment", macAddress), false, mac.ARP.Comment)
		l.publish(fmt.Sprintf("wifi/clients/%s/hostname", macAddress), false, mac.DHCP.Hostname)
		l.publish(fmt.Sprintf("wifi/clients/%s/interface", macAddress), false, args[2])
		l.publish(fmt.Sprintf("wifi/clients/%s/state", macAddress), true, ValueOn)

	case boggart.DeviceEventWifiClientDisconnected:
		mac := args[1].(*devices.MikrotikRouterMac)
		macAddress := l.macAddress(mac.Address)

		l.publish(fmt.Sprintf("wifi/clients/%s/state", macAddress), true, ValueOff)

	case boggart.DeviceEventVPNClientConnected:
		macAddress := l.macAddress(args[1].(string))

		l.publish(fmt.Sprintf("vpn/clients/%s/ip", macAddress), false, args[2])
		l.publish(fmt.Sprintf("vpn/clients/%s/state", macAddress), true, ValueOn)

	case boggart.DeviceEventVPNClientDisconnected:
		macAddress := l.macAddress(args[1].(string))

		l.publish(fmt.Sprintf("vpn/clients/%s/state", macAddress), true, ValueOff)

	case boggart.DeviceEventSoftVideoBalanceChanged:
		l.publish(fmt.Sprintf("service/softvideo/%s/balance", args[2]), true, args[1])

	case boggart.DeviceEventMegafonBalanceChanged:
		l.publish(fmt.Sprintf("service/megafon/%s/balance", args[2]), true, args[1])

	case boggart.DeviceEventPulsarChanged:
		values := args[1].(devices.PulsarHeadMeterChange)

		l.publish(fmt.Sprintf("meter/pulsar/%s/temperature_in", args[2]), true, values.TemperatureIn)
		l.publish(fmt.Sprintf("meter/pulsar/%s/temperature_out", args[2]), true, values.TemperatureOut)
		l.publish(fmt.Sprintf("meter/pulsar/%s/temperature_delta", args[2]), true, values.TemperatureDelta)
		l.publish(fmt.Sprintf("meter/pulsar/%s/energy", args[2]), true, values.Energy)
		l.publish(fmt.Sprintf("meter/pulsar/%s/consumption", args[2]), true, values.Consumption)

	case boggart.DeviceEventPulsarPulsedChanged:
		values := args[1].(devices.PulsarPulsedWaterMeterChanged)

		l.publish(fmt.Sprintf("meter/pulsar/%s/volume", args[2]), true, values.Volume)
		l.publish(fmt.Sprintf("meter/pulsar/%s/pulses", args[2]), true, values.Pulses)

	case boggart.DeviceEventMercury200Changed:
		values := args[1].(devices.Mercury200ElectricityMeterChange)

		l.publish(fmt.Sprintf("meter/mercury200/%s/tariff_1", args[2]), true, values.Tariff1)
		l.publish(fmt.Sprintf("meter/mercury200/%s/tariff_2", args[2]), true, values.Tariff2)
		l.publish(fmt.Sprintf("meter/mercury200/%s/tariff_3", args[2]), true, values.Tariff3)
		l.publish(fmt.Sprintf("meter/mercury200/%s/tariff_4", args[2]), true, values.Tariff4)
		l.publish(fmt.Sprintf("meter/mercury200/%s/voltage", args[2]), true, values.Voltage)
		l.publish(fmt.Sprintf("meter/mercury200/%s/amperage", args[2]), true, values.Amperage)
		l.publish(fmt.Sprintf("meter/mercury200/%s/power", args[2]), true, values.Power)
		l.publish(fmt.Sprintf("meter/mercury200/%s/battery_voltage", args[2]), true, values.BatteryVoltage)

	case boggart.DeviceEventBME280Changed:
		values := args[1].(devices.BME280SensorChange)

		l.publish(fmt.Sprintf("meter/bme280/%s/temperature", args[2]), true, values.Temperature)
		l.publish(fmt.Sprintf("meter/bme280/%s/altitude", args[2]), true, values.Altitude)
		l.publish(fmt.Sprintf("meter/bme280/%s/humidity", args[2]), true, values.Humidity)
		l.publish(fmt.Sprintf("meter/bme280/%s/pressure", args[2]), true, values.Pressure)

	case boggart.DeviceEventGPIOPinChanged:
		if args[2].(bool) {
			l.publish(fmt.Sprintf("gpio/%d", args[1]), true, ValueOn)
		} else {
			l.publish(fmt.Sprintf("gpio/%d", args[1]), true, ValueOff)
		}
	}
}

func (l *MQTTListener) Name() string {
	return boggart.ComponentName + ".mqtt"
}

func (l *MQTTListener) publish(topic string, retained bool, payload interface{}) {
	switch value := payload.(type) {
	case float64:
		payload = fmt.Sprintf("%.2f", value)
	case uint64, int64, int:
		payload = fmt.Sprintf("%d", value)
	}

	l.mqtt.Publish(TopicPrefix+topic, 0, retained, payload)
}

func (l *MQTTListener) macAddress(address string) string {
	return strings.Replace(address, ":", "-", -1)
}
