package telegram

import (
	"strconv"

	"github.com/kihamo/boggart/components/boggart"
	"gopkg.in/telegram-bot-api.v4"
)

type Type struct{}

func (t Type) CreateBind(c interface{}) (boggart.Bind, error) {
	config := c.(*Config)

	client, err := tgbotapi.NewBotAPI(config.Token)
	if err != nil {
		return nil, err
	}

	client.Debug = config.Debug

	bind := &Bind{
		client: client,
	}
	bind.SetSerialNumber(strconv.Itoa(client.Self.ID))

	return bind, nil
}
