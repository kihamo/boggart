package samsung_tizen

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/kihamo/go-workers"
	"github.com/kihamo/go-workers/task"
)

func (b *Bind) Tasks() []workers.Task {
	taskSerialNumber := task.NewFunctionTillSuccessTask(b.taskSerialNumber)
	taskSerialNumber.SetRepeats(-1)
	taskSerialNumber.SetRepeatInterval(time.Second * 30)
	taskSerialNumber.SetName("serial-number")

	return []workers.Task{
		taskSerialNumber,
	}
}

func (b *Bind) taskSerialNumber(ctx context.Context) (interface{}, error) {
	if !b.IsStatusOnline() {
		return nil, errors.New("bind isn't online")
	}

	info, err := b.client.Device(ctx)
	if err != nil {
		return nil, err
	}

	parts := strings.Split(info.ID, ":")
	if len(parts) > 1 {
		b.SetSerialNumber(parts[1])
	} else {
		b.SetSerialNumber(info.ID)
	}

	b.mutex.Lock()
	b.mac = info.Device.WifiMac
	b.mutex.Unlock()

	sn := b.SerialNumber()
	// TODO:
	_ = b.MQTTPublishAsync(ctx, b.config.TopicDeviceID.Format(sn), info.Device.ID)
	_ = b.MQTTPublishAsync(ctx, b.config.TopicDeviceModelName.Format(sn), info.Device.Name)

	return nil, nil
}
