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
}

func (b *Bind) config() *Config {
	return b.Config().Bind().(*Config)
}

func (b *Bind) Check(ctx context.Context) error {
	cfg := b.config()

	pinger, err := ping.NewPinger(cfg.Hostname)
	if err != nil {
		return err
	}

	pinger.SetPrivileged(cfg.Privileged)

	pinger.Count = cfg.Retry
	pinger.Timeout = cfg.ReadinessProbeTimeout()

	pinger.Run()
	stats := pinger.Statistics()

	online := stats.PacketsRecv != 0

	var mqttError error

	if e := b.MQTT().PublishAsync(ctx, cfg.TopicOnline.Format(cfg.Hostname), online); e != nil {
		mqttError = multierr.Append(mqttError, e)
	}

	if online {
		latency := uint32(stats.MaxRtt.Nanoseconds() / 1e+6)
		metricLatency.With("host", cfg.Hostname).Set(float64(latency))

		if e := b.MQTT().PublishAsync(ctx, cfg.TopicLatency.Format(cfg.Hostname), latency); e != nil {
			mqttError = multierr.Append(mqttError, e)
		}
	}

	if mqttError != nil {
		b.Logger().Error(mqttError.Error())
	}

	return err
}
