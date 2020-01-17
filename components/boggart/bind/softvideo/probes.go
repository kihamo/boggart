package softvideo

import (
	"context"

	"github.com/kihamo/boggart/components/boggart/probes"
	"github.com/kihamo/boggart/providers/softvideo"
)

func (b *Bind) ReadinessProbe(ctx context.Context) (err error) {
	return probes.HTTPProbe(ctx, softvideo.AccountURL, nil)
}
