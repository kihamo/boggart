package mosenergosbyt

import (
	"context"
	"errors"

	"github.com/kihamo/boggart/components/boggart/di"
	"github.com/kihamo/boggart/providers/integratorit/mosenergosbyt"
)

var services = map[string]string{
	"взнос на капитальный ремонт": "10001",
	"обращение с тко":             "10002",
}

const (
	layoutPeriod = "2006-01-02"
)

type Bind struct {
	di.MetaBind
	di.MQTTBind
	di.WorkersBind
	di.LoggerBind
	di.ProbesBind
	di.WidgetBind

	config  *Config
	client  *mosenergosbyt.Client
	account *string
}

func (b *Bind) Run() error {
	accounts, err := b.client.Accounts(context.Background())
	if err != nil {
		return err
	}

	for _, account := range accounts {
		if account.Provider.IDAbonent == 0 {
			continue
		}

		if account.Provider.IDAbonent > 0 && ((b.config.Account == "" && b.account == nil) || b.config.Account == account.NNAccount) {
			b.config.Account = account.NNAccount
			b.account = &account.NNAccount

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
			if account.NNAccount == *b.account {
				return &account, nil
			}
		}
	}

	return nil, errors.New("account not found")
}
