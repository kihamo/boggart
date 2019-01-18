package xiaomi

import (
	"github.com/kihamo/boggart/components/boggart"
)

type RoborockRootType struct{}

func (t RoborockRootType) CreateBind(c interface{}) (boggart.Bind, error) {
	config := c.(*RoborockRootConfig)

	device := &RoborockRootBind{
		cacheRuntimeConfig: make(map[string]string, 11),
		watchFiles:         make(map[string]func(string) error, 0),
	}

	device.Init()

	if config.DeviceIDFile != "" {
		if err := device.InitDeviceID(config.DeviceIDFile); err != nil {
			return nil, err
		}
	}

	if config.RuntimeConfigFile != "" {
		if err := device.AddWatchRuntimeConfig(config.RuntimeConfigFile); err != nil {
			return nil, err
		}
	}

	if err := device.StartWatch(); err != nil {
		return nil, err
	}

	device.UpdateStatus(boggart.BindStatusOnline)

	return device, nil
}
