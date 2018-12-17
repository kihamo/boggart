package devices

import (
	"bytes"
	"context"
	"strings"
	"sync/atomic"
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/mqtt"
	"github.com/kihamo/go-workers"
	"github.com/kihamo/go-workers/task"
	"github.com/vikstrous/zengge-lightcontrol/control"
	"github.com/vikstrous/zengge-lightcontrol/local"
	"github.com/vikstrous/zengge-lightcontrol/remote"
)

const (
	ZenggeLEDPortCommand = 5577
	ZenggeLEDPortWifi    = 48899

	ZenggeLEDUpdateInterval = time.Second * 3

	ZenggeLEDMQTTTopicPower      mqtt.Topic = boggart.ComponentName + "/led/+/power"
	ZenggeLEDMQTTTopicColor      mqtt.Topic = boggart.ComponentName + "/led/+/color"
	ZenggeLEDMQTTTopicStatePower mqtt.Topic = boggart.ComponentName + "/led/+/state/power"
	ZenggeLEDMQTTTopicStateColor mqtt.Topic = boggart.ComponentName + "/led/+/state/color"
	ZenggeLEDMQTTTopicStateMode  mqtt.Topic = boggart.ComponentName + "/led/+/state/mode"
)

type ZenggeLED struct {
	statePower int64
	stateMode  uint64

	boggart.DeviceBase
	boggart.DeviceSerialNumber
	boggart.DeviceMQTT

	controller *control.Controller
}

func NewZenggeLED(controller *control.Controller) *ZenggeLED {
	device := &ZenggeLED{
		controller: controller,
	}
	device.Init()
	device.SetDescription("LED Zengge")

	var sn string

	switch l := controller.Transport.(type) {
	case *local.Transport:
		address := l.Conn.RemoteAddr().String()
		index := strings.IndexByte(address, ':')
		if index == -1 {
			sn = address
		} else if i := strings.IndexByte(address, ']'); i != -1 {
			sn = strings.TrimPrefix(address[:i], "[")
		} else {
			sn = address[:index]
		}

		sn = strings.Replace(sn, ".", "-", -1)

	case *remote.RemoteTransport:
		sn = strings.Replace(l.MAC, ":", "-", -1)
	}

	if sn != "" {
		device.SetSerialNumber(sn)
		device.SetDescription(device.Description() + " " + sn)
	}

	return device
}

func (d *ZenggeLED) Types() []boggart.DeviceType {
	return []boggart.DeviceType{
		boggart.DeviceTypeLED,
	}
}

func (d *ZenggeLED) Ping(_ context.Context) bool {
	_, err := d.controller.GetState()
	return err == nil
}

func (d *ZenggeLED) Tasks() []workers.Task {
	taskUpdater := task.NewFunctionTask(d.taskUpdater)
	taskUpdater.SetRepeats(-1)
	taskUpdater.SetRepeatInterval(ZenggeLEDUpdateInterval)
	taskUpdater.SetName("device-led-zengge-updater-" + d.SerialNumber())

	return []workers.Task{
		taskUpdater,
	}
}

func (d *ZenggeLED) taskUpdater(ctx context.Context) (interface{}, error) {
	if !d.IsEnabled() {
		return nil, nil
	}

	state, err := d.controller.GetState()
	if err != nil {
		return nil, err
	}

	sn := d.SerialNumber()

	prevPower := atomic.LoadInt64(&d.statePower)
	if prevPower == 0 || (prevPower == 1) != state.IsOn {
		var mqttValue []byte

		if state.IsOn {
			atomic.StoreInt64(&d.statePower, 1)
			mqttValue = []byte(`1`)
		} else {
			atomic.StoreInt64(&d.statePower, -1)
			mqttValue = []byte(`0`)
		}

		d.MQTTPublishAsync(ctx, ZenggeLEDMQTTTopicStatePower.Format(sn), 0, true, mqttValue)
	}

	currentMode := uint64(state.Mode)
	prevMode := atomic.LoadUint64(&d.stateMode)
	if prevMode != currentMode {
		atomic.StoreUint64(&d.stateMode, currentMode)

		d.MQTTPublishAsync(ctx, ZenggeLEDMQTTTopicStateMode.Format(sn), 0, true, currentMode)
	}

	d.MQTTPublishAsync(ctx, ZenggeLEDMQTTTopicStateColor.Format(sn), 0, true, state.Color.Format())

	return nil, nil
}

func (d *ZenggeLED) On(ctx context.Context) error {
	err := d.controller.SetPower(true)
	if err == nil {
		_, err = d.taskUpdater(ctx)
	}

	return err
}

func (d *ZenggeLED) Off(ctx context.Context) error {
	err := d.controller.SetPower(false)
	if err == nil {
		_, err = d.taskUpdater(ctx)
	}

	return err
}

func (d *ZenggeLED) MQTTTopics() []mqtt.Topic {
	return []mqtt.Topic{
		ZenggeLEDMQTTTopicPower,
		ZenggeLEDMQTTTopicColor,
		ZenggeLEDMQTTTopicStatePower,
		ZenggeLEDMQTTTopicStateColor,
		ZenggeLEDMQTTTopicStateMode,
	}
}

func (d *ZenggeLED) MQTTSubscribers() []mqtt.Subscriber {
	sn := d.SerialNumber()

	return []mqtt.Subscriber{
		mqtt.NewSubscriber(ZenggeLEDMQTTTopicPower.Format(sn), 0, func(ctx context.Context, client mqtt.Component, message mqtt.Message) {
			if !d.IsEnabled() {
				return
			}

			if bytes.Equal(message.Payload(), []byte(`1`)) {
				d.On(ctx)
			} else {
				d.Off(ctx)
			}
		}),
		mqtt.NewSubscriber(ZenggeLEDMQTTTopicColor.Format(sn), 0, func(ctx context.Context, client mqtt.Component, message mqtt.Message) {
			if !d.IsEnabled() {
				return
			}

			color := control.ParseColorString(string(message.Payload()))
			if color != nil {
				d.controller.SetColor(*color)
			}
		}),
	}
}
