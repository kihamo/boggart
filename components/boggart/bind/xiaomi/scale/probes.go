package scale

import (
	"context"
	"errors"
)

func (b *Bind) LivenessProbe(ctx context.Context) error {
	if b.disconnected.IsTrue() {
		return errors.New("disconnected")
	}

	return nil
}
