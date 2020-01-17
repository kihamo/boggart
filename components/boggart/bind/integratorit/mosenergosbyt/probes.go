package mosenergosbyt

import (
	"context"

	"github.com/kihamo/boggart/components/boggart/probes"
	"github.com/kihamo/boggart/providers/integratorit/mosenergosbyt"
)

func (b *Bind) ReadinessProbe(ctx context.Context) (err error) {
	return probes.HTTPProbe(ctx, mosenergosbyt.BaseURL, nil)
}
