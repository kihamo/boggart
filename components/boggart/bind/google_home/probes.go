package home

import (
	"context"
)

func (b *Bind) ReadinessProbe(ctx context.Context) (err error) {
	_, err = b.GetMAC(ctx)
	return err
}
