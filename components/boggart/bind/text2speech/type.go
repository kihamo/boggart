package text2speech

import (
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/protocols/swagger"
	"github.com/kihamo/boggart/providers/yandex_speechkit_cloud"
)

type Type struct {
	boggart.BindTypeWidget
}

func (t Type) CreateBind(c interface{}) (boggart.Bind, error) {
	cfg := c.(*Config)

	bind := &Bind{
		config: cfg,
	}

	l := swagger.NewLogger(
		func(message string) {
			bind.Logger().Info(message)
		},
		func(message string) {
			bind.Logger().Debug(message)
		})

	bind.provider = speechkit.New(cfg.Debug, l)

	return bind, nil
}
