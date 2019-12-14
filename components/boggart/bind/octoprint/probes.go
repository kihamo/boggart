package octoprint

import (
	"context"

	"github.com/kihamo/boggart/providers/octoprint/client/version"
)

func (b *Bind) ReadinessProbe(ctx context.Context) error {
	_, err := b.provider.Version.GetVersion(version.NewGetVersionParamsWithContext(ctx), nil)

	return err
}
