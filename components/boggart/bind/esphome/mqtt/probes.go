package mqtt

import (
	"context"
	"errors"
)

func (b *Bind) ReadinessProbe(ctx context.Context) (err error) {
	if b.status.IsTrue() {
		return nil
	}

	return errors.New("status isn't online")
}
