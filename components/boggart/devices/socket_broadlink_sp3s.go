package devices

import (
	"bytes"
	"context"
	"strings"
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

func NewBroadlinkSP3SSocket(provider *broadlink.SP3S) *BroadlinkSP3SSocket {
	device := &BroadlinkSP3SSocket{
		provider: provider,
		state:    0,
		power:    -1,
	}
	device.Init()
	device.SetSerialNumber(provider.MAC().String())
	device.SetDescription("Socket of Broadlink")

	return device
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
	serialNumberMQTT := strings.Replace(serialNumber, ":", "-", -1)

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
	return []mqtt.Topic{
		SocketBroadlinkSP3SMQTTTopicState,
		SocketBroadlinkSP3SMQTTTopicPower,
		SocketBroadlinkSP3SMQTTTopicSet,
	}
}

func (d *BroadlinkSP3SSocket) MQTTSubscribers() []mqtt.Subscriber {
	mac := strings.Replace(d.provider.MAC().String(), ":", "-", -1)

	return []mqtt.Subscriber{
		mqtt.NewSubscriber(SocketBroadlinkSP3SMQTTTopicSet.Format(mac), 0, func(ctx context.Context, client mqtt.Component, message mqtt.Message) {
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
