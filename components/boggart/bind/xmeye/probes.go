package xmeye

import (
	"context"
)

func (b *Bind) ReadinessProbe(ctx context.Context) error {
	client, err := b.client(ctx)
	if err != nil {
		return err
	}
	defer client.Close()

	_, err = client.SystemInfo(ctx)

	return err
}
