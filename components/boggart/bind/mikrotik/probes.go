package mikrotik

import (
	"context"
	"errors"

	"github.com/kihamo/boggart/components/boggart/probes"
)

func (b *Bind) LivenessProbe(ctx context.Context) (err error) {
	return probes.ConnErrorProbe(b.ReadinessProbe(ctx))
}

func (b *Bind) ReadinessProbe(ctx context.Context) error {
	system, err := b.provider.SystemRouterBoard(ctx)
	if err != nil {
		return err
	}

	if system.SerialNumber == "" {
		return errors.New("serial number is empty")
	}

	b.SetSerialNumber(system.SerialNumber)

	return nil
}
