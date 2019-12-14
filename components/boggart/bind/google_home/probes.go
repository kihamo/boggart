package google_home

import (
	"context"

	"github.com/kihamo/boggart/providers/google/home/client/info"
)

func (b *Bind) ReadinessProbe(ctx context.Context) (err error) {
	_, err = b.provider.Info.GetEurekaInfo(info.NewGetEurekaInfoParamsWithContext(ctx))
	return err
}
