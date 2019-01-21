package owntracks

import (
	"github.com/kihamo/boggart/components/boggart"
)

type Type struct{}

func (t Type) CreateBind(c interface{}) (boggart.Bind, error) {
	config := c.(*Config)

	if len(config.Regions) > 0 {
		for name, region := range config.Regions {
			if region.GeoFence <= 0 {
				region.GeoFence = DefaultGeoFence
			}

			config.Regions[name] = region
		}
	}

	return &Bind{
		user:    config.User,
		device:  config.Device,
		regions: config.Regions,
	}, nil
}
