package service

import (
	"context"
	"net"
	"net/url"
	"strconv"
	"time"

	"github.com/kihamo/boggart/components/boggart/di"
)

type Bind struct {
	di.ConfigBind
	di.LoggerBind
	di.MetaBind
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
	b.Meta().SetLink(&url.URL{
		Scheme: "http",
		Host:   b.address,
	})

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

	if err != nil {
		return err
	}

	metricLatency.With("address", b.address).Set(float64(latency))

	return b.MQTT().PublishAsync(ctx, cfg.TopicLatency.Format(b.Meta().ID()), latency)
}
