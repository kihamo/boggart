package softvideo

import (
	"context"
	"net/url"

	"github.com/hashicorp/go-multierror"
	"github.com/kihamo/boggart/components/boggart/di"
	"github.com/kihamo/boggart/providers/softvideo"
)

var (
	link, _ = url.Parse(softvideo.BaseURL)
)

type Bind struct {
	di.ConfigBind
	di.LoggerBind
	di.MetaBind
	di.MetricsBind
	di.MQTTBind
	di.ProbesBind

	provider *softvideo.Client
}

func (b *Bind) config() *Config {
	return b.Config().Bind().(*Config)
}

func (b *Bind) Run() error {
	cfg := b.config()

	b.Meta().SetSerialNumber(cfg.Login)
	b.Meta().SetLink(link)
	b.provider = softvideo.New(cfg.Login, cfg.Password, cfg.Debug)

	return nil
}

func (b *Bind) Balance(ctx context.Context) (balance, promise float64, err error) {
	balance, promise, err = b.provider.Balance(ctx)

	if err == nil {
		cfg := b.config()

		metricBalance.With("account", cfg.Login).Set(balance)

		if e := b.MQTT().PublishAsync(ctx, cfg.TopicBalance.Format(cfg.Login), balance); e != nil {
			err = multierror.Append(err, e)
		}

		if e := b.MQTT().PublishAsync(ctx, cfg.TopicPromise.Format(cfg.Login), promise); e != nil {
			err = multierror.Append(err, e)
		}
	}

	return balance, promise, err
}
