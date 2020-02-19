package mikrotik

import (
	"context"
	"errors"
)

func (b *Bind) ReadinessProbe(ctx context.Context) error {
	system, err := b.provider.SystemRouterboard(ctx)
	if err != nil {
		return err
	}

	if system.SerialNumber == "" {
		return errors.New("serial number is empty")
	}

	b.SetSerialNumber(system.SerialNumber)

	return nil
}
