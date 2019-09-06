package service

import (
	"context"
	"net"
	"time"

	"github.com/kihamo/boggart/components/mqtt"
	"github.com/kihamo/go-workers"
	"github.com/kihamo/go-workers/task"
	"go.uber.org/multierr"
)

func (b *Bind) Tasks() []workers.Task {
	taskUpdater := task.NewFunctionTask(b.taskUpdater)
	taskUpdater.SetRepeats(-1)
	taskUpdater.SetRepeatInterval(b.updaterInterval)
	taskUpdater.SetName("updater-" + b.address)

	return []workers.Task{
		taskUpdater,
	}
}

func (b *Bind) taskUpdater(ctx context.Context) (interface{}, error) {
	var (
		err      error
		attempts int
		latency  uint32
		conn     net.Conn
	)

	for {
		attempts++

		startTime := time.Now()
		conn, err = net.DialTimeout("tcp", b.address, b.timeout)
		latency = uint32(time.Since(startTime).Nanoseconds() / 1e+6)

		if err == nil {
			conn.Close()
			break
		}

		if b.retry > 0 && attempts >= b.retry {
			break
		}
	}

	h := mqtt.NameReplace(b.address)

	online := err == nil
	if e := b.MQTTPublishAsync(ctx, MQTTPublishTopicOnline.Format(h), online); e != nil {
		err = multierr.Append(err, e)
	}

	if online {
		metricLatency.With("address", b.address).Set(float64(latency))

		if e := b.MQTTPublishAsync(ctx, MQTTPublishTopicLatency.Format(h), latency); e != nil {
			err = multierr.Append(err, e)
		}
	}

	return nil, err
}
