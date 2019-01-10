package bind

import (
	"bytes"
	"context"
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
	BroadlinkSP3SUpdateInterval = time.Second * 3 // as e-control app, refresh every 3 sec

	BroadlinkSP3SMQTTTopicState mqtt.Topic = boggart.ComponentName + "/socket/+/state"
	BroadlinkSP3SMQTTTopicPower mqtt.Topic = boggart.ComponentName + "/socket/+/power"
	BroadlinkSP3SMQTTTopicSet   mqtt.Topic = boggart.ComponentName + "/socket/+/set"
)

var (
	metricBroadlinkSP3SPower = snitch.NewGauge(boggart.ComponentName+"_device_socket_broadlink_sp3s_power_watt", "Broadlink SP3S socket current power")
)

type BroadlinkSP3S struct {
	state int64
	power int64

	boggart.DeviceBindBase
	boggart.DeviceBindMQTT

	provider *broadlink.SP3S
}

type BroadlinkSP3SConfig struct {
	IP  string `valid:"ip,required"`
	MAC string `valid:"mac,required"`
}

func (d BroadlinkSP3S) Config() interface{} {
	return &BroadlinkSP3SConfig{}
}

func (d BroadlinkSP3S) CreateBind(c interface{}) (boggart.DeviceBind, error) {
	config := c.(*BroadlinkSP3SConfig)

	localAddr, err := broadlink.LocalAddr()
	if err != nil {
		return nil, err
	}

	mac, err := net.ParseMAC(config.MAC)
	if err != nil {
		return nil, err
	}

	ip := net.UDPAddr{
		IP:   net.ParseIP(config.IP),
		Port: broadlink.DevicePort,
	}

	device := &BroadlinkSP3S{
		provider: broadlink.NewSP3S(mac, ip, *localAddr),
		state:    0,
		power:    -1,
	}
	device.Init()
	device.SetSerialNumber(mac.String())

	return device, nil
}

func (d *BroadlinkSP3S) Describe(ch chan<- *snitch.Description) {
	serialNumber := d.SerialNumber()

	metricBroadlinkSP3SPower.With("serial_number", serialNumber).Describe(ch)
}

func (d *BroadlinkSP3S) Collect(ch chan<- snitch.Metric) {
	serialNumber := d.SerialNumber()

	metricBroadlinkSP3SPower.With("serial_number", serialNumber).Collect(ch)
}

func (d *BroadlinkSP3S) Tasks() []workers.Task {
	taskUpdater := task.NewFunctionTask(d.taskStateUpdater)
	taskUpdater.SetRepeats(-1)
	taskUpdater.SetRepeatInterval(BroadlinkSP3SUpdateInterval)
	taskUpdater.SetName("bind-broadlink-sp3s-updater-" + d.SerialNumber())

	return []workers.Task{
		taskUpdater,
	}
}

func (d *BroadlinkSP3S) taskStateUpdater(ctx context.Context) (interface{}, error) {
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

		d.MQTTPublishAsync(ctx, BroadlinkSP3SMQTTTopicState.Format(serialNumberMQTT), 0, true, mqttValue)
	}

	value, err := d.Power()
	if err != nil {
		return nil, nil
	}

	metricBroadlinkSP3SPower.With("serial_number", serialNumber).Set(value)

	currentPower := int64(value * 100)
	prevPower := atomic.LoadInt64(&d.power)

	if currentPower != prevPower {
		atomic.StoreInt64(&d.power, currentPower)

		d.MQTTPublishAsync(ctx, BroadlinkSP3SMQTTTopicPower.Format(serialNumberMQTT), 0, true, value)
	}

	return nil, nil
}

func (d *BroadlinkSP3S) State() (bool, error) {
	return d.provider.State()
}

func (d *BroadlinkSP3S) On(ctx context.Context) error {
	err := d.provider.On()
	if err == nil {
		_, err = d.taskStateUpdater(ctx)
	}

	return err
}

func (d *BroadlinkSP3S) Off(ctx context.Context) error {
	err := d.provider.Off()
	if err == nil {
		_, err = d.taskStateUpdater(ctx)
	}

	return err
}

func (d *BroadlinkSP3S) Power() (float64, error) {
	return d.provider.Power()
}

func (d *BroadlinkSP3S) MQTTTopics() []mqtt.Topic {
	sn := d.SerialNumberMQTTEscaped()

	return []mqtt.Topic{
		mqtt.Topic(BroadlinkSP3SMQTTTopicState.Format(sn)),
		mqtt.Topic(BroadlinkSP3SMQTTTopicPower.Format(sn)),
		mqtt.Topic(BroadlinkSP3SMQTTTopicSet.Format(sn)),
	}
}

func (d *BroadlinkSP3S) MQTTSubscribers() []mqtt.Subscriber {
	sn := d.SerialNumberMQTTEscaped()

	return []mqtt.Subscriber{
		mqtt.NewSubscriber(BroadlinkSP3SMQTTTopicSet.Format(sn), 0, func(ctx context.Context, client mqtt.Component, message mqtt.Message) {
			if d.Status() != boggart.DeviceStatusOnline {
				return
			}

			span, ctx := tracing.StartSpanFromContext(ctx, "socket", "set")
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
