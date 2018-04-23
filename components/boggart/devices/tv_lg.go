package devices

import (
	"context"

	"github.com/dhickie/go-lgtv/control"
	"github.com/kihamo/boggart/components/boggart"
)

const (
	LGTVControlTimeout = 1000
)

type LGTV struct {
	boggart.DeviceBase
	boggart.DeviceWOL

	control *control.LgTv
}

func NewLGTV(control *control.LgTv) *LGTV {
	device := &LGTV{
		control: control,
	}
	device.Init()
	device.SetDescription("LG TV")

	return device
}

func (d *LGTV) connect() (err error) {
	if !d.control.IsConnected {
		_, err = d.control.Connect(d.control.ClientKey, LGTVControlTimeout)
	}

	return err
}

func (d *LGTV) Types() []boggart.DeviceType {
	return []boggart.DeviceType{
		boggart.DeviceTypeTV,
	}
}

func (d *LGTV) Ping(ctx context.Context) bool {
	if d.connect() != nil {
		return false
	}

	_, err := d.control.GetVolume()
	return err == nil
}

func (d *LGTV) Enable() error {
	err := d.WakeUp()
	if err != nil {
		return err
	}

	return d.DeviceBase.Enable()
}

func (d *LGTV) Disable() error {
	if d.control.IsConnected {
		d.control.TurnOff()
		d.control.Disconnect()
	}

	return d.DeviceBase.Disable()
}
