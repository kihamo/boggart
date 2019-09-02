package ping

import (
	"github.com/kihamo/boggart/components/boggart"
)

type Type struct{}

func (t Type) CreateBind(c interface{}) (boggart.Bind, error) {
	config := c.(*Config)

	bind := &Bind{
		hostname:        config.Hostname,
		retry:           config.Retry,
		timeout:         config.Timeout,
		updaterInterval: config.UpdaterInterval,
	}

	return bind, nil
}
