package ledwifi

import (
	"context"
)

func (b *Bind) ReadinessProbe(ctx context.Context) error {
	_, err := b.bulb.State(ctx)

	return err
}
