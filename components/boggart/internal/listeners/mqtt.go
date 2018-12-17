package listeners

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/devices"
	"github.com/kihamo/boggart/components/boggart/providers/hikvision"
	"github.com/kihamo/boggart/components/mqtt"
	"github.com/kihamo/go-workers"
	"github.com/kihamo/go-workers/listener"
	"github.com/kihamo/shadow/components/tracing"
)

const (
	TopicPrefix = boggart.ComponentName + "/"

	ValueOff = 0
	ValueOn  = 1
)

type MQTTListener struct {
	listener.BaseListener

	client mqtt.Component
}

func NewMQTTListener(client mqtt.Component) *MQTTListener {
	l := &MQTTListener{
		client: client,
	}
	l.Init()

	return l
}

func (l *MQTTListener) Events() []workers.Event {
	return []workers.Event{
		boggart.DeviceEventHikvisionEventNotificationAlert,
		boggart.DeviceEventWifiClientConnected,
		boggart.DeviceEventWifiClientDisconnected,
		boggart.DeviceEventVPNClientConnected,
		boggart.DeviceEventVPNClientDisconnected,
		boggart.DeviceEventMegafonBalanceChanged,
		boggart.DeviceEventPulsarChanged,
		boggart.DeviceEventMercury200Changed,
		boggart.DeviceEventBME280Changed,
		boggart.DeviceEventGPIOPinChanged,
	}
}

