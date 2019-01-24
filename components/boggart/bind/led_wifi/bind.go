package led_wifi

import (
	"context"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/atomic"
	"github.com/kihamo/boggart/components/boggart/providers/wifiled"
)

type Bind struct {
	boggart.BindBase
	boggart.BindMQTT

	bulb *wifiled.Bulb

	power *atomic.BoolNull
	mode  *atomic.Uint32Null
	speed *atomic.Uint32Null
	color *atomic.Uint32Null
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
