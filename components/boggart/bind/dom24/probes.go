package dom24

import (
	"context"

	"github.com/kihamo/boggart/providers/dom24/client/config"
)

func (b *Bind) ReadinessProbe(ctx context.Context) (err error) {
	_, err = b.provider.Config.MobileAppSettings(config.NewMobileAppSettingsParamsWithContext(ctx))

	return err
}
