package miio

import (
	"context"
)

func (b *Bind) ReadinessProbe(ctx context.Context) error {
	sn, err := b.device.SerialNumber(ctx)
	if err == nil {
		b.Meta().SetSerialNumber(sn)
	}

	return err
}
