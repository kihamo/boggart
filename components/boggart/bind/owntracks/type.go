package owntracks

import (
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/atomic"
)

type Type struct{}

func (t Type) CreateBind(c interface{}) (boggart.Bind, error) {
	device := &Bind{
		config:   c.(*Config),
		lat:      atomic.NewFloat64(),
		lon:      atomic.NewFloat64(),
		geoHash:  atomic.NewString(),
		conn:     atomic.NewString(),
		acc:      atomic.NewInt64(),
		alt:      atomic.NewInt64(),
		batt:     atomic.NewFloat64(),
		vel:      atomic.NewInt64(),
		regions:  make(map[string]Point),
		checkers: make(map[string]*atomic.BoolNull),
	}

	for name, region := range device.config.Regions {
		if region.Radius <= 0 {
			region.Radius = DefaultPointRadius
		}

		device.config.Regions[name] = region
		device.registerRegion(name, region)
	}

	return device, nil
}
