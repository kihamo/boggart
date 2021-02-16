package sp3s

import (
	"github.com/kihamo/boggart/atomic"
	"github.com/kihamo/boggart/components/boggart"
)

type Type struct{}

func (t Type) CreateBind() boggart.Bind {
	return &Bind{
		state: atomic.NewBoolNull(),
		power: atomic.NewFloat32Null(),
	}
}
