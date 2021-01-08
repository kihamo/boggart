package nut

import (
	"github.com/kihamo/boggart/atomic"
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/providers/nut"
)

type Type struct{}

func (t Type) CreateBind(c interface{}) (boggart.Bind, error) {
	cfg := c.(*Config)

	var username, password string

	if cfg.Address.User != nil {
		username = cfg.Address.User.Username()
		password, _ = cfg.Address.User.Password()
	}

	bind := &Bind{
		config:          cfg,
		provider:        nut.New(cfg.Address.Host, username, password),
		updaterInterval: atomic.NewDurationDefault(cfg.UpdaterInterval),
	}

	return bind, nil
}
