package devices

import (
	"context"

	"github.com/kihamo/boggart/components/boggart"
)

type DesktopPC struct {
	boggart.DeviceBase
	boggart.DeviceWOL
}

func NewDesktopPC() *DesktopPC {
	device := &DesktopPC{}
	device.Init()
	device.SetDescription("PC")

	return device
}

func (d *DesktopPC) Types() []boggart.DeviceType {
	return []boggart.DeviceType{
		boggart.DeviceTypePC,
	}
}

func (d *DesktopPC) Ping(ctx context.Context) bool {
	return false
}

func (d *DesktopPC) Enable() error {
	err := d.WakeUp()
	if err != nil {
		return err
	}

	return d.DeviceBase.Enable()
}
