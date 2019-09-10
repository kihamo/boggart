package softvideo

import (
	"context"
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/providers/softvideo"
)

type Bind struct {
	boggart.BindBase
	boggart.BindMQTT

	provider        *softvideo.Client
	updaterInterval time.Duration
}

func (b *Bind) Balance(ctx context.Context) (float64, error) {
	return b.provider.Balance(ctx)
}
