package herospeed

import (
	"strconv"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/providers/herospeed"
)

type Type struct {
	boggart.BindTypeWidget
}

func (t Type) CreateBind(c interface{}) (boggart.Bind, error) {
	config := c.(*Config)

	port, _ := strconv.ParseInt(config.Address.Port(), 10, 64)
	password, _ := config.Address.User.Password()

	return &Bind{
		client: herospeed.New(config.Address.Hostname(), port, config.Address.User.Username(), password),
		config: config,
	}, nil
}
