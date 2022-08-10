package modbus

import (
	"github.com/kihamo/boggart/atomic"
	"github.com/kihamo/boggart/components/boggart"
)

type Type struct{}

func (t Type) CreateBind() boggart.Bind {
	return &Bind{
		providerOnce:    &atomic.Once{},
		statePower:      atomic.NewBool(),
		stateDeviceType: atomic.NewUint32(),
		//stateTemperatureFormat: atomic.NewUint32(),
		stateSetTemperature:     atomic.NewFloat64(),
		stateAway:               atomic.NewBool(),
		stateAwayTemperature:    atomic.NewUint32(),
		stateHoldingTemperature: atomic.NewUint32(),
	}
}
