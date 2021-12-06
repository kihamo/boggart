package dom24

import (
	"context"
	"fmt"
	"net/url"

	"github.com/kihamo/boggart/components/boggart/di"
	"github.com/kihamo/boggart/protocols/swagger"
	"github.com/kihamo/boggart/providers/dom24"
	"github.com/kihamo/boggart/providers/dom24/client/user"
)

var (
	link, _ = url.Parse("https://dom-24.net")
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

	provider *dom24.Client
}

func (b *Bind) config() *Config {
	return b.Config().Bind().(*Config)
}

func (b *Bind) Run() error {
	cfg := b.config()

	b.Meta().SetSerialNumber(cfg.Phone)
	b.Meta().SetLink(link)

	b.provider = dom24.New(cfg.Phone, cfg.Password, cfg.Debug, swagger.NewLogger(
		func(message string) {
			b.Logger().Info(message)
		},
		func(message string) {
			b.Logger().Debug(message)
		}))

	// пытаемся зарегистрировать счета, которые еще не зарегистрированы, но определены в конфиге
	if cfg.AutoRegisterIfNotExists && len(cfg.Accounts) > 0 {
		exists := make(map[string]bool, len(cfg.Accounts))
		for _, account := range cfg.Accounts {
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
