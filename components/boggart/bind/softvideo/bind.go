package softvideo

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/kihamo/boggart/components/boggart/di"
	"github.com/kihamo/boggart/providers/softvideo"
)

type Bind struct {
	di.LoggerBind
	di.MetaBind
	di.MetricsBind
	di.MQTTBind
	di.ProbesBind

	config   *Config
	provider *softvideo.Client
}

func (b *Bind) Run() error {
	b.Meta().SetSerialNumber(b.config.Login)

	return nil
}

func (b *Bind) Balance(ctx context.Context) (balance, promise float64, err error) {
	balance, promise, err = b.provider.Balance(ctx)

	if err == nil {
		metricBalance.With("account", b.config.Login).Set(balance)

		if e := b.MQTT().PublishAsync(ctx, b.config.TopicBalance, balance); e != nil {
			err = multierror.Append(err, e)
		}

		if e := b.MQTT().PublishAsync(ctx, b.config.TopicPromise, promise); e != nil {
			err = multierror.Append(err, e)
		}
	}

	return balance, promise, err
}
