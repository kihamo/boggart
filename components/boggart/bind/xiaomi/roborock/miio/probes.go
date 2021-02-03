package miio

import (
	"context"
)

func (b *Bind) ReadinessProbe(ctx context.Context) (err error) {
	_, err = b.device.Timezone(ctx)

	return err
}
