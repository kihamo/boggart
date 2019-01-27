package root

import (
	"github.com/kihamo/boggart/components/boggart"
)

type Type struct{}

func (t Type) CreateBind(c interface{}) (boggart.Bind, error) {
	config := c.(*Config)

	bind := &Bind{
		cacheRuntimeConfig: make(map[string]string, 11),
		watchFiles:         make(map[string]func(string) error, 0),
	}

	if config.DeviceIDFile != "" {
		if err := bind.InitDeviceID(config.DeviceIDFile); err != nil {
			return nil, err
		}
	}

	if config.RuntimeConfigFile != "" {
		if err := bind.AddWatchRuntimeConfig(config.RuntimeConfigFile); err != nil {
			return nil, err
		}
	}

	if err := bind.StartWatch(); err != nil {
		return nil, err
	}

	return bind, nil
}
