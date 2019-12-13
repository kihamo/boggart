package v3

import (
	"context"
)

func (b *Bind) ReadinessProbe(ctx context.Context) error {
	return b.provider.ChannelTest()
}
