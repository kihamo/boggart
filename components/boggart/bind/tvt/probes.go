package tvt

import (
	"context"

	"github.com/kihamo/boggart/providers/tvt/client/information"
)

func (b *Bind) ReadinessProbe(ctx context.Context) error {
	_, err := b.client.Information.GetBasicConfig(information.NewGetBasicConfigParamsWithContext(ctx), nil)

	return err
}
