package rm

import (
	"context"

	"github.com/kihamo/boggart/components/boggart/probes"
)

func (b *Bind) ReadinessProbe(ctx context.Context) error {
	return probes.PingProbe(ctx, b.config().Host)
}
