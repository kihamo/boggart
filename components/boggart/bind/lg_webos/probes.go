package webos

import (
	"context"
	"errors"
)

func (b *Bind) LivenessProbe(ctx context.Context) error {
	if client := b.Client(); client != nil {
		_, err := client.SystemGetSystemInfo()
		return err
	}

	return nil
}

func (b *Bind) ReadinessProbe(ctx context.Context) (err error) {
	client := b.Client()

	if client == nil {
		err = b.initClient()
	} else if !b.power.IsNil() && b.power.IsFalse() {
		err = errors.New("application ID is empty")
	}

	return err
}
