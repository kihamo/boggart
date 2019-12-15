package miio

import (
	"context"
)

func (b *Bind) ReadinessProbe(ctx context.Context) (err error) {
	_, err = b.device.SerialNumber(ctx)

	return err
}
