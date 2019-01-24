package owntracks

import (
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/atomic"
)

type Type struct{}

func (t Type) CreateBind(c interface{}) (boggart.Bind, error) {
	bind := &Bind{
		config:   c.(*Config),
		lat:      atomic.NewFloat32Null(),
		lon:      atomic.NewFloat32Null(),
		geoHash:  atomic.NewString(),
		conn:     atomic.NewString(),
		acc:      atomic.NewInt32Null(),
		alt:      atomic.NewInt32Null(),
		batt:     atomic.NewFloat32Null(),
		vel:      atomic.NewInt32Null(),
		regions:  make(map[string]Point),
		checkers: make(map[string]*atomic.BoolNull),
	}

	for name, region := range bind.config.Regions {
		if region.Radius <= 0 {
			region.Radius = DefaultPointRadius
		}

		bind.config.Regions[name] = region
		bind.registerRegion(name, region)
	}

	return bind, nil
}
