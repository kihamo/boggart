package tvt

import (
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/protocols/swagger"
	"github.com/kihamo/boggart/providers/tvt"
)

type Type struct{}

func (t Type) CreateBind(c interface{}) (boggart.Bind, error) {
	config := c.(*Config)

	password, _ := config.Address.User.Password()

	bind := &Bind{
		config: config,
	}

	l := swagger.NewLogger(
		func(message string) {
			bind.Logger().Info(message)
		},
		func(message string) {
			bind.Logger().Debug(message)
		})

	bind.client = tvt.New(config.Address.Host, config.Address.User.Username(), password, config.Debug, l)

	return bind, nil
}
