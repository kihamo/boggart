package sp3s

import (
	"context"
)

func (b *Bind) ReadinessProbe(_ context.Context) (err error) {
	_, err = b.State()

	return err
}
