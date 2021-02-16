package ping

import (
	"context"

	"github.com/kihamo/boggart/components/boggart/di"
	"github.com/sparrc/go-ping"
	"go.uber.org/multierr"
)

type Bind struct {
	di.ConfigBind
	di.LoggerBind
	di.MetricsBind
	di.MQTTBind
	di.ProbesBind

	config *Config
}

func (b *Bind) Check(ctx context.Context) error {
	pinger, err := ping.NewPinger(b.config.Hostname)
	if err != nil {
		return err
	}

	pinger.SetPrivileged(b.config.Privileged)

	pinger.Count = b.config.Retry
	pinger.Timeout = b.config.ReadinessProbeTimeout()

	pinger.Run()
	stats := pinger.Statistics()

	online := stats.PacketsRecv != 0

	var mqttError error

	if e := b.MQTT().PublishAsync(ctx, b.config.TopicOnline, online); e != nil {
		mqttError = multierr.Append(mqttError, e)
	}

	if online {
		latency := uint32(stats.MaxRtt.Nanoseconds() / 1e+6)
		metricLatency.With("host", b.config.Hostname).Set(float64(latency))

		if e := b.MQTT().PublishAsync(ctx, b.config.TopicLatency, latency); e != nil {
			mqttError = multierr.Append(mqttError, e)
		}
	}

	if mqttError != nil {
		b.Logger().Error(mqttError.Error())
	}

	return err
}
