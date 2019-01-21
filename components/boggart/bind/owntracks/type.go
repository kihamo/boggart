package owntracks

import (
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/atomic"
)

type Type struct{}

func (t Type) CreateBind(c interface{}) (boggart.Bind, error) {
	config := c.(*Config)
	region := make(map[string]*atomic.Bool, len(config.Regions))

	if len(config.Regions) > 0 {
		for n, r := range config.Regions {
			if r.GeoFence <= 0 {
				r.GeoFence = DefaultGeoFence
			}

			config.Regions[n] = r
			region[n] = atomic.NewBool()
		}
	}

	return &Bind{
		user:    config.User,
		device:  config.Device,
		regions: config.Regions,

		lat:     atomic.NewFloat64(),
		lon:     atomic.NewFloat64(),
		geoHash: atomic.NewString(),
		conn:    atomic.NewString(),
		acc:     atomic.NewInt64(),
		alt:     atomic.NewInt64(),
		batt:    atomic.NewFloat64(),
		vel:     atomic.NewInt64(),
		region:  region,
	}, nil
}
