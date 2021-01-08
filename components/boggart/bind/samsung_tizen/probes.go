package tizen

import (
	"context"
)

func (b *Bind) ReadinessProbe(ctx context.Context) (err error) {
	_, err = b.client.Device(ctx)

	return err
}
