package led_wifi

import (
	"context"
)

func (b *Bind) ReadinessProbe(ctx context.Context) (err error) {
	_, err = b.bulb.State(ctx)
	return err
}
