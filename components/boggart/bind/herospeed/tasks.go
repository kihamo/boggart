package herospeed

import (
	"context"
	"errors"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/go-workers"
	"github.com/kihamo/go-workers/task"
	"go.uber.org/multierr"
)

func (b *Bind) Tasks() []workers.Task {
	taskLiveness := task.NewFunctionTask(b.taskLiveness)
	taskLiveness.SetTimeout(b.config.LivenessTimeout)
	taskLiveness.SetRepeats(-1)
	taskLiveness.SetRepeatInterval(b.config.LivenessInterval)
	taskLiveness.SetName("liveness-" + b.config.Address.Host)

	return []workers.Task{
		taskLiveness,
	}
}

func (b *Bind) taskLiveness(ctx context.Context) (interface{}, error) {
	configuration, err := b.client.Configuration(ctx)

	if err != nil {
		b.UpdateStatus(boggart.BindStatusOffline)
		return nil, nil
	}

	sn, ok := configuration["serialnumber"]
	if !ok || len(sn) == 0 {
		b.UpdateStatus(boggart.BindStatusOffline)
		return nil, errors.New("device returns empty serial number")
	}

	if b.SerialNumber() == "" {
		b.SetSerialNumber(sn)

		if model, ok := configuration["modelname"]; ok {
			if e := b.MQTTPublishAsync(ctx, b.config.TopicStateModel.Format(sn), model); e != nil {
				err = multierr.Append(err, e)
			}
		}

		if fw, ok := configuration["firmwareversion"]; ok {
			if e := b.MQTTPublishAsync(ctx, b.config.TopicStateFirmwareVersion.Format(sn), fw); e != nil {
				err = multierr.Append(err, e)
			}
		}
	}

	b.UpdateStatus(boggart.BindStatusOnline)

	return nil, err
}
