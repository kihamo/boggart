package softvideo

import (
	"context"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/providers/softvideo"
)

type Bind struct {
	boggart.BindBase
	boggart.BindMQTT
	config *Config

	provider *softvideo.Client
}

func (b *Bind) Balance(ctx context.Context) (float64, error) {
	return b.provider.Balance(ctx)
}
