package smcenter

import (
	"context"

	"github.com/kihamo/boggart/providers/smcenter/client/config"
)

func (b *Bind) ReadinessProbe(ctx context.Context) (err error) {
	_, err = b.provider.Config.MobileAppSettings(config.NewMobileAppSettingsParamsWithContext(ctx))

	return err
}
