package softvideo

import (
	"context"
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/providers/softvideo"
)

type Bind struct {
	lastValue int64

	boggart.DeviceBindBase
	boggart.DeviceBindMQTT

	provider        *softvideo.Client
	updaterInterval time.Duration
}

func (b *Bind) Balance(ctx context.Context) (float64, error) {
	return b.provider.Balance(ctx)
}
