package dom24

import (
	"context"
	"fmt"

	"github.com/kihamo/boggart/components/boggart/di"
	"github.com/kihamo/boggart/providers/dom24"
	"github.com/kihamo/boggart/providers/dom24/client/user"
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

	config   *Config
	provider *dom24.Client
}

func (b *Bind) Run() error {
	b.Meta().SetSerialNumber(b.config.Phone)

	// пытаемся зарегистрировать счета, которые еще не зарегистрированы, но определены в конфиге
	if b.config.AutoRegisterIfNotExists && len(b.config.Accounts) > 0 {
		exists := make(map[string]bool, len(b.config.Accounts))
		for _, account := range b.config.Accounts {
			exists[account] = false
		}

		ctx := context.Background()

		userResponse, err := b.provider.User.UserInfo(user.NewUserInfoParamsWithContext(ctx))
		if err != nil {
			return err
		}

		for _, item := range userResponse.GetPayload().Accounts {
			if _, ok := exists[item.Ident]; ok {
				exists[item.Ident] = true
			}
		}

		for ident, exist := range exists {
			if exist {
				continue
			}

			params := user.NewAddByIdentParamsWithContext(ctx)
			params.Request.Ident = &[]string{ident}[0]

			if _, err = b.provider.User.AddByIdent(params); err != nil {
				return fmt.Errorf("add account #%s failed: %w", ident, err)
			}
		}
	}

	return nil
}
