package led_wifi

import (
	"context"
)

func (b *Bind) ReadinessProbe(ctx context.Context) error {
	return b.State(ctx)
}
