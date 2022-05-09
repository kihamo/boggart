package nut

import (
	"context"
)

func (b *Bind) ReadinessProbe(_ context.Context) (err error) {
	_, err = b.provider.Session().Version()
	return err
}
