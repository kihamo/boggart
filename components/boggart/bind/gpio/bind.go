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

func (b *Bind) Run() error {
	if _, ok := b.pin.(gpio.PinIn); ok {
		go func() {
			b.waitForEdge()
		}()
	}

	b.UpdateStatus(boggart.BindStatusOnline)
	return nil
}

func (b *Bind) Mode() Mode {
	return b.mode
}

func (b *Bind) High(ctx context.Context) error {
	if b.Mode() == ModeIn {
		return nil
	}

	if g, ok := b.pin.(gpio.PinOut); ok {
		if err := g.Out(gpio.High); err != nil {
			return err
		}

		if err := b.MQTTPublishAsync(ctx, MQTTPublishTopicPinState.Format(b.pin.Number()), true); err != nil {
			return err
		}
	}

	return nil
}

func (b *Bind) Low(ctx context.Context) error {
	if b.Mode() == ModeIn {
		return nil
	}

	if g, ok := b.pin.(gpio.PinOut); ok {
		if err := g.Out(gpio.Low); err != nil {
			return err
		}

		if err := b.MQTTPublishAsync(ctx, MQTTPublishTopicPinState.Format(b.pin.Number()), false); err != nil {
			return err
		}
	}

	return nil
}

func (b *Bind) Read() bool {
	if b.Mode() == ModeOut {
		return false
	}

	if g, ok := b.pin.(gpio.PinIn); ok {
		return g.Read() == gpio.High
	}

	return false
}

func (b *Bind) waitForEdge() {
	p := b.pin.(gpio.PinIn)
	_ = p.In(gpio.PullNoChange, gpio.BothEdges)
	ctx := context.Background()

	for p.WaitForEdge(-1) {
		// TODO: log
		_ = b.MQTTPublishAsync(ctx, MQTTPublishTopicPinState.Format(b.pin.Number()), b.Read())
	}
}
