package esphome

import (
	"context"
)

func (b *Bind) ReadinessProbe(ctx context.Context) (err error) {
	_, err = b.provider.DeviceInfo(ctx)

	return err
}
