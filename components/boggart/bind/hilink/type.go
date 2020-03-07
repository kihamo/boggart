package hilink

import (
	"github.com/kihamo/boggart/atomic"
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/protocols/swagger"
	"github.com/kihamo/boggart/providers/hilink"
)

type Type struct {
	boggart.BindTypeWidget
}

func (t Type) CreateBind(c interface{}) (boggart.Bind, error) {
	config := c.(*Config)

	bind := &Bind{
		config:                    config,
		operator:                  atomic.NewString(),
		limitInternetTrafficIndex: atomic.NewInt64(),
		simStatus:                 atomic.NewUint32(),
		ussdEnabled:               atomic.NewBool(),
		smsEnabled:                atomic.NewBool(),
	}

	l := swagger.NewLogger(
		func(message string) {
			bind.Logger().Info(message)
		},
		func(message string) {
			bind.Logger().Debug(message)
		})

	bind.client = hilink.New(config.Address.Host, config.Debug, l)

	return bind, nil
}
