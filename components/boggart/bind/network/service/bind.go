package service

import (
	"context"
	"net"
	"time"

	"github.com/kihamo/boggart/components/boggart/di"
	"go.uber.org/multierr"
)

type Bind struct {
	di.LoggerBind
	di.MetricsBind
	di.MQTTBind
	di.ProbesBind

	config  *Config
	address string
}

func (b *Bind) Check(ctx context.Context) error {
	var (
		err      error
		attempts int
		latency  uint32
		conn     net.Conn
	)

	for {
		attempts++

		startTime := time.Now()
		conn, err = net.DialTimeout("tcp", b.address, b.config.ReadinessProbeTimeout())
		latency = uint32(time.Since(startTime).Nanoseconds() / 1e+6)

		if err == nil {
			conn.Close()
			break
		}

		if b.config.Retry > 0 && attempts >= b.config.Retry {
			break
		}
	}

	online := err == nil
	if e := b.MQTT().PublishAsync(ctx, b.config.TopicOnline, online); e != nil {
		err = multierr.Append(err, e)
	}

	if online {
		metricLatency.With("address", b.address).Set(float64(latency))

		if e := b.MQTT().PublishAsync(ctx, b.config.TopicLatency, latency); e != nil {
			err = multierr.Append(err, e)
		}
	}

	return err
}
