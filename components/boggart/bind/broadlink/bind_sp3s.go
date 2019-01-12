package broadlink

import (
	"context"
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/providers/broadlink"
)

type BindSP3S struct {
	state int64
	power int64

	boggart.DeviceBindBase
	boggart.DeviceBindMQTT

	provider        *broadlink.SP3S
	updaterInterval time.Duration
}

func (b *BindSP3S) State() (bool, error) {
	return b.provider.State()
}

func (b *BindSP3S) On(ctx context.Context) error {
	err := b.provider.On()
	if err == nil {
		_, err = b.taskStateUpdater(ctx)
	}

	return err
}

func (b *BindSP3S) Off(ctx context.Context) error {
	err := b.provider.Off()
	if err == nil {
		_, err = b.taskStateUpdater(ctx)
	}

	return err
}

func (b *BindSP3S) Power() (float64, error) {
	return b.provider.Power()
}
