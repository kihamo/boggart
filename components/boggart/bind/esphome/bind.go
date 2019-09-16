package esphome

import (
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/providers/esphome/native_api"
)

type Bind struct {
	boggart.BindBase

	config   *Config
	provider *native_api.Client
}

func (b *Bind) Close() error {
	return b.provider.Close()
}
