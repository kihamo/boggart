package service

import (
	"context"
)

func (b *Bind) ReadinessProbe(ctx context.Context) error {
	return b.Check(ctx)
}
