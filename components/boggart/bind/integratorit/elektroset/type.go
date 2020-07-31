package elektroset

import (
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/providers/integratorit/elektroset"
)

type Type struct {
	boggart.BindTypeWidget
}

func (t Type) CreateBind(c interface{}) (boggart.Bind, error) {
	config := c.(*Config)

	bind := &Bind{
		config: config,
		client: elektroset.New(config.Login, config.Password),
	}

	return bind, nil
}
