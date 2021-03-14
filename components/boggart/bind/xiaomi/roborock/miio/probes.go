package miio

import (
	"context"
	"errors"
)

func (b *Bind) ReadinessProbe(ctx context.Context) (err error) {
	_, err = b.device.Timezone(ctx)

	// скорее всего этой кейс, когда счетчик пакетов совпал и был заблокирован на пылесосе
	if errors.Is(err, context.DeadlineExceeded) && b.config().PacketsCounter == 0 {
		b.device.Client().SetPacketsCounter(100)
	}

	return err
}
