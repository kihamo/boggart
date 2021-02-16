package mqtt

import (
	"github.com/kihamo/boggart/atomic"
	"github.com/kihamo/boggart/components/boggart"
)

type Type struct{}

func (t Type) CreateBind() boggart.Bind {
	return &Bind{
		ip:                     atomic.NewValue(),
		ipSubscriber:           atomic.NewBool(),
		connectivitySubscriber: atomic.NewBool(),
	}
}
