package google_home

import (
	"context"
	"sync/atomic"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/providers/google/home"
	"github.com/kihamo/boggart/components/boggart/providers/google/home/client/info"
	"github.com/kihamo/boggart/components/mqtt"
	"github.com/kihamo/go-workers"
	"github.com/kihamo/go-workers/task"
)

func (b *Bind) Tasks() []workers.Task {
	taskLiveness := task.NewFunctionTask(b.taskLiveness)
	taskLiveness.SetTimeout(b.livenessTimeout)
	taskLiveness.SetRepeats(-1)
	taskLiveness.SetRepeatInterval(b.livenessInterval)
	taskLiveness.SetName("bind-google-home:mini-liveness-" + b.host)

	taskUpdater := task.NewFunctionTask(b.taskUpdater)
	taskUpdater.SetRepeats(-1)
	taskUpdater.SetRepeatInterval(b.updaterInterval)
	taskUpdater.SetName("bind-google-home:mini-updater-" + b.host)

	return []workers.Task{
		taskLiveness,
		taskUpdater,
	}
}

func (b *Bind) taskLiveness(ctx context.Context) (interface{}, error) {
	ctrl := b.ClientGoogleHome()

	response, err := ctrl.Info.GetEurekaInfo(info.NewGetEurekaInfoParams().
		WithOptions(home.EurekaInfoOptionDetail.Value()).
		WithParams(home.EurekaInfoParamDeviceInfo.Value()))
	if err != nil {
		b.UpdateStatus(boggart.BindStatusOffline)
		return nil, nil
	}

	if b.Status() == boggart.BindStatusOnline {
		return nil, nil
	}

	if b.SerialNumber() == "" && response.Payload != nil {
		b.SetSerialNumber(response.Payload.DeviceInfo.MacAddress)
	}

	b.UpdateStatus(boggart.BindStatusOnline)
	return nil, nil
}

func (b *Bind) taskUpdater(ctx context.Context) (interface{}, error) {
	ctrl := b.ClientChromeCast()
	sn := mqtt.NameReplace(b.SerialNumber())

	status := ctrl.Status()
	prev := atomic.LoadInt64(&b.status)
	if status.Int64() != prev {
		atomic.StoreInt64(&b.status, status.Int64())

		// TODO:
		_ = b.MQTTPublishAsync(ctx, MQTTPublishTopicStateStatus.Format(sn), 0, true, status.String())
	}

	if current, err := ctrl.Volume(); err == nil {
		prev := atomic.LoadInt64(&b.volume)
		if current != prev {
			atomic.StoreInt64(&b.volume, current)

			// TODO:
			_ = b.MQTTPublishAsync(ctx, MQTTPublishTopicStateVolume.Format(sn), 0, true, current)
		}
	} else {
		// TODO: log
	}

	if current, err := ctrl.Mute(); err == nil {
		prev := atomic.LoadInt64(&b.mute)

		if prev == 0 || (prev == 1) != current {
			if current {
				atomic.StoreInt64(&b.mute, 1)
			} else {
				atomic.StoreInt64(&b.mute, -1)
			}

			// TODO:
			_ = b.MQTTPublishAsync(ctx, MQTTPublishTopicStateMute.Format(sn), 0, true, current)
		}
	} else {
		// TODO:
	}

	return nil, nil
}
