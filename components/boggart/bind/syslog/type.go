package syslog

import (
	"net"
	"strconv"

	"github.com/kihamo/boggart/components/boggart"
)

type Type struct {
	boggart.BindTypeWidget
}

func (t Type) CreateBind(c interface{}) (boggart.Bind, error) {
	cfg := c.(*Config)

	bind := &Bind{
		config: cfg,
		addr:   net.JoinHostPort(cfg.Hostname, strconv.FormatInt(cfg.Port, 10)),
	}

	return bind, nil
}
