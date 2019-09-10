package elektroset

import (
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/providers/integratorit/elektroset"
)

type Type struct{}

func (t Type) CreateBind(c interface{}) (boggart.Bind, error) {
	cfg := c.(*Config)

	bind := &Bind{
		config: cfg,
		client: elektroset.New(cfg.Login, cfg.Password),
	}
	bind.SetSerialNumber(bind.config.Login)

	return bind, nil
}
