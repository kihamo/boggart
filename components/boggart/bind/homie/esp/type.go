package esp

import (
	"sync"

	a "github.com/kihamo/boggart/atomic"
	"github.com/kihamo/boggart/components/boggart"
)

type Type struct {
	boggart.BindTypeWidget
}

func (t Type) CreateBind(c interface{}) (boggart.Bind, error) {
	config := c.(*Config)

	bind := &Bind{
		config:     config,
		lastUpdate: a.NewTimeNull(),

		deviceAttributes: &sync.Map{},

		nodes: &sync.Map{},

		otaEnabled:  a.NewBool(),
		otaRun:      a.NewBool(),
		otaWritten:  a.NewUint32(),
		otaTotal:    a.NewUint32(),
		otaChecksum: a.NewString(),
		otaFlash:    make(chan struct{}, 1),

		settings: &sync.Map{},
	}
	bind.SetSerialNumber(config.DeviceID)

	return bind, nil
}
