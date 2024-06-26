package mosenergosbyt

import (
	"context"
	"errors"
	"net/url"

	"github.com/kihamo/boggart/components/boggart/di"
	"github.com/kihamo/boggart/providers/integratorit/mosenergosbyt"
)

var link *url.URL

func init() {
	l, _ := url.Parse(mosenergosbyt.BaseURL)

	link = &url.URL{
		Scheme: l.Scheme,
		Host:   l.Host,
	}
}

type Bind struct {
	di.ConfigBind
	di.LoggerBind
	di.MetaBind
	di.MetricsBind
	di.MQTTBind
	di.ProbesBind
	di.WidgetBind
	di.WorkersBind

	client  *mosenergosbyt.Client
	account *string
}

func (b *Bind) config() *Config {
	return b.Config().Bind().(*Config)
}

func (b *Bind) Run() error {
	b.Meta().SetLink(link)

	cfg := b.config()

	b.client = mosenergosbyt.New(cfg.Login, cfg.Password)

	accounts, err := b.client.Accounts(context.Background())
	if err != nil {
		return err
	}

	for _, account := range accounts {
		if (cfg.Account == "" && b.account == nil) || cfg.Account == account.AccountID {
			cfg.Account = account.AccountID
			b.account = &account.AccountID

			break
		}
	}

	return nil
}

func (b *Bind) Account(ctx context.Context) (*mosenergosbyt.Account, error) {
	if b.account != nil {
		accounts, err := b.client.Accounts(ctx)
		if err != nil {
			return nil, err
		}

		for _, account := range accounts {
			if account.AccountID == *b.account {
				return &account, nil
			}
		}
	}

	return nil, errors.New("account not found")
}
