package mikrotik

import (
	"net/url"
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/providers/mikrotik"
)

type Type struct{}

func (t Type) CreateBind(c interface{}) (boggart.Bind, error) {
	config := c.(*Config)

	u, _ := url.Parse(config.Address)
	username := u.User.Username()
	password, _ := u.User.Password()

	device := &Bind{
		provider:         mikrotik.NewClient(u.Host, username, password, time.Second*10),
		host:             u.Host + "-" + u.Port(),
		syslogClient:     config.SyslogClient,
		livenessInterval: config.LivenessInterval,
		livenessTimeout:  config.LivenessTimeout,
		updaterInterval:  config.UpdaterInterval,
	}

	device.Init()

	return device, nil
}
