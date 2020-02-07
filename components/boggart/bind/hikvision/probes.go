package hikvision

import (
	"context"

	"github.com/kihamo/boggart/providers/hikvision/client/system"
)

func (b *Bind) ReadinessProbe(ctx context.Context) error {
	_, err := b.client.System.GetStatus(system.NewGetStatusParamsWithContext(ctx), nil)

	return err
}
