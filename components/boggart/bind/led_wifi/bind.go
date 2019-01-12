package led_wifi

import (
	"context"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/providers/wifiled"
)

type Bind struct {
	statePower int64
	stateMode  uint64
	stateSpeed uint64
	stateColor uint64

	boggart.DeviceBindBase
	boggart.DeviceBindMQTT

	bulb *wifiled.Bulb
}

func (b *Bind) On(ctx context.Context) error {
	err := b.bulb.PowerOn(ctx)
	if err == nil {
		_, err = b.taskUpdater(ctx)
	}

	return err
}

func (b *Bind) Off(ctx context.Context) error {
	err := b.bulb.PowerOff(ctx)
	if err == nil {
		_, err = b.taskUpdater(ctx)
	}

	return err
}
