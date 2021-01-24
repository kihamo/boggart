package ds18b20

import (
	"context"

	"github.com/yryz/ds18b20"
)

func (b *Bind) ReadinessProbe(_ context.Context) error {
	_, err := ds18b20.Sensors()

	return err
}
