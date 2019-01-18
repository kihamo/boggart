package alsa

import (
	"context"
	"sync/atomic"

	"github.com/kihamo/boggart/components/mqtt"
	"github.com/kihamo/go-workers"
	"github.com/kihamo/go-workers/task"
	"go.uber.org/multierr"
)

func (b *Bind) Tasks() []workers.Task {
	taskUpdater := task.NewFunctionTask(b.taskUpdater)
	taskUpdater.SetRepeats(-1)
	taskUpdater.SetRepeatInterval(b.updaterInterval)
	taskUpdater.SetName("bind-google-home:mini-updater-" + b.SerialNumber())

	return []workers.Task{
		taskUpdater,
	}
}

func (b *Bind) taskUpdater(ctx context.Context) (interface{}, error) {
	sn := mqtt.NameReplace(b.SerialNumber())
	var result error

	status := b.player.Status()
	prev := atomic.LoadInt64(&b.status)
	if status.Int64() != prev {
		atomic.StoreInt64(&b.status, status.Int64())

		if err := b.MQTTPublishAsync(ctx, MQTTPublishTopicStateStatus.Format(sn), 0, true, status.String()); err != nil {
			result = multierr.Append(result, err)
		}
	}

	if current, err := b.player.Volume(); err == nil {
		prev := atomic.LoadInt64(&b.volume)
		if current != prev {
			atomic.StoreInt64(&b.volume, current)

			if err := b.MQTTPublishAsync(ctx, MQTTPublishTopicStateVolume.Format(sn), 0, true, current); err != nil {
				result = multierr.Append(result, err)
			}
		}
	} else {
		result = multierr.Append(result, err)
	}

	if current, err := b.player.Mute(); err == nil {
		prev := atomic.LoadInt64(&b.mute)

		if prev == 0 || (prev == 1) != current {
			if current {
				atomic.StoreInt64(&b.mute, 1)
			} else {
				atomic.StoreInt64(&b.mute, -1)
			}

			if err := b.MQTTPublishAsync(ctx, MQTTPublishTopicStateMute.Format(sn), 0, true, current); err != nil {
				result = multierr.Append(result, err)
			}
		}
	} else {
		result = multierr.Append(result, err)
	}

	return nil, result
}
