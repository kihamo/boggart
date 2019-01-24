package softvideo

import (
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/atomic"
	"github.com/kihamo/boggart/components/boggart/providers/softvideo"
)

type Type struct{}

func (t Type) CreateBind(c interface{}) (boggart.Bind, error) {
	config := c.(*Config)

	bind := &Bind{
		provider:        softvideo.NewClient(config.Login, config.Password, config.Debug),
		balance:         atomic.NewFloat32Null(),
		updaterInterval: config.UpdaterInterval,
	}
	bind.SetSerialNumber(config.Login)

	return bind, nil
}
