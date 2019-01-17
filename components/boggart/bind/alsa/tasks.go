package alsa

import (
	"context"
	"sync/atomic"

	"github.com/kihamo/boggart/components/mqtt"
	"github.com/kihamo/go-workers"
	"github.com/kihamo/go-workers/task"
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

	status := b.player.Status()
	prev := atomic.LoadInt64(&b.status)
	if status.Int64() != prev {
		atomic.StoreInt64(&b.status, status.Int64())

		// TODO:
		_ = b.MQTTPublishAsync(ctx, MQTTPublishTopicStateStatus.Format(sn), 0, true, status.String())
	}

	if current, err := b.player.Volume(); err == nil {
		prev := atomic.LoadInt64(&b.volume)
		if current != prev {
			atomic.StoreInt64(&b.volume, current)

			// TODO:
			_ = b.MQTTPublishAsync(ctx, MQTTPublishTopicStateVolume.Format(sn), 0, true, current)
		}
	} else {
		// TODO: log
	}

	if current, err := b.player.Mute(); err == nil {
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
