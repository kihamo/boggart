package service

import (
	"net"
	"strconv"

	"github.com/kihamo/boggart/components/boggart"
)

type Type struct{}

func (t Type) CreateBind(c interface{}) (boggart.Bind, error) {
	config := c.(*Config)
	addr := net.JoinHostPort(config.Hostname, strconv.Itoa(config.Port))

	config.TopicOnline = config.TopicOnline.Format(addr)
	config.TopicLatency = config.TopicLatency.Format(addr)
	config.TopicCheck = config.TopicCheck.Format(addr)

	bind := &Bind{
		config:  config,
		address: addr,
	}

	return bind, nil
}
