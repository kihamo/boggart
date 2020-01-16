package ping

import (
	"context"

	"github.com/kihamo/go-workers"
	"github.com/kihamo/go-workers/task"
	"github.com/sparrc/go-ping"
	"go.uber.org/multierr"
)

func (b *Bind) Tasks() []workers.Task {
	taskUpdater := task.NewFunctionTask(b.taskUpdater)
	taskUpdater.SetRepeats(-1)
	taskUpdater.SetRepeatInterval(b.config.UpdaterInterval)
	taskUpdater.SetName("updater")

	return []workers.Task{
		taskUpdater,
	}
}

func (b *Bind) taskUpdater(ctx context.Context) (interface{}, error) {
	pinger, err := ping.NewPinger(b.config.Hostname)
	if err != nil {
		return nil, err
	}

	pinger.SetPrivileged(true)

	pinger.Count = b.config.Retry
	pinger.Timeout = b.config.Timeout

	pinger.Run()
	stats := pinger.Statistics()

	online := stats.PacketsRecv != 0
	if e := b.MQTTContainer().PublishAsync(ctx, b.config.TopicOnline, online); e != nil {
		err = multierr.Append(err, e)
	}

	if online {
		latency := uint32(stats.MaxRtt.Nanoseconds() / 1e+6)
		metricLatency.With("host", b.config.Hostname).Set(float64(latency))

		if e := b.MQTTContainer().PublishAsync(ctx, b.config.TopicLatency, latency); e != nil {
			err = multierr.Append(err, e)
		}
	}

	return nil, err
}
