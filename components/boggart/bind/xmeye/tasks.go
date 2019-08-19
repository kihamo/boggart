package xmeye

import (
	"context"
	"errors"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/go-workers"
	"github.com/kihamo/go-workers/task"
)

func (b *Bind) Tasks() []workers.Task {
	taskLiveness := task.NewFunctionTask(b.taskLiveness)
	taskLiveness.SetTimeout(b.config.LivenessTimeout)
	taskLiveness.SetRepeats(-1)
	taskLiveness.SetRepeatInterval(b.config.LivenessInterval)
	taskLiveness.SetName("bind-xmeye-liveness-" + b.config.Address.Host)

	tasks := []workers.Task{
		taskLiveness,
	}

	return tasks
}

func (b *Bind) taskLiveness(ctx context.Context) (interface{}, error) {
	info, err := b.client.SystemInfo(ctx)
	if err != nil {
		b.UpdateStatus(boggart.BindStatusOffline)
		return nil, nil
	}

	if info.SerialNo == "" {
		b.UpdateStatus(boggart.BindStatusOffline)
		return nil, errors.New("device returns empty serial number")
	}

	if b.SerialNumber() == "" {
		b.SetSerialNumber(info.SerialNo)

		if b.config.AlarmStreamingEnabled {
			b.startAlarmStreaming()
		}
	}

	b.UpdateStatus(boggart.BindStatusOnline)

	return nil, err
}
