package ds18b20

import (
	"context"
	"errors"

	"github.com/yryz/ds18b20"
)

func (b *Bind) ReadinessProbe(_ context.Context) error {
	devices, err := ds18b20.Sensors()
	if err != nil {
		return err
	}

	sn := b.SerialNumber()

	for _, device := range devices {
		if device == sn {
			return nil
		}
	}

	return errors.New("device with ID " + sn + " not found")
}
