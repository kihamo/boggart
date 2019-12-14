package pulsar

import (
	"context"
)

func (b *Bind) ReadinessProbe(_ context.Context) (err error) {
	_, err = b.provider.Version()
	return err
}
