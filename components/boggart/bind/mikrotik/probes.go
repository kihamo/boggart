package mikrotik

import (
	"context"
)

func (b *Bind) ReadinessProbe(ctx context.Context) (err error) {
	_, err = b.provider.SystemRouterboard(ctx)

	return err
}
