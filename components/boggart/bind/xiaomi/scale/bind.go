package scale

import (
	"context"

	"github.com/kihamo/boggart/atomic"
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/providers/xiaomi/scale"
)

type Bind struct {
	boggart.BindBase
	boggart.BindMQTT
	config *Config

	provider *scale.Client

	sex    *atomic.BoolNull
	height *atomic.Uint32Null
	age    *atomic.Uint32Null
}

func (b *Bind) Run() error {
	b.updateProfile(context.Background())
	return nil
}

func (b *Bind) Close() error {
	return b.provider.Close()
}
