package octoprint

import (
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/protocols/swagger"
	"github.com/kihamo/boggart/providers/octoprint"
	"github.com/kihamo/shadow/components/dashboard"
)

type Type struct {
	dashboard.Handler
}

func (t Type) CreateBind(c interface{}) (_ boggart.Bind, err error) {
	config := c.(*Config)

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

	bind.provider = octoprint.New(config.Address.Host, config.APIKey, config.Debug, l)

	return bind, nil
}
