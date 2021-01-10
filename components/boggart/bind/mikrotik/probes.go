package mikrotik

import (
	"context"

	"github.com/kihamo/boggart/components/boggart/probes"
)

func (b *Bind) LivenessProbe(ctx context.Context) (err error) {
	return probes.ConnErrorProbe(b.ReadinessProbe(ctx))
}

func (b *Bind) ReadinessProbe(ctx context.Context) error {
	_, err := b.provider.SystemClock(ctx)
	return err
}
