package xmeye

import (
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/providers/xmeye"
)

type Type struct {
	boggart.BindTypeWidget
}

func (t Type) CreateBind(c interface{}) (boggart.Bind, error) {
	config := c.(*Config)

	password, _ := config.Address.User.Password()

	provider, err := xmeye.New(config.Address.Host, config.Address.User.Username(), password)
	if err != nil {
		return nil, err
	}

	bind := &Bind{
		config: config,
		client: provider,
	}

	return bind, nil
}
