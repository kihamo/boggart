package timelapse

import (
	"context"
)

func (b *Bind) ReadinessProbe(ctx context.Context) error {
	return b.Capture(ctx, nil)
}
