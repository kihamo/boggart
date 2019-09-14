package ping

import (
	"github.com/kihamo/boggart/components/boggart"
)

type Bind struct {
	boggart.BindBase
	boggart.BindMQTT
	config *Config
}

func (b *Bind) Run() error {
	b.UpdateStatus(boggart.BindStatusOnline)
	return nil
}
