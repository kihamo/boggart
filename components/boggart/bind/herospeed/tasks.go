package herospeed

import (
	"context"
	"errors"
	"time"

	"github.com/kihamo/go-workers"
	"github.com/kihamo/go-workers/task"
	"go.uber.org/multierr"
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
	if !b.Meta().IsStatusOnline() {
		return nil, errors.New("bind isn't online")
	}

	configuration, err := b.client.Configuration(ctx)
	if err != nil {
		return nil, err
	}

	sn, ok := configuration["serialnumber"]
	if !ok || len(sn) == 0 {
		return nil, errors.New("device returns empty serial number")
	}

	b.Meta().SetSerialNumber(sn)

	if model, ok := configuration["modelname"]; ok {
		if e := b.MQTT().PublishAsync(ctx, b.config.TopicStateModel.Format(sn), model); e != nil {
			err = multierr.Append(err, e)
		}
	}

	if fw, ok := configuration["firmwareversion"]; ok {
		if e := b.MQTT().PublishAsync(ctx, b.config.TopicStateFirmwareVersion.Format(sn), fw); e != nil {
			err = multierr.Append(err, e)
		}
	}

	return nil, err
}
