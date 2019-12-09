package chromecast

import (
	"context"
)

func (b *Bind) ReadinessProbe(ctx context.Context) error {
	return b.Connect(ctx)
}
