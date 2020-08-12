package timelapse

import (
	"github.com/kihamo/boggart/components/boggart"
)

type Type struct{}

func (t Type) CreateBind(c interface{}) (boggart.Bind, error) {
	bind := &Bind{
		config: c.(*Config),
	}

	return bind, nil
}
