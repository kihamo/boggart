package service

import (
	"github.com/kihamo/boggart/components/boggart"
)

type Bind struct {
	boggart.BindBase
	boggart.BindMQTT
	config  *Config
	address string
}

func (b *Bind) Run() error {
	b.UpdateStatus(boggart.BindStatusOnline)
	return nil
}
