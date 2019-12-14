package hilink

import (
	"context"

	"github.com/kihamo/boggart/providers/hilink/client/config"
)

func (b *Bind) ReadinessProbe(ctx context.Context) (err error) {
	cfg, err := b.client.Config.GetGlobalConfig(config.NewGetGlobalConfigParamsWithContext(ctx))
	if err != nil {
		return err
	}

	if cfg.Payload.Login == 1 {
		if err := b.client.Login(ctx, b.config.Username, b.config.Password); err != nil {
			return err
		}
	}

	return nil
}
