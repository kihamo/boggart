package sun

import (
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/atomic"
)

type Type struct{}

func (t Type) CreateBind(c interface{}) (boggart.Bind, error) {
	config := c.(*Config)

	bind := &Bind{
		config:        config,
		riseStart:     atomic.NewTimeNull(),
		riseEnd:       atomic.NewTimeNull(),
		riseDuration:  atomic.NewDuration(),
		setStart:      atomic.NewTimeNull(),
		setEnd:        atomic.NewTimeNull(),
		setDuration:   atomic.NewDuration(),
		nightStart:    atomic.NewTimeNull(),
		nightEnd:      atomic.NewTimeNull(),
		nightDuration: atomic.NewDuration(),
		nadir:         atomic.NewTimeNull(),
		solarNoon:     atomic.NewTimeNull(),
	}
	bind.SetSerialNumber(config.Name)

	return bind, nil
}
