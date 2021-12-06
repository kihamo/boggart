package premiergc

import (
	"context"
	"net/url"

	"github.com/hashicorp/go-multierror"
	"github.com/kihamo/boggart/components/boggart/di"
	"github.com/kihamo/boggart/providers/premiergc"
)

var (
	link, _ = url.Parse("https://my.premier-gc.ru/")
)

type Bind struct {
	di.ConfigBind
	di.LoggerBind
	di.MetaBind
	di.MetricsBind
	di.MQTTBind
	di.ProbesBind

	provider *premiergc.Client
}

func (b *Bind) config() *Config {
	return b.Config().Bind().(*Config)
}

func (b *Bind) Run() error {
	cfg := b.config()

	b.Meta().SetLink(link)
	b.provider = premiergc.New(cfg.Login, cfg.Password, cfg.Debug)

	return nil
}

func (b *Bind) Balance(ctx context.Context) (contract string, balance float64, err error) {
	contract, balance, err = b.provider.Balance(ctx)
	if err == nil {
		b.Meta().SetSerialNumber(contract)

		if e := b.MQTT().PublishAsync(ctx, b.config().TopicBalance.Format(contract), balance); e != nil {
			err = multierror.Append(err, e)
		}

		metricBalance.With("contract", contract).Set(balance)
	}

	return contract, balance, err
}
