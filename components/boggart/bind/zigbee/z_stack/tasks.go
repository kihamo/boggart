package z_stack

import (
	"context"
	"encoding/hex"
	"time"

	"github.com/kihamo/go-workers"
)

func (b *Bind) Tasks() []workers.Task {
	taskSerialNumber := b.Workers().WrapTaskIsOnlineOnceSuccess(b.taskSerialNumber)
	taskSerialNumber.SetRepeats(-1)
	taskSerialNumber.SetRepeatInterval(time.Second * 30)
	taskSerialNumber.SetName("serial-number")

	return []workers.Task{
		taskSerialNumber,
	}
}

func (b *Bind) taskSerialNumber(ctx context.Context) error {
	client, err := b.getClient(ctx)
	if err != nil {
		return err
	}

	info, err := client.UtilGetDeviceInfo(ctx)
	if err != nil {
		return err
	}

	b.Meta().SetSerialNumber(hex.EncodeToString(info.IEEEAddr))
	b.syncPermitJoin()

	return nil
}
