package softvideo

import (
	"context"

	"github.com/kihamo/boggart/components/boggart/di"
	"github.com/kihamo/boggart/providers/softvideo"
)

type Bind struct {
	di.MQTTBind
	di.WorkersBind

	config   *Config
	provider *softvideo.Client
}

func (b *Bind) Balance(ctx context.Context) (float64, error) {
	return b.provider.Balance(ctx)
}
