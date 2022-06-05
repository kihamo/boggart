package ping

import (
	"context"
	"errors"

	"github.com/go-ping/ping"
	"github.com/kihamo/boggart/components/boggart/di"
)

type Bind struct {
	di.ConfigBind
	di.LoggerBind
	di.MetaBind
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

	if stats.PacketsRecv == 0 {
		return errors.New("packets receive is 0")
	}

	latency := uint32(stats.MaxRtt.Nanoseconds() / 1e+6)
	metricLatency.With("host", cfg.Hostname).Set(float64(latency))

	return b.MQTT().PublishAsync(ctx, cfg.TopicLatency.Format(b.Meta().ID()), latency)
}
