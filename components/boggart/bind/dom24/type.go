package dom24

import (
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/protocols/swagger"
	"github.com/kihamo/boggart/providers/dom24"
)

type Type struct{}

func (t Type) CreateBind(c interface{}) (boggart.Bind, error) {
	bind := &Bind{
		config: c.(*Config),
	}

	l := swagger.NewLogger(
		func(message string) {
			bind.Logger().Info(message)
		},
		func(message string) {
			bind.Logger().Debug(message)
		})

	bind.provider = dom24.New(bind.config.Phone, bind.config.Password, bind.config.Debug, l)

	return bind, nil
}
