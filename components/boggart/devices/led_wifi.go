package devices

import (
	"bytes"
	"context"
	"fmt"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/providers/wifiled"
	"github.com/kihamo/boggart/components/mqtt"
	"github.com/kihamo/go-workers"
	"github.com/kihamo/go-workers/task"
)

const (
	WiFiLEDUpdateInterval = time.Second * 3

	WiFiLEDMQTTTopicPower         mqtt.Topic = boggart.ComponentName + "/led/+/power"
	WiFiLEDMQTTTopicColor         mqtt.Topic = boggart.ComponentName + "/led/+/color"
	WiFiLEDMQTTTopicMode          mqtt.Topic = boggart.ComponentName + "/led/+/mode"
	WiFiLEDMQTTTopicSpeed         mqtt.Topic = boggart.ComponentName + "/led/+/speed"
	WiFiLEDMQTTTopicStatePower    mqtt.Topic = boggart.ComponentName + "/led/+/state/power"
	WiFiLEDMQTTTopicStateColor    mqtt.Topic = boggart.ComponentName + "/led/+/state/color"
	WiFiLEDMQTTTopicStateColorHSV mqtt.Topic = boggart.ComponentName + "/led/+/state/color/hsv"
	WiFiLEDMQTTTopicStateMode     mqtt.Topic = boggart.ComponentName + "/led/+/state/mode"
	WiFiLEDMQTTTopicStateSpeed    mqtt.Topic = boggart.ComponentName + "/led/+/state/speed"
)

type WiFiLED struct {
	statePower int64
	stateMode  uint64
	stateSpeed uint64
	stateColor uint64

	boggart.DeviceBase
	boggart.DeviceSerialNumber
	boggart.DeviceMQTT

	bulb *wifiled.Bulb
}

func NewWiFiLED(bulb *wifiled.Bulb) *WiFiLED {
	device := &WiFiLED{
		bulb: bulb,
	}
	device.Init()
	device.SetDescription("LED WiFi " + bulb.Host())
	device.SetSerialNumber(bulb.Host())

	return device
}

func (d *WiFiLED) Types() []boggart.DeviceType {
	return []boggart.DeviceType{
		boggart.DeviceTypeLED,
	}
}

func (d *WiFiLED) Ping(ctx context.Context) bool {
	_, err := d.bulb.Time(ctx)
	return err == nil
}

func (d *WiFiLED) Tasks() []workers.Task {
	taskUpdater := task.NewFunctionTask(d.taskUpdater)
	taskUpdater.SetRepeats(-1)
	taskUpdater.SetRepeatInterval(WiFiLEDUpdateInterval)
	taskUpdater.SetName("device-led-wifi-updater-" + d.SerialNumber())

	return []workers.Task{
		taskUpdater,
	}
}

func (d *WiFiLED) taskUpdater(ctx context.Context) (interface{}, error) {
	if !d.IsEnabled() {
		return nil, nil
	}

	state, err := d.bulb.State(ctx)
	if err != nil {
		return nil, err
	}

	sn := d.serialNumberEscaped()

	prevPower := atomic.LoadInt64(&d.statePower)
	if prevPower == 0 || (prevPower == 1) != state.Power {
		var mqttValue []byte

		if state.Power {
			atomic.StoreInt64(&d.statePower, 1)
			mqttValue = []byte(`1`)
		} else {
			atomic.StoreInt64(&d.statePower, -1)
			mqttValue = []byte(`0`)
		}

		d.MQTTPublishAsync(ctx, WiFiLEDMQTTTopicStatePower.Format(sn), 0, true, mqttValue)
	}

	currentMode := uint64(state.Mode)
	prevMode := atomic.LoadUint64(&d.stateMode)
	if prevMode != currentMode {
		atomic.StoreUint64(&d.stateMode, currentMode)

		d.MQTTPublishAsync(ctx, WiFiLEDMQTTTopicStateMode.Format(sn), 0, true, currentMode)
	}

	currentSpeed := uint64(state.Speed)
	prevSpeed := atomic.LoadUint64(&d.stateSpeed)
	if prevSpeed != currentSpeed {
		atomic.StoreUint64(&d.stateSpeed, currentSpeed)

		d.MQTTPublishAsync(ctx, WiFiLEDMQTTTopicStateSpeed.Format(sn), 0, true, currentSpeed)
	}

	currentColor := state.Color.Uint64()
	prevColor := atomic.LoadUint64(&d.stateColor)
	if prevColor != currentColor {
		atomic.StoreUint64(&d.stateColor, currentColor)

		// in HEX
		d.MQTTPublishAsync(ctx, WiFiLEDMQTTTopicStateColor.Format(sn), 0, true, state.Color.String())

		// in HSV
		h, s, v := state.Color.HSV()
		d.MQTTPublishAsync(ctx, WiFiLEDMQTTTopicStateColorHSV.Format(sn), 0, true, fmt.Sprintf("%d,%.2f,%.2f", h, s, v))
	}

	return nil, nil
}

