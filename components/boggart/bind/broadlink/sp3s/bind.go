package sp3s

import (
	"context"
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/atomic"
	"github.com/kihamo/boggart/providers/broadlink"
)

type Bind struct {
	boggart.BindBase
	boggart.BindMQTT

	state *atomic.BoolNull
	power *atomic.Float32Null

	provider        *broadlink.SP3S
	updaterInterval time.Duration
}

func (b *Bind) State() (bool, error) {
	return b.provider.State()
}

func (b *Bind) On(ctx context.Context) error {
	err := b.provider.On()
	if err == nil {
		_, err = b.taskUpdater(ctx)
	}

	return err
}

func (b *Bind) Off(ctx context.Context) error {
	err := b.provider.Off()
	if err == nil {
		_, err = b.taskUpdater(ctx)
	}

	return err
}

func (b *Bind) Power() (float64, error) {
	return b.provider.Power()
}
