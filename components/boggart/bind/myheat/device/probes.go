package device

import (
	"context"

	"github.com/kihamo/boggart/providers/myheat/device/client/state"
)

func (b *Bind) ReadinessProbe(ctx context.Context) error {
	_, err := b.client.State.GetState(state.NewGetStateParamsWithContext(ctx), nil)

	return err
}
