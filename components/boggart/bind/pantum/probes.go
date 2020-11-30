package pantum

import (
	"context"

	"github.com/kihamo/boggart/components/boggart/probes"
)

func (b *Bind) ReadinessProbe(ctx context.Context) error {
	return probes.HTTPProbe(ctx, b.config.Address.String(), nil)
}
