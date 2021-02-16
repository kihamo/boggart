package service

import (
	"context"
	"net"
	"strconv"
	"time"

	"github.com/kihamo/boggart/components/boggart/di"
	"go.uber.org/multierr"
)

type Bind struct {
	di.ConfigBind
	di.LoggerBind
	di.MetricsBind
	di.MQTTBind
	di.ProbesBind

	address string
}

func (b *Bind) config() *Config {
	return b.Config().Bind().(*Config)
}

func (b *Bind) Run() error {
	cfg := b.config()
	b.address = net.JoinHostPort(cfg.Hostname, strconv.Itoa(cfg.Port))

	return nil
}

func (b *Bind) Check(ctx context.Context) error {
	var (
		err      error
		attempts int
		latency  uint32
		conn     net.Conn
	)

	cfg := b.config()

	for {
		attempts++

		startTime := time.Now()
		conn, err = net.DialTimeout("tcp", b.address, cfg.ReadinessProbeTimeout())
		latency = uint32(time.Since(startTime).Nanoseconds() / 1e+6)

		if err == nil {
			conn.Close()
			break
		}

		if cfg.Retry > 0 && attempts >= cfg.Retry {
			break
		}
	}

	online := err == nil
	if e := b.MQTT().PublishAsync(ctx, cfg.TopicOnline.Format(b.address), online); e != nil {
		err = multierr.Append(err, e)
	}

	if online {
		metricLatency.With("address", b.address).Set(float64(latency))

		if e := b.MQTT().PublishAsync(ctx, cfg.TopicLatency.Format(b.address), latency); e != nil {
			err = multierr.Append(err, e)
		}
	}

	return err
}
