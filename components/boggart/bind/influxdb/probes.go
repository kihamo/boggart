package influxdb

import (
	"context"
)

func (b *Bind) ReadinessProbe(ctx context.Context) (err error) {
	_, err = b.client.Health(ctx)
	return err
}
