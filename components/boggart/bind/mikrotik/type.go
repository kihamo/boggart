package mikrotik

import (
	"github.com/kihamo/boggart/atomic"
	"github.com/kihamo/boggart/components/boggart"
)

type Type struct{}

func (t Type) CreateBind() boggart.Bind {
	return &Bind{
		connectionsZombieKiller: &atomic.Once{},
		connectionsFirstLoad: map[string]*atomic.Once{
			InterfaceWireless:   {},
			InterfaceL2TPServer: {},
		},
	}
}