func (l *MQTTListener) Run(ctx context.Context, event workers.Event, t time.Time, args ...interface{}) {
	span, ctx := tracing.StartSpanFromContext(ctx, "mqtt_listener", "run")
	span = span.SetTag("event", event.Name())
	defer span.Finish()

	switch event {
	case boggart.DeviceEventHikvisionEventNotificationAlert:
		event := args[1].(*hikvision.EventNotificationAlertStreamResponse)
		id := strings.Replace(args[2].(string), "/", "-", -1)

		l.publish(ctx, fmt.Sprintf("cctv/%s/%d/%s", id, event.DynChannelID, event.EventType), false, event.EventDescription)

	case boggart.DeviceEventWifiClientConnected:
		mac := args[1].(*devices.MikrotikRouterMac)
		macAddress := l.macAddress(mac.Address)

		topicPrefix := "router/" + args[3].(string) + "/wifi/clients/last/on/"
		l.publish(ctx, topicPrefix+"mac", false, macAddress)
		l.publish(ctx, topicPrefix+"ip", false, mac.ARP.IP)
		l.publish(ctx, topicPrefix+"comment", false, mac.ARP.Comment)
		l.publish(ctx, topicPrefix+"hostname", false, mac.DHCP.Hostname)
		l.publish(ctx, topicPrefix+"interface", false, args[2])

		topicPrefix = "router/" + args[3].(string) + "/wifi/clients/" + macAddress + "/"
		l.publish(ctx, topicPrefix+"ip", false, mac.ARP.IP)
		l.publish(ctx, topicPrefix+"comment", false, mac.ARP.Comment)
		l.publish(ctx, topicPrefix+"hostname", false, mac.DHCP.Hostname)
		l.publish(ctx, topicPrefix+"interface", false, args[2])
		l.publish(ctx, topicPrefix+"state", true, ValueOn)

	case boggart.DeviceEventWifiClientDisconnected:
		mac := args[1].(*devices.MikrotikRouterMac)
		macAddress := l.macAddress(mac.Address)

		topicPrefix := "router/" + args[3].(string) + "/wifi/clients/last/off/"
		l.publish(ctx, topicPrefix+"mac", false, macAddress)
		l.publish(ctx, topicPrefix+"ip", false, mac.ARP.IP)
		l.publish(ctx, topicPrefix+"comment", false, mac.ARP.Comment)
		l.publish(ctx, topicPrefix+"hostname", false, mac.DHCP.Hostname)
		l.publish(ctx, topicPrefix+"interface", false, args[2])

		topicPrefix = "router/" + args[3].(string) + "/wifi/clients/" + macAddress + "/"
		l.publish(ctx, topicPrefix+"ip", false, mac.ARP.IP)
		l.publish(ctx, topicPrefix+"comment", false, mac.ARP.Comment)
		l.publish(ctx, topicPrefix+"hostname", false, mac.DHCP.Hostname)
		l.publish(ctx, topicPrefix+"interface", false, args[2])
		l.publish(ctx, topicPrefix+"state", true, ValueOff)

	case boggart.DeviceEventVPNClientConnected:
		login := l.macAddress(args[1].(string))

		topicPrefix := "router/" + args[3].(string) + "/vpn/clients/last/on/"
		l.publish(ctx, topicPrefix+"login", false, login)
		l.publish(ctx, topicPrefix+"ip", false, args[2])

		topicPrefix = "router/" + args[3].(string) + "/vpn/clients/" + login + "/"
		l.publish(ctx, topicPrefix+"ip", false, args[2])
		l.publish(ctx, topicPrefix+"state", false, ValueOn)

	case boggart.DeviceEventVPNClientDisconnected:
		login := l.macAddress(args[1].(string))
		topicPrefix := "router/" + args[2].(string) + "/vpn/clients/"

		l.publish(ctx, topicPrefix+"last/off/login", false, login)
		l.publish(ctx, topicPrefix+login+"/state", true, ValueOff)

	case boggart.DeviceEventMegafonBalanceChanged:
		l.publish(ctx, fmt.Sprintf("service/megafon/%s/balance", args[2]), true, args[1])

	case boggart.DeviceEventPulsarChanged:
		values := args[1].(devices.PulsarHeadMeterChange)

		l.publish(ctx, fmt.Sprintf("meter/pulsar/%s/temperature_in", args[2]), true, values.TemperatureIn)
		l.publish(ctx, fmt.Sprintf("meter/pulsar/%s/temperature_out", args[2]), true, values.TemperatureOut)
		l.publish(ctx, fmt.Sprintf("meter/pulsar/%s/temperature_delta", args[2]), true, values.TemperatureDelta)
		l.publish(ctx, fmt.Sprintf("meter/pulsar/%s/energy", args[2]), true, values.Energy)
		l.publish(ctx, fmt.Sprintf("meter/pulsar/%s/consumption", args[2]), true, values.Consumption)

	case boggart.DeviceEventMercury200Changed:
		values := args[1].(devices.Mercury200ElectricityMeterChange)

		l.publish(ctx, fmt.Sprintf("meter/mercury200/%s/tariff_1", args[2]), true, values.Tariff1)
		l.publish(ctx, fmt.Sprintf("meter/mercury200/%s/tariff_2", args[2]), true, values.Tariff2)
		l.publish(ctx, fmt.Sprintf("meter/mercury200/%s/tariff_3", args[2]), true, values.Tariff3)
		l.publish(ctx, fmt.Sprintf("meter/mercury200/%s/tariff_4", args[2]), true, values.Tariff4)
		l.publish(ctx, fmt.Sprintf("meter/mercury200/%s/voltage", args[2]), true, values.Voltage)
		l.publish(ctx, fmt.Sprintf("meter/mercury200/%s/amperage", args[2]), true, values.Amperage)
		l.publish(ctx, fmt.Sprintf("meter/mercury200/%s/power", args[2]), true, values.Power)
		l.publish(ctx, fmt.Sprintf("meter/mercury200/%s/battery_voltage", args[2]), true, values.BatteryVoltage)

	case boggart.DeviceEventBME280Changed:
		values := args[1].(devices.BME280SensorChange)

		l.publish(ctx, fmt.Sprintf("meter/bme280/%s/temperature", args[2]), true, values.Temperature)
		l.publish(ctx, fmt.Sprintf("meter/bme280/%s/altitude", args[2]), true, values.Altitude)
		l.publish(ctx, fmt.Sprintf("meter/bme280/%s/humidity", args[2]), true, values.Humidity)
		l.publish(ctx, fmt.Sprintf("meter/bme280/%s/pressure", args[2]), true, values.Pressure)

	case boggart.DeviceEventGPIOPinChanged:
		if args[2].(bool) {
			l.publishWithQOS(ctx, fmt.Sprintf("gpio/%d", args[1]), 2, true, ValueOn)
		} else {
			l.publishWithQOS(ctx, fmt.Sprintf("gpio/%d", args[1]), 2, true, ValueOff)
		}
	}
}

func (l *MQTTListener) Name() string {
	return boggart.ComponentName + ".mqtt"
}

func (l *MQTTListener) publishWithQOS(ctx context.Context, topic string, qos byte, retained bool, payload interface{}) {
	go func() {
		switch value := payload.(type) {
		case float64:
			payload = fmt.Sprintf("%.2f", value)
		case uint64, int64, int:
			payload = fmt.Sprintf("%d", value)
		}

		l.client.Publish(ctx, TopicPrefix+topic, qos, retained, payload)
	}()
}

func (l *MQTTListener) publish(ctx context.Context, topic string, retained bool, payload interface{}) {
	l.publishWithQOS(ctx, topic, 0, retained, payload)
}

func (l *MQTTListener) macAddress(address string) string {
	address = strings.Replace(address, ":", "-", -1)
	address = strings.Replace(address, ",", "-", -1)

	return address
}
