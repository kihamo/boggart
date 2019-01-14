package gpio

import (
	"context"

	"github.com/kihamo/boggart/components/boggart"
	"periph.io/x/periph/conn/gpio"
	"periph.io/x/periph/conn/pin"
)

type Mode int64

const (
	ModeDefault Mode = iota
	ModeIn
	ModeOut
)

type Bind struct {
	boggart.BindBase
	boggart.BindMQTT

	pin  pin.Pin
	mode Mode
}

func (d *Bind) Mode() Mode {
	return d.mode
}

func (d *Bind) High(ctx context.Context) error {
	if d.Mode() == ModeIn {
		return nil
	}

	if g, ok := d.pin.(gpio.PinOut); ok {
		if err := g.Out(gpio.High); err != nil {
			return err
		}

		if err := d.MQTTPublishAsync(ctx, MQTTPublishTopicPinState.Format(d.pin.Number()), 2, true, true); err != nil {
			return err
		}
	}

	return nil
}

func (d *Bind) Low(ctx context.Context) error {
	if d.Mode() == ModeIn {
		return nil
	}

	if g, ok := d.pin.(gpio.PinOut); ok {
		if err := g.Out(gpio.Low); err != nil {
			return err
		}

		if err := d.MQTTPublishAsync(ctx, MQTTPublishTopicPinState.Format(d.pin.Number()), 2, true, false); err != nil {
			return err
		}
	}

	return nil
}

func (d *Bind) Read() bool {
	if d.Mode() == ModeOut {
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
		// TODO: log
		_ = d.MQTTPublishAsync(ctx, MQTTPublishTopicPinState.Format(d.pin.Number()), 2, true, d.Read())
	}
}
