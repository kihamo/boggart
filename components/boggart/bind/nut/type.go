package nut

import (
	"github.com/kihamo/boggart/components/boggart"
)

type Type struct{}

func (t Type) CreateBind(c interface{}) (boggart.DeviceBind, error) {
	device := &Bind{
		config:    c.(*Config),
		variables: make(map[string]interface{}, 0),
	}
	device.Init()

	return device, nil
}
