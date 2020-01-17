package mosenergosbyt

import (
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/providers/integratorit/mosenergosbyt"
)

type Type struct{}

func (t Type) CreateBind(c interface{}) (boggart.Bind, error) {
	config := c.(*Config)

	bind := &Bind{
		config: config,
		client: mosenergosbyt.New(config.Login, config.Password),
	}

	return bind, nil
}
