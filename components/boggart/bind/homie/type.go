package homie

import (
	"sync"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/shadow/components/dashboard"
)

type Type struct {
	dashboard.Handler
}

func (t Type) CreateBind(c interface{}) (boggart.Bind, error) {
	config := c.(*Config)

	bind := &Bind{
		config:               config,
		deviceAttributes:     &sync.Map{},
		implementationConfig: &sync.Map{},
	}
	bind.SetSerialNumber(config.DeviceID)

	return bind, nil
}
