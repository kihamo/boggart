package softvideo

import (
	"context"
)

func (b *Bind) ReadinessProbe(ctx context.Context) (err error) {
	_, _, err = b.Balance(ctx)
	return err
}
