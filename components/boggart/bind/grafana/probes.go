package grafana

import (
	"context"
	"errors"

	"github.com/kihamo/boggart/providers/grafana/client/other"
)

func (b *Bind) ReadinessProbe(ctx context.Context) error {
	response, err := b.provider.Other.Health(other.NewHealthParams(), nil)
	if err != nil {
		return err
	}

	if response.Payload.Database != "ok" {
		return errors.New("databse checker returns " + response.Payload.Database)
	}

	return nil
}
