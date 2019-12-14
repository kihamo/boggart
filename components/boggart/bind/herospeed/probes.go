package herospeed

import (
	"context"
)

func (b *Bind) ReadinessProbe(ctx context.Context) error {
	_, err := b.client.Configuration(ctx)

	return err
}
