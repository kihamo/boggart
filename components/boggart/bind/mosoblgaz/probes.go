package mosoblgaz

import (
	"context"
)

func (b *Bind) ReadinessProbe(ctx context.Context) error {
	return b.taskUpdaterHandler(ctx)
}
