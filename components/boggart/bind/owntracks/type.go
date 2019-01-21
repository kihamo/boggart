package owntracks

import (
	"github.com/kihamo/boggart/components/boggart"
)

type Type struct{}

func (t Type) CreateBind(c interface{}) (boggart.Bind, error) {
	config := c.(*Config)

	return &Bind{
		user:   config.User,
		device: config.Device,
	}, nil
}
