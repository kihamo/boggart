package devices

import (
	"bytes"
	"context"
	"fmt"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/mqtt"
	"github.com/kihamo/shadow/components/tracing"
	"github.com/opentracing/opentracing-go/log"
	"periph.io/x/periph/conn/gpio"
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

type GPIOPin struct {
	boggart.DeviceBase
	boggart.DeviceSerialNumber
	boggart.DeviceMQTT

	pin  pin.Pin
	mode GPIOMode
}

func NewGPIOPin(p pin.Pin, m GPIOMode) *GPIOPin {
	device := &GPIOPin{
		pin:  p,
		mode: m,
	}

	device.Init()
	device.SetSerialNumber(p.Name())
	device.SetDescription(fmt.Sprintf("%s %s", p.Name(), p.Function()))

	if _, ok := p.(gpio.PinIn); ok {
		go func() {
			device.waitForEdge()
		}()
	}

	device.UpdateStatus(boggart.DeviceStatusOnline)

	return device
}

func (d *GPIOPin) Types() []boggart.DeviceType {
	return []boggart.DeviceType{
		boggart.DeviceTypeGPIO,
	}
}

func (d *GPIOPin) Ping(_ context.Context) bool {
	return true
}

func (d *GPIOPin) Mode() GPIOMode {
	return d.mode
}

func (d *GPIOPin) High(ctx context.Context) error {
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

func (d *GPIOPin) Low(ctx context.Context) error {
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

func (d *GPIOPin) Read() bool {
	if d.Mode() == GPIOModeOut {
		return false
	}

	if g, ok := d.pin.(gpio.PinIn); ok {
		return g.Read() == gpio.High
	}

	return false
}

func (d *GPIOPin) waitForEdge() {
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

func (d *GPIOPin) MQTTTopics() []mqtt.Topic {
	return []mqtt.Topic{
		GPIOMQTTTopicPinState,
		GPIOMQTTTopicPinSet,
	}
}

func (d *GPIOPin) MQTTSubscribers() []mqtt.Subscriber {
	if d.Mode() != GPIOModeOut {
		return nil
	}

	return []mqtt.Subscriber{
		mqtt.NewSubscriber(
			GPIOMQTTTopicPinSet.Format(d.pin.Number()),
			0,
			func(ctx context.Context, client mqtt.Component, message mqtt.Message) {
				if !d.IsEnabled() {
					return
				}

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
