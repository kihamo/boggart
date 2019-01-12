package mikrotik

import (
	"net/url"
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/providers/mikrotik"
)

type Type struct{}

func (t Type) CreateBind(c interface{}) (boggart.DeviceBind, error) {
	config := c.(*Config)

	u, _ := url.Parse(config.Address)
	username := u.User.Username()
	password, _ := u.User.Password()

	device := &Bind{
		provider:     mikrotik.NewClient(u.Host, username, password, time.Second*10),
		host:         u.Host + "-" + u.Port(),
		syslogClient: config.SyslogClient,
	}

	if device.syslogClient == "" {
		device.syslogClient = u.Hostname() + ":514"
	}

	device.Init()

	return device, nil
}
