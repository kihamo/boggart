package hilink

import (
	"github.com/kihamo/boggart/atomic"
	"github.com/kihamo/boggart/components/boggart"
)

type Type struct{}

func (t Type) CreateBind() boggart.Bind {
	return &Bind{
		operator:                  atomic.NewString(),
		limitInternetTrafficIndex: atomic.NewInt64(),
		simStatus:                 atomic.NewUint32(),
	}
}
