package bind

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/mqtt"
	"github.com/kihamo/shadow/components/tracing"
	"github.com/opentracing/opentracing-go/log"
	"periph.io/x/periph/conn/gpio"
	"periph.io/x/periph/conn/gpio/gpioreg"
	"periph.io/x/periph/conn/pin"
)

type GPIOMode int64

const (
	GPIOMQTTTopicPinState mqtt.Topic = boggart.ComponentName + "/gpio/+"
	GPIOMQTTTopicPinSet   mqtt.Topic = boggart.ComponentName + "/gpio/+/set"

	GPIOModeDefault GPIOMode = iota
	GPIOModeIn
	GPIOModeOut
)

type GPIO struct {
	boggart.DeviceBindBase
	boggart.DeviceBindSerialNumber
	boggart.DeviceBindMQTT

	pin  pin.Pin
	mode GPIOMode
}

func (d GPIO) CreateBind(config map[string]interface{}) (boggart.DeviceBind, error) {
	pinConfig, ok := config["pin"]
	if !ok {
		return nil, errors.New("config option pin isn't set")
	}

	modeConfig, ok := config["mode"]
	if !ok {
		return nil, errors.New("config option mode isn't set")
	}

	g := gpioreg.ByName(fmt.Sprintf("GPIO%d", pinConfig))
	if g == nil {
		return nil, fmt.Errorf("GPIO %d not found", pinConfig)
	}

	var mode GPIOMode
	switch strings.ToLower(modeConfig.(string)) {
	case "in":
		mode = GPIOModeIn
	case "out":
		mode = GPIOModeOut
	default:
		mode = GPIOModeDefault
	}

	device := &GPIO{
		pin:  g,
		mode: mode,
	}

	device.Init()
	device.SetSerialNumber(g.Name())

	if _, ok := g.(gpio.PinIn); ok {
		go func() {
			device.waitForEdge()
		}()
	}

	device.UpdateStatus(boggart.DeviceStatusOnline)

	return device, nil
}

func (d *GPIO) Mode() GPIOMode {
	return d.mode
}

func (d *GPIO) High(ctx context.Context) error {
	if d.Mode() == GPIOModeIn {
		return nil
	}

	if g, ok := d.pin.(gpio.PinOut); ok {
		if err := g.Out(gpio.High); err != nil {
			return err
		}

		d.MQTTPublishAsync(ctx, GPIOMQTTTopicPinState.Format(d.pin.Number()), 2, true, []byte(`1`))
	}

	return nil
}

func (d *GPIO) Low(ctx context.Context) error {
	if d.Mode() == GPIOModeIn {
		return nil
	}

	if g, ok := d.pin.(gpio.PinOut); ok {
		if err := g.Out(gpio.Low); err != nil {
			return err
		}

		d.MQTTPublishAsync(ctx, GPIOMQTTTopicPinState.Format(d.pin.Number()), 2, true, []byte(`0`))
	}

	return nil
}

func (d *GPIO) Read() bool {
	if d.Mode() == GPIOModeOut {
		return false
	}

	if g, ok := d.pin.(gpio.PinIn); ok {
		return g.Read() == gpio.High
	}

	return false
}

func (d *GPIO) waitForEdge() {
	p := d.pin.(gpio.PinIn)
	p.In(gpio.PullNoChange, gpio.BothEdges)
	ctx := context.Background()

	for p.WaitForEdge(-1) {
		var mqttValue []byte

		if d.Read() {
			mqttValue = []byte(`1`)
		} else {
			mqttValue = []byte(`0`)
		}

		d.MQTTPublishAsync(ctx, GPIOMQTTTopicPinState.Format(d.pin.Number()), 2, true, mqttValue)
	}
}

func (d *GPIO) MQTTTopics() []mqtt.Topic {
	return []mqtt.Topic{
		mqtt.Topic(GPIOMQTTTopicPinState.Format(d.pin.Number())),
		mqtt.Topic(GPIOMQTTTopicPinSet.Format(d.pin.Number())),
	}
}

func (d *GPIO) MQTTSubscribers() []mqtt.Subscriber {
	if d.Mode() != GPIOModeOut {
		return nil
	}

	return []mqtt.Subscriber{
		mqtt.NewSubscriber(
			GPIOMQTTTopicPinSet.Format(d.pin.Number()),
			0,
			func(ctx context.Context, client mqtt.Component, message mqtt.Message) {
				span, ctx := tracing.StartSpanFromContext(ctx, "gpio", "set")
				span.LogFields(
					log.String("name", d.pin.Name()),
					log.Int("number", d.pin.Number()))
				defer span.Finish()

				var err error

				if bytes.Equal(message.Payload(), []byte(`1`)) {
					err = d.High(ctx)
					span.LogFields(log.String("out", "high"))
				} else {
					err = d.Low(ctx)
					span.LogFields(log.String("out", "low"))
				}

				if err != nil {
					tracing.SpanError(span, err)
				}
			}),
	}
}
