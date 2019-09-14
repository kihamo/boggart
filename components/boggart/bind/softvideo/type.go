package softvideo

import (
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/providers/softvideo"
)

type Type struct{}

func (t Type) CreateBind(c interface{}) (boggart.Bind, error) {
	config := c.(*Config)

	config.TopicBalance = config.TopicBalance.Format(config.Login)

	bind := &Bind{
		config:   config,
		provider: softvideo.NewClient(config.Login, config.Password, config.Debug),
	}
	bind.SetSerialNumber(config.Login)

	return bind, nil
}
