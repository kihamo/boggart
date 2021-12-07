package mosenergosbyt

import (
	"context"
	"errors"
	"net/url"

	"github.com/kihamo/boggart/components/boggart/di"
	"github.com/kihamo/boggart/providers/integratorit/mosenergosbyt"
)

var services = map[string]string{
	"взнос на капитальный ремонт": "10001",
	"обращение с тко":             "10002",
}

var link *url.URL

func init() {
	l, _ := url.Parse(mosenergosbyt.BaseURL)

	link = &url.URL{
		Scheme: l.Scheme,
		Host:   l.Host,
	}
}

const (
	layoutPeriod = "2006-01-02"
)

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
	cfg := b.config()

	b.client = mosenergosbyt.New(cfg.Login, cfg.Password)

	accounts, err := b.client.Accounts(context.Background())
	if err != nil {
		return err
	}

	for _, account := range accounts {
		if account.Provider.IDAbonent == 0 {
			continue
		}

		if account.Provider.IDAbonent > 0 && ((cfg.Account == "" && b.account == nil) || cfg.Account == account.NNAccount) {
			cfg.Account = account.NNAccount
			b.account = &account.NNAccount

			break
		}
	}

	b.Meta().SetLink(link)

	return nil
}

func (b *Bind) Account(ctx context.Context) (*mosenergosbyt.Account, error) {
	if b.account != nil {
		accounts, err := b.client.Accounts(ctx)
		if err != nil {
			return nil, err
		}

		for _, account := range accounts {
			if account.NNAccount == *b.account {
				return &account, nil
			}
		}
	}

	return nil, errors.New("account not found")
}
