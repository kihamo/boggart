package elektroset

import (
	"github.com/kihamo/boggart/atomic"
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/providers/integratorit/elektroset"
)

type Type struct{}

func (t Type) CreateBind(c interface{}) (boggart.Bind, error) {
	config := c.(*Config)

	bind := &Bind{
		config:      config,
		client:      elektroset.New(config.Login, config.Password),
		metersCount: atomic.NewUint32Null(),
	}

	return bind, nil
}
