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
	di.MQTTBind
	di.LoggerBind
	di.WidgetBind

	config *Config

	pin  pin.Pin
	mode Mode
	out  *atomic.Bool
}

func (b *Bind) Run() error {
	err := b.MQTT().PublishAsync(context.Background(), b.config.TopicPinState, b.Read())
	if err != nil {
		return err
	}

	if _, ok := b.pin.(gpio.PinIn); ok {
		go func() {
			b.waitForEdge()
		}()
	}

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

		if err := b.MQTT().PublishAsync(ctx, b.config.TopicPinState, true); err != nil {
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

		if err := b.MQTT().PublishAsync(ctx, b.config.TopicPinState, false); err != nil {
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
	_ = p.In(gpio.PullNoChange, gpio.BothEdges)
	ctx := context.Background()

	for p.WaitForEdge(-1) {
		v := b.Read()

		if v {
			b.Logger().Debugf("Pin %s edge high", p.String())
		} else {
			b.Logger().Debugf("Pin %s edge log", p.String())
		}

		if err := b.MQTT().PublishAsync(ctx, b.config.TopicPinState, v); err != nil {
			b.Logger().Errorf("Publish to %s topic failed with error %v", b.config.TopicPinState, err)
		}
	}
}
