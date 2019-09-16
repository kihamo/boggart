package mikrotik

import (
	"net/url"
	"time"

	"github.com/kihamo/boggart/atomic"
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/providers/mikrotik"
)

type Type struct{}

func (t Type) CreateBind(c interface{}) (boggart.Bind, error) {
	config := c.(*Config)

	u, err := url.Parse(config.Address)
	if err != nil {
		return nil, err
	}

	username := u.User.Username()
	password, _ := u.User.Password()

	bind := &Bind{
		config:            config,
		address:           u,
		provider:          mikrotik.NewClient(u.Host, username, password, time.Second*10),
		serialNumberLock:  make(chan struct{}),
		serialNumberReady: atomic.NewBool(),
	}

	bind.clientWiFi = NewPreloadMap()
	bind.clientVPN = NewPreloadMap()

	return bind, nil
}