func (d *WiFiLED) On(ctx context.Context) error {
	err := d.bulb.PowerOn(ctx)
	if err == nil {
		_, err = d.taskUpdater(ctx)
	}

	return err
}

func (d *WiFiLED) Off(ctx context.Context) error {
	err := d.bulb.PowerOff(ctx)
	if err == nil {
		_, err = d.taskUpdater(ctx)
	}

	return err
}

func (d *WiFiLED) MQTTTopics() []mqtt.Topic {
	return []mqtt.Topic{
		WiFiLEDMQTTTopicPower,
		WiFiLEDMQTTTopicColor,
		WiFiLEDMQTTTopicMode,
		WiFiLEDMQTTTopicSpeed,
		WiFiLEDMQTTTopicStatePower,
		WiFiLEDMQTTTopicStateColor,
		WiFiLEDMQTTTopicStateColorHSV,
		WiFiLEDMQTTTopicStateMode,
		WiFiLEDMQTTTopicStateSpeed,
	}
}

func (d *WiFiLED) MQTTSubscribers() []mqtt.Subscriber {
	sn := d.serialNumberEscaped()

	return []mqtt.Subscriber{
		mqtt.NewSubscriber(WiFiLEDMQTTTopicPower.Format(sn), 0, func(ctx context.Context, client mqtt.Component, message mqtt.Message) {
			if !d.IsEnabled() {
				return
			}

			if bytes.Equal(message.Payload(), []byte(`1`)) {
				d.On(ctx)
			} else {
				d.Off(ctx)
			}
		}),
		mqtt.NewSubscriber(WiFiLEDMQTTTopicColor.Format(sn), 0, func(ctx context.Context, client mqtt.Component, message mqtt.Message) {
			if !d.IsEnabled() {
				return
			}

			color, err := wifiled.ColorFromString(string(message.Payload()))
			if err != nil {
				return
			}

			d.bulb.SetColorPersist(ctx, *color)
		}),
		mqtt.NewSubscriber(WiFiLEDMQTTTopicMode.Format(sn), 0, func(ctx context.Context, client mqtt.Component, message mqtt.Message) {
			if !d.IsEnabled() {
				return
			}

			mode, err := wifiled.ModeFromString(strings.TrimSpace(string(message.Payload())))
			if err != nil {
				return
			}

			state, err := d.bulb.State(ctx)
			if err != nil {
				return
			}

			d.bulb.SetMode(ctx, *mode, state.Speed)
		}),
		mqtt.NewSubscriber(WiFiLEDMQTTTopicSpeed.Format(sn), 0, func(ctx context.Context, client mqtt.Component, message mqtt.Message) {
			if !d.IsEnabled() {
				return
			}

			speed, err := strconv.ParseInt(strings.TrimSpace(string(message.Payload())), 10, 64)
			if err != nil {
				return
			}

			state, err := d.bulb.State(ctx)
			if err != nil {
				return
			}

			d.bulb.SetMode(ctx, state.Mode, uint8(speed))
		}),
	}
}

func (d *WiFiLED) serialNumberEscaped() string {
	sn := strings.Replace(d.SerialNumber(), ":", "-", -1)
	return strings.Replace(sn, ".", "-", -1)
}
