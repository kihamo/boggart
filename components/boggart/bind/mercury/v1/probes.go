package v1

import (
	"context"
)

func (b *Bind) ReadinessProbe(_ context.Context) (err error) {
	_, _, err = b.provider.Version()
	return err
}
