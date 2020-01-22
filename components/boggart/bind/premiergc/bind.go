package premiergc

import (
	"context"

	"github.com/kihamo/boggart/components/boggart/di"
	"github.com/kihamo/boggart/providers/premiergc"
)

type Bind struct {
	di.MetaBind
	di.MQTTBind
	di.ProbesBind
	di.LoggerBind

	config   *Config
	provider *premiergc.Client
}

func (b *Bind) Balance(ctx context.Context) (contract string, balance float64, err error) {
	contract, balance, err = b.provider.Balance(ctx)
	if err == nil {
		b.Meta().SetSerialNumber(contract)

		metricBalance.With("contract", contract).Set(balance)
		err = b.MQTT().PublishAsync(ctx, b.config.TopicBalance.Format(contract), balance)
	}

	return contract, balance, err
}
