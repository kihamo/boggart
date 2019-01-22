package owntracks

import (
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/atomic"
)

type Type struct{}

func (t Type) CreateBind(c interface{}) (boggart.Bind, error) {
	config := c.(*Config)
	wayPointsCheck := make(map[string]*atomic.BoolNull, len(config.WayPoints))

	if len(config.WayPoints) > 0 {
		for n, r := range config.WayPoints {
			if r.Radius <= 0 {
				r.Radius = DefaultPointRadius
			}

			config.WayPoints[n] = r
			wayPointsCheck[n] = atomic.NewBoolNull()
		}
	}

	return &Bind{
		config:         config,
		lat:            atomic.NewFloat64(),
		lon:            atomic.NewFloat64(),
		geoHash:        atomic.NewString(),
		conn:           atomic.NewString(),
		acc:            atomic.NewInt64(),
		alt:            atomic.NewInt64(),
		batt:           atomic.NewFloat64(),
		vel:            atomic.NewInt64(),
		wayPointsCheck: wayPointsCheck,
	}, nil
}
