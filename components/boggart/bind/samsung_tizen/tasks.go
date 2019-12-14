package samsung_tizen

import (
	"context"
	"strings"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/go-workers"
	"github.com/kihamo/go-workers/task"
)

func (b *Bind) Tasks() []workers.Task {
	taskLiveness := task.NewFunctionTask(b.taskLiveness)
	taskLiveness.SetTimeout(b.config.LivenessTimeout)
	taskLiveness.SetRepeats(-1)
	taskLiveness.SetRepeatInterval(b.config.LivenessInterval)
	taskLiveness.SetName("liveness")

	return []workers.Task{
		taskLiveness,
	}
}

func (b *Bind) taskLiveness(ctx context.Context) (interface{}, error) {
	info, err := b.client.Device(ctx)
	if err != nil {
		b.UpdateStatus(boggart.BindStatusOffline)
		return nil, nil
	}

	if b.IsStatusOnline() {
		return nil, nil
	}

	if b.SerialNumber() == "" {
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
	}
	b.UpdateStatus(boggart.BindStatusOnline)

	return nil, nil
}
