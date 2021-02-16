package gpio

import (
	"context"
	"errors"

	"github.com/kihamo/boggart/atomic"
	"github.com/kihamo/boggart/components/boggart/di"
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
	di.ConfigBind
	di.LoggerBind
	di.MetaBind
	di.MQTTBind
	di.WidgetBind

	config *Config
	done   chan struct{}

	pin  pin.Pin
	mode Mode
	out  *atomic.Bool
}

func (b *Bind) Run() error {
	err := b.publishState(context.Background(), b.Read())
	if err != nil {
		return err
	}

	b.done = make(chan struct{})

	if p, ok := b.pin.(gpio.PinIn); ok {
		err = p.In(gpio.PullNoChange, gpio.BothEdges)
		if err != nil {
			return err
		}

		go b.waitForEdge()
	}

	return nil
}

func (b *Bind) Close() error {
	close(b.done)

	return nil
}

func (b *Bind) Mode() Mode {
	return b.mode
}

func (b *Bind) High(ctx context.Context) error {
	if b.Mode() == ModeIn {
		return errors.New("is read only")
	}

	if g, ok := b.pin.(gpio.PinOut); ok {
		if err := g.Out(gpio.High); err != nil {
			return err
		}

		b.out.True()

		if err := b.publishState(ctx, true); err != nil {
			return err
		}
	}

	return nil
}

func (b *Bind) Low(ctx context.Context) error {
	if b.Mode() == ModeIn {
		return errors.New("is read only")
	}

	if g, ok := b.pin.(gpio.PinOut); ok {
		if err := g.Out(gpio.Low); err != nil {
			return err
		}

		b.out.False()

		if err := b.publishState(ctx, false); err != nil {
			return err
		}
	}

	return nil
}

func (b *Bind) Read() bool {
	if b.Mode() == ModeOut {
		return b.out.IsTrue()
	}

	if g, ok := b.pin.(gpio.PinIn); ok {
		return g.Read() == gpio.High
	}

	return false
}

func (b *Bind) waitForEdge() {
	p := b.pin.(gpio.PinIn)
	ctx := context.Background()

	for {
		select {
		case <-b.done:
			return

		default:
			if p.WaitForEdge(-1) {
				v := b.Read()

				if v {
					b.Logger().Debugf("Pin %s edge high", p.String())
				} else {
					b.Logger().Debugf("Pin %s edge low", p.String())
				}

				if err := b.publishState(ctx, v); err != nil {
					b.Logger().Errorf("Publish to %s topic failed with error %v", b.config.TopicPinState, err)
				}
			}
		}
	}
}
