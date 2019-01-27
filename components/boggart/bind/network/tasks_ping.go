package network

import (
	"context"

	"github.com/kihamo/boggart/components/mqtt"
	"github.com/kihamo/go-workers"
	"github.com/kihamo/go-workers/task"
	"github.com/sparrc/go-ping"
	"go.uber.org/multierr"
)

func (b *BindPing) Tasks() []workers.Task {
	taskUpdater := task.NewFunctionTask(b.taskUpdater)
	taskUpdater.SetRepeats(-1)
	taskUpdater.SetRepeatInterval(b.updaterInterval)
	taskUpdater.SetName("bind-network:ping-updater-" + b.SerialNumber())

	return []workers.Task{
		taskUpdater,
	}
}

func (b *BindPing) taskUpdater(ctx context.Context) (interface{}, error) {
	pinger, err := ping.NewPinger(b.hostname)
	if err != nil {
		return nil, err
	}

	pinger.Count = b.retry
	pinger.Timeout = b.timeout

	pinger.Run()
	stats := pinger.Statistics()

	h := mqtt.NameReplace(b.hostname)

	online := stats.PacketsRecv != 0
	if ok := b.online.Set(online); ok {
		if e := b.MQTTPublishAsync(ctx, PingMQTTPublishTopicOnline.Format(h), 0, true, online); e != nil {
			err = multierr.Append(err, e)
		}
	}

	latency := uint32(stats.MaxRtt.Nanoseconds() / 1e+6)
	if ok := b.latency.Set(latency); ok {
		if e := b.MQTTPublishAsync(ctx, PingMQTTPublishTopicLatency.Format(h), 0, true, latency); e != nil {
			err = multierr.Append(err, e)
		}
	}

	return nil, err
}
