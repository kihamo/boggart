package elektroset

import (
	"context"

	"github.com/kihamo/boggart/components/boggart/probes"
	"github.com/kihamo/boggart/providers/integratorit/elektroset"
)

func (b *Bind) ReadinessProbe(ctx context.Context) (err error) {
	return probes.HTTPProbe(ctx, elektroset.BaseURL, nil)
}
