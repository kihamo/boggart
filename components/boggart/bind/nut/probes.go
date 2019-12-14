package nut

import (
	"context"
)

func (b *Bind) ReadinessProbe(ctx context.Context) error {
	client, err := b.connect()
	if err != nil {
		return err
	}
	defer func() {
		_, _ = client.Disconnect()
	}()

	_, err = client.GetVersion()
	return err
}
