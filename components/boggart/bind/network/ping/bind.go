package ping

import (
	"context"

	"github.com/kihamo/boggart/components/boggart/di"
	"github.com/sparrc/go-ping"
	"go.uber.org/multierr"
)

type Bind struct {
	di.MQTTBind
	di.LoggerBind
	di.ProbesBind

	config *Config
}

func (b *Bind) Check(ctx context.Context) error {
	pinger, err := ping.NewPinger(b.config.Hostname)
	if err != nil {
		return err
	}

	pinger.SetPrivileged(true)

	pinger.Count = b.config.Retry
	pinger.Timeout = b.config.ReadinessProbeTimeout()

	pinger.Run()
	stats := pinger.Statistics()

	online := stats.PacketsRecv != 0
	if e := b.MQTT().PublishAsync(ctx, b.config.TopicOnline, online); e != nil {
		err = multierr.Append(err, e)
	}

	if online {
		latency := uint32(stats.MaxRtt.Nanoseconds() / 1e+6)
		metricLatency.With("host", b.config.Hostname).Set(float64(latency))

		if e := b.MQTT().PublishAsync(ctx, b.config.TopicLatency, latency); e != nil {
			err = multierr.Append(err, e)
		}
	}

	return err
}
