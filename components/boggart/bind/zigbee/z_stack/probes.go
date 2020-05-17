package z_stack

import (
	"context"
)

func (b *Bind) LivenessProbe(_ context.Context) (err error) {
	_, err = b.client.SysPing()
	return err
}
