package mosenergosbyt

import (
	"context"

	"github.com/kihamo/boggart/components/boggart/probes/http"
	"github.com/kihamo/boggart/providers/integratorit/mosenergosbyt"
)

func (b *Bind) ReadinessProbe(ctx context.Context) (err error) {
	return http.Probe(ctx, mosenergosbyt.BaseURL, nil)
}
