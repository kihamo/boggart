package devices

import (
	"bytes"
	"context"
	"fmt"
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
	SocketBroadlinkSP3SUpdateInterval  = time.Second * 3 // // as e-control app, refresh every 3 sec
	SocketBroadlinkSP3SMQTTTopicPrefix = boggart.ComponentName + "/socket/"
)

var (
	metricSocketBroadlinkSP3SPower = snitch.NewGauge(boggart.ComponentName+"_device_socket_broadlink_sp3s_power_watt", "Broadlink SP3S socket current power")
)

type BroadlinkSP3SSocket struct {
	state     int64
	lastValue int64

	boggart.DeviceBase
	boggart.DeviceSerialNumber

	provider *broadlink.SP3S
}

func NewBroadlinkSP3SSocket(provider *broadlink.SP3S) *BroadlinkSP3SSocket {
	device := &BroadlinkSP3SSocket{
		provider: provider,
	}
	device.Init()
	device.SetSerialNumber(provider.MAC().String())
	device.SetDescription("Socket of Broadlink with IP " + provider.Addr().String() + " and MAC" + provider.MAC().String())

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

func (d *BroadlinkSP3SSocket) Ping(_ context.Context) bool {
	_, err := d.Power()
	return err == nil
}

func (d *BroadlinkSP3SSocket) Tasks() []workers.Task {
	taskUpdater := task.NewFunctionTask(d.taskUpdater)
	taskUpdater.SetRepeats(-1)
	taskUpdater.SetRepeatInterval(SocketBroadlinkSP3SUpdateInterval)
	taskUpdater.SetName("device-socket-broadlink-sp3s-updater-" + d.SerialNumber())

	return []workers.Task{
		taskUpdater,
	}
}

func (d *BroadlinkSP3SSocket) taskUpdater(ctx context.Context) (interface{}, error) {
	if !d.IsEnabled() {
		return nil, nil
	}

	state, err := d.State()
	if err != nil {
		return nil, err
	}

	serialNumber := d.SerialNumber()

	last := atomic.LoadInt64(&d.state)
	if last == 0 || (last == 1) != state {
		d.TriggerEvent(boggart.DeviceEventSocketStateChanged, state, serialNumber)

		if state {
			atomic.StoreInt64(&d.state, 1)
		} else {
			atomic.StoreInt64(&d.state, -1)
		}
	}

	value, err := d.Power()
	if err != nil {
		return nil, nil
	}

	metricSocketBroadlinkSP3SPower.With("serial_number", serialNumber).Set(value)

	current := int64(value * 100)
	prev := atomic.LoadInt64(&d.lastValue)

	if current != prev {
		atomic.StoreInt64(&d.lastValue, current)
		d.TriggerEvent(boggart.DeviceEventSocketPowerChanged, value, serialNumber)
	}

	return nil, nil
}

func (d *BroadlinkSP3SSocket) State() (bool, error) {
	return d.provider.State()
}

func (d *BroadlinkSP3SSocket) On() error {
	err := d.provider.On()
	if err == nil {
		d.taskUpdater(context.Background())
	}

	return err
}

func (d *BroadlinkSP3SSocket) Off() error {
	err := d.provider.Off()
	if err == nil {
		d.taskUpdater(context.Background())
	}

	return err
}

func (d *BroadlinkSP3SSocket) Power() (float64, error) {
	return d.provider.Power()
}

func (d *BroadlinkSP3SSocket) MQTTSubscribers() []mqtt.Subscriber {
	mac := strings.Replace(d.provider.MAC().String(), ":", "-", -1)
	topic := fmt.Sprintf("%s%s/set", SocketBroadlinkSP3SMQTTTopicPrefix, mac)

	return []mqtt.Subscriber{
		mqtt.NewSubscriber(topic, 0, func(ctx context.Context, client mqtt.Component, message mqtt.Message) {
			if !d.IsEnabled() {
				return
			}

			span, ctx := tracing.StartSpanFromContext(ctx, "socket", "set")
			span.LogFields(
				log.String("mac", d.provider.MAC().String()),
				log.String("ip", d.provider.Addr().String()))
			defer span.Finish()

			var err error

			if bytes.Equal(message.Payload(), []byte(`1`)) {
				err = d.On()
				span.LogFields(log.String("state", "on"))
			} else {
				err = d.Off()
				span.LogFields(log.String("state", "off"))
			}

			if err != nil {
				tracing.SpanError(span, err)
			}
		}),
	}
}
