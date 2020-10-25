package premiergc

import (
	"context"

	"github.com/hashicorp/go-multierror"

	"github.com/kihamo/boggart/components/boggart/di"
	"github.com/kihamo/boggart/providers/premiergc"
)

type Bind struct {
	di.LoggerBind
	di.MetaBind
	di.MetricsBind
	di.MQTTBind
	di.ProbesBind

	config   *Config
	provider *premiergc.Client
}

func (b *Bind) Balance(ctx context.Context) (contract string, balance float64, err error) {
	contract, balance, err = b.provider.Balance(ctx)
	if err == nil {
		b.Meta().SetSerialNumber(contract)

		if e := b.MQTT().PublishAsync(ctx, b.config.TopicBalance.Format(contract), balance); e != nil {
			err = multierror.Append(err, e)
		}

		metricBalance.With("contract", contract).Set(balance)
	}

	return contract, balance, err
}
