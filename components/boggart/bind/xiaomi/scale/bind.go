package scale

import (
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/providers/xiaomi/scale"
)

type Bind struct {
	boggart.BindBase
	boggart.BindMQTT
	config *Config

	provider *scale.Client
}

func (b *Bind) Close() error {
	return b.provider.Close()
}
