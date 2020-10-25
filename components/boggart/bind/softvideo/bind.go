package softvideo

import (
	"context"

	"github.com/kihamo/boggart/components/boggart/di"
	"github.com/kihamo/boggart/providers/softvideo"
)

type Bind struct {
	di.LoggerBind
	di.MetricsBind
	di.MQTTBind
	di.ProbesBind

	config   *Config
	provider *softvideo.Client
}

func (b *Bind) Balance(ctx context.Context) (balance float64, err error) {
	balance, err = b.provider.Balance(ctx)

	if err == nil {
		metricBalance.With("account", b.config.Login).Set(balance)

		if e := b.MQTT().PublishAsync(ctx, b.config.TopicBalance, balance); e != nil {
			b.Logger().Error(e.Error())
		}
	}

	return balance, err
}
