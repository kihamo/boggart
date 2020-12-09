package telegram

import (
	"github.com/kihamo/boggart/components/boggart"
	"gopkg.in/telegram-bot-api.v4"
)

type Type struct{}

func (t Type) CreateBind(c interface{}) (boggart.Bind, error) {
	bind := &Bind{
		config: c.(*Config),
	}

	l := NewLogger(
		func(message string) {
			bind.Logger().Debug(message)
		},
		func(message string) {
			bind.Logger().Warn(message)
		})

	tgbotapi.SetLogger(l)

	return bind, nil
}
