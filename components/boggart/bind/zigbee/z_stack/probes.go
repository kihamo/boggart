package z_stack

import (
	"context"
	"errors"

	"github.com/kihamo/boggart/providers/zigbee/z_stack"
)

func (b *Bind) LivenessProbe(_ context.Context) error {
	if b.disconnected.IsTrue() {
		return errors.New("disconnected")
	}

	return nil
}

func (b *Bind) ReadinessProbe(ctx context.Context) (err error) {
	if b.disconnected.IsNil() {
		var client *z_stack.Client
		client, err = b.getClient()

		if err != nil {
			_, err = client.SysPing(ctx)
		}

		if err != nil {
			b.disconnected.True()
		}
	} else if b.disconnected.IsTrue() {
		err = errors.New("disconnected")
	}

	return err
}
