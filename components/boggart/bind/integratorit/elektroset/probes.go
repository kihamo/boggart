package elektroset

import (
	"context"

	"github.com/kihamo/boggart/components/boggart/probes/http"
	"github.com/kihamo/boggart/providers/integratorit/elektroset"
)

func (b *Bind) ReadinessProbe(ctx context.Context) (err error) {
	return http.Probe(ctx, elektroset.BaseURL, nil)
}
