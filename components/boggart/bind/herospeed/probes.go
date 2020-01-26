package herospeed

import (
	"context"
)

func (b *Bind) ReadinessProbe(ctx context.Context) error {
	_, err := b.GetSerialNumber(ctx)

	return err
}
