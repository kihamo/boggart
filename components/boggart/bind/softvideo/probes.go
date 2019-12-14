package softvideo

import (
	"context"

	"github.com/kihamo/boggart/components/boggart/probes/http"
	"github.com/kihamo/boggart/providers/softvideo"
)

func (b *Bind) ReadinessProbe(ctx context.Context) (err error) {
	return http.Probe(ctx, softvideo.AccountURL, nil)
}
