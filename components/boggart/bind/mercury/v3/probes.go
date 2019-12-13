package v3

import (
	"context"
)

func (b *Bind) ReadinessProbe(_ context.Context) error {
	return b.provider.ChannelTest()
}
