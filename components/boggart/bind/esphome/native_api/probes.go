package nativeapi

import (
	"context"
)

func (b *Bind) ReadinessProbe(ctx context.Context) (err error) {
	_, err = b.provider.Ping(ctx)

	return err
}
