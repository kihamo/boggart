package led_wifi

import (
	"context"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/providers/wifiled"
)

type Bind struct {
	boggart.BindBase
	boggart.BindMQTT

	bulb *wifiled.Bulb
}

func (b *Bind) On(ctx context.Context) error {
	err := b.bulb.PowerOn(ctx)
	if err == nil {
		_, err = b.taskUpdater(ctx)
	}

	return err
}

func (b *Bind) Off(ctx context.Context) error {
	err := b.bulb.PowerOff(ctx)
	if err == nil {
		_, err = b.taskUpdater(ctx)
	}

	return err
}
