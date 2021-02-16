package zstack

import (
	"github.com/kihamo/boggart/atomic"
	"github.com/kihamo/boggart/components/boggart"
)

type Type struct{}

func (t Type) CreateBind() boggart.Bind {
	return &Bind{
		disconnected: atomic.NewBoolNull(),
		onceClient:   &atomic.Once{},
	}
}
