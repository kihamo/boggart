package v3

import (
	"context"

	"github.com/kihamo/boggart/components/boggart/probes"
)

func (b *Bind) LivenessProbe(ctx context.Context) (err error) {
	return probes.ConnErrorProbe(b.ReadinessProbe(ctx))
}

func (b *Bind) ReadinessProbe(_ context.Context) error {
	if provider, e := b.Provider(); e == nil {
		return provider.ChannelTest()
	}

	return nil
}
