package gpio

import (
	"context"

	"github.com/kihamo/boggart/components/boggart"
	"periph.io/x/periph/conn/gpio"
	"periph.io/x/periph/conn/pin"
)

type GPIOMode int64

const (
	GPIOModeDefault GPIOMode = iota
	GPIOModeIn
	GPIOModeOut
)

type Bind struct {
	boggart.DeviceBindBase
	boggart.DeviceBindMQTT

	pin  pin.Pin
	mode GPIOMode
}

func (d *Bind) Mode() GPIOMode {
	return d.mode
}

func (d *Bind) High(ctx context.Context) error {
	if d.Mode() == GPIOModeIn {
		return nil
	}

	if g, ok := d.pin.(gpio.PinOut); ok {
		if err := g.Out(gpio.High); err != nil {
			return err
		}

		d.MQTTPublishAsync(ctx, TopicPinState.Format(d.pin.Number()), 2, true, []byte(`1`))
	}

	return nil
}

func (d *Bind) Low(ctx context.Context) error {
	if d.Mode() == GPIOModeIn {
		return nil
	}

	if g, ok := d.pin.(gpio.PinOut); ok {
		if err := g.Out(gpio.Low); err != nil {
			return err
		}

		d.MQTTPublishAsync(ctx, TopicPinState.Format(d.pin.Number()), 2, true, []byte(`0`))
	}

	return nil
}

func (d *Bind) Read() bool {
	if d.Mode() == GPIOModeOut {
		return false
	}

	if g, ok := d.pin.(gpio.PinIn); ok {
		return g.Read() == gpio.High
	}

	return false
}

func (d *Bind) waitForEdge() {
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

		d.MQTTPublishAsync(ctx, TopicPinState.Format(d.pin.Number()), 2, true, mqttValue)
	}
}
