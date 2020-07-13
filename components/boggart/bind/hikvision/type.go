package hikvision

import (
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/protocols/swagger"
	"github.com/kihamo/boggart/providers/hikvision"
)

type Type struct {
	boggart.BindTypeWidget
}

func (t Type) CreateBind(c interface{}) (boggart.Bind, error) {
	config := c.(*Config)

	password, _ := config.Address.User.Password()

	bind := &Bind{
		alertStreamingHistory: make(map[string]time.Time),
		address:               config.Address.URL,
		config:                config,
	}

	l := swagger.NewLogger(
		func(message string) {
			bind.Logger().Info(message)
		},
		func(message string) {
			bind.Logger().Debug(message)
		})

	bind.client = hikvision.New(config.Address.Host, config.Address.User.Username(), password, config.Debug, l)

	return bind, nil
}
