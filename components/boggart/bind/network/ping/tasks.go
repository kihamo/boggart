package ping

import (
	"context"

	"github.com/kihamo/boggart/components/mqtt"
	"github.com/kihamo/go-workers"
	"github.com/kihamo/go-workers/task"
	"github.com/sparrc/go-ping"
	"go.uber.org/multierr"
)

func (b *Bind) Tasks() []workers.Task {
	taskUpdater := task.NewFunctionTask(b.taskUpdater)
	taskUpdater.SetRepeats(-1)
	taskUpdater.SetRepeatInterval(b.updaterInterval)
	taskUpdater.SetName("bind-network:ping-updater-" + b.hostname)

	return []workers.Task{
		taskUpdater,
	}
}

func (b *Bind) taskUpdater(ctx context.Context) (interface{}, error) {
	pinger, err := ping.NewPinger(b.hostname)
	if err != nil {
		return nil, err
	}

	pinger.SetPrivileged(true)

	pinger.Count = b.retry
	pinger.Timeout = b.timeout

	pinger.Run()
	stats := pinger.Statistics()

	h := mqtt.NameReplace(b.hostname)

	online := stats.PacketsRecv != 0
	if ok := b.online.Set(online); ok {
		if e := b.MQTTPublishAsync(ctx, MQTTPublishTopicOnline.Format(h), online); e != nil {
			err = multierr.Append(err, e)
		}
	}

	if online {
		latency := uint32(stats.MaxRtt.Nanoseconds() / 1e+6)
		if ok := b.latency.Set(latency); ok {
			metricLatency.With("host", b.hostname).Set(float64(latency))

			if e := b.MQTTPublishAsync(ctx, MQTTPublishTopicLatency.Format(h), latency); e != nil {
				err = multierr.Append(err, e)
			}
		}
	}

	return nil, err
}
