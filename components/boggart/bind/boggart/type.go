package boggart

import (
	"github.com/kihamo/boggart/components/boggart"
)

type Type struct{}

func (t Type) CreateBind(c interface{}) (boggart.Bind, error) {
	config := c.(*Config)

	bind := &Bind{}
	bind.SetSerialNumber(config.Build)

	return bind, nil
}
