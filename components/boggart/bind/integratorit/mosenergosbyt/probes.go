package mosenergosbyt

import (
	"context"
)

func (b *Bind) ReadinessProbe(ctx context.Context) (err error) {
	_, err = b.Account(ctx)
	return err
}
