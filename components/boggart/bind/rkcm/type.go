package rkcm

import (
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/protocols/swagger"
	"github.com/kihamo/boggart/components/boggart/providers/rkcm"
)

type Type struct{}

func (t Type) CreateBind(c interface{}) (boggart.Bind, error) {
	bind := &Bind{
		config: c.(*Config),
	}
	bind.SetSerialNumber(bind.config.Login)

	l := swagger.NewLogger(
		func(message string) {
			bind.Logger().Info(message)
		},
		func(message string) {
			bind.Logger().Debug(message)
		})

	bind.client = rkcm.New(bind.config.Debug, l)

	return bind, nil
}
