package devices

import (
	"bytes"
	"context"
	"errors"
	"net"
	"sync/atomic"
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/providers/broadlink"
	"github.com/kihamo/boggart/components/mqtt"
	"github.com/kihamo/go-workers"
	"github.com/kihamo/go-workers/task"
	"github.com/kihamo/shadow/components/tracing"
	"github.com/kihamo/snitch"
	"github.com/opentracing/opentracing-go/log"
)

const (
	SocketBroadlinkSP3SUpdateInterval = time.Second * 3 // as e-control app, refresh every 3 sec

	SocketBroadlinkSP3SMQTTTopicState mqtt.Topic = boggart.ComponentName + "/socket/+/state"
	SocketBroadlinkSP3SMQTTTopicPower mqtt.Topic = boggart.ComponentName + "/socket/+/power"
	SocketBroadlinkSP3SMQTTTopicSet   mqtt.Topic = boggart.ComponentName + "/socket/+/set"
)

var (
	metricSocketBroadlinkSP3SPower = snitch.NewGauge(boggart.ComponentName+"_device_socket_broadlink_sp3s_power_watt", "Broadlink SP3S socket current power")
)

type BroadlinkSP3SSocket struct {
	state int64
	power int64

	boggart.DeviceBase
	boggart.DeviceSerialNumber
	boggart.DeviceMQTT

	provider *broadlink.SP3S
}

func (d BroadlinkSP3SSocket) Create(config map[string]interface{}) (boggart.Device, error) {
	localAddr, err := broadlink.LocalAddr()
	if err != nil {
		return nil, err
	}

	ipConfig, ok := config["ip"]
	if !ok {
		return nil, errors.New("config option ip isn't set")
	}

	if ipConfig == "" {
		return nil, errors.New("config option ip is empty")
	}

	macConfig, ok := config["mac"]
	if !ok {
		return nil, errors.New("config option mac isn't set")
	}

	if macConfig == "" {
		return nil, errors.New("config option mac is empty")
	}

	mac, err := net.ParseMAC(macConfig.(string))
	if err != nil {
		return nil, err
	}

	ip := net.UDPAddr{
		IP:   net.ParseIP(ipConfig.(string)),
		Port: broadlink.DevicePort,
	}

	device := &BroadlinkSP3SSocket{
		provider: broadlink.NewSP3S(mac, ip, *localAddr),
		state:    0,
		power:    -1,
	}
	device.Init()
	device.SetSerialNumber(mac.String())
	device.SetDescription("Socket of Broadlink")

	return device, nil
}

func (d *BroadlinkSP3SSocket) Types() []boggart.DeviceType {
	return []boggart.DeviceType{
		boggart.DeviceTypeSocket,
	}
}

func (d *BroadlinkSP3SSocket) Describe(ch chan<- *snitch.Description) {
	serialNumber := d.SerialNumber()

	metricSocketBroadlinkSP3SPower.With("serial_number", serialNumber).Describe(ch)
}

func (d *BroadlinkSP3SSocket) Collect(ch chan<- snitch.Metric) {
	serialNumber := d.SerialNumber()

	metricSocketBroadlinkSP3SPower.With("serial_number", serialNumber).Collect(ch)
}

func (d *BroadlinkSP3SSocket) Tasks() []workers.Task {
	taskUpdater := task.NewFunctionTask(d.taskStateUpdater)
	taskUpdater.SetRepeats(-1)
	taskUpdater.SetRepeatInterval(SocketBroadlinkSP3SUpdateInterval)
	taskUpdater.SetName("device-socket-broadlink-sp3s-updater-" + d.SerialNumber())

	return []workers.Task{
		taskUpdater,
	}
}

func (d *BroadlinkSP3SSocket) taskStateUpdater(ctx context.Context) (interface{}, error) {
	state, err := d.State()
	if err != nil {
		d.UpdateStatus(boggart.DeviceStatusOffline)
		return nil, err
	}

	d.UpdateStatus(boggart.DeviceStatusOnline)

	serialNumber := d.SerialNumber()
	serialNumberMQTT := d.SerialNumberMQTTEscaped()

	prevState := atomic.LoadInt64(&d.state)
	if prevState == 0 || (prevState == 1) != state {
		var mqttValue []byte

		if state {
			atomic.StoreInt64(&d.state, 1)
			mqttValue = []byte(`1`)
		} else {
			atomic.StoreInt64(&d.state, -1)
			mqttValue = []byte(`0`)
		}

		d.MQTTPublishAsync(ctx, SocketBroadlinkSP3SMQTTTopicState.Format(serialNumberMQTT), 0, true, mqttValue)
	}

	value, err := d.Power()
	if err != nil {
		return nil, nil
	}

	metricSocketBroadlinkSP3SPower.With("serial_number", serialNumber).Set(value)

	currentPower := int64(value * 100)
	prevPower := atomic.LoadInt64(&d.power)

	if currentPower != prevPower {
		atomic.StoreInt64(&d.power, currentPower)

		d.MQTTPublishAsync(ctx, SocketBroadlinkSP3SMQTTTopicPower.Format(serialNumberMQTT), 0, true, value)
	}

	return nil, nil
}

func (d *BroadlinkSP3SSocket) State() (bool, error) {
	return d.provider.State()
}

func (d *BroadlinkSP3SSocket) On(ctx context.Context) error {
	err := d.provider.On()
	if err == nil {
		_, err = d.taskStateUpdater(ctx)
	}

	return err
}

func (d *BroadlinkSP3SSocket) Off(ctx context.Context) error {
	err := d.provider.Off()
	if err == nil {
		_, err = d.taskStateUpdater(ctx)
	}

	return err
}

func (d *BroadlinkSP3SSocket) Power() (float64, error) {
	return d.provider.Power()
}

func (d *BroadlinkSP3SSocket) MQTTTopics() []mqtt.Topic {
	sn := d.SerialNumberMQTTEscaped()

	return []mqtt.Topic{
		mqtt.Topic(SocketBroadlinkSP3SMQTTTopicState.Format(sn)),
		mqtt.Topic(SocketBroadlinkSP3SMQTTTopicPower.Format(sn)),
		mqtt.Topic(SocketBroadlinkSP3SMQTTTopicSet.Format(sn)),
	}
}

func (d *BroadlinkSP3SSocket) MQTTSubscribers() []mqtt.Subscriber {
	sn := d.SerialNumberMQTTEscaped()

	return []mqtt.Subscriber{
		mqtt.NewSubscriber(SocketBroadlinkSP3SMQTTTopicSet.Format(sn), 0, func(ctx context.Context, client mqtt.Component, message mqtt.Message) {
			if d.Status() != boggart.DeviceStatusOnline {
				return
			}

			span, ctx := tracing.StartSpanFromContext(ctx, boggart.DeviceTypeSocket.String(), "set")
			span.LogFields(
				log.String("mac", d.provider.MAC().String()),
				log.String("ip", d.provider.Addr().String()))
			defer span.Finish()

			var err error

			if bytes.Equal(message.Payload(), []byte(`1`)) {
				err = d.On(ctx)
				span.LogFields(log.String("state", "on"))
			} else {
				err = d.Off(ctx)
				span.LogFields(log.String("state", "off"))
			}

			if err != nil {
				tracing.SpanError(span, err)
			}
		}),
	}
}
