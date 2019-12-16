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
		if user := b.config.Address.User; user != nil {
			password, _ := user.Password()

			if err := b.client.Login(ctx, user.Username(), password); err != nil {
				return err
			}
		}
	}

	return nil
}
