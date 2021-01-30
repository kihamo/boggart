package mosenergosbyt

import (
	"context"
	"strconv"
	"strings"

	"github.com/kihamo/boggart/components/boggart/config_generators"
	"github.com/kihamo/boggart/components/boggart/config_generators/openhab"
)

func (b *Bind) GenerateConfigOpenHab() ([]generators.Step, error) {
	ctx := context.Background()

	account, err := b.Account(ctx)
	if err != nil {
		return nil, err
	}

	balance, err := b.client.CurrentBalance(ctx, account.Provider.IDAbonent)
	if err != nil {
		return nil, err
	}

	itemPrefix := openhab.ItemPrefixFromBindMeta(b.Meta())

	const (
		idBalance        = "Balance"
		idAccountBalance = "Account" + idBalance
		idBill           = "Bill"
	)

	channels := []*openhab.Channel{
		openhab.NewChannel(idBill, openhab.ChannelTypeString).
			WithStateTopic(b.config.TopicLastBill.Format(account.NNAccount)).
			AddItems(
				openhab.NewItem(itemPrefix+idBill, openhab.ItemTypeString).
					WithLabel("Bill [%s]").
					WithIcon("text"),
			),
		openhab.NewChannel(idAccountBalance, openhab.ChannelTypeNumber).
			WithStateTopic(b.config.TopicBalance.Format(account.NNAccount)).
			AddItems(
				openhab.NewItem(itemPrefix+idAccountBalance, openhab.ItemTypeNumber).
					WithLabel("Account balance [%.2f ₽]").
					WithIcon("price"),
			),
	}

	var serviceID string

	for i, service := range balance.Services {
		if id, ok := services[strings.ToLower(service.Service)]; ok {
			serviceID = id
		} else {
			serviceID = strconv.FormatInt(int64(i), 10)
		}

		id := openhab.IDNormalizeCamelCase("Service "+serviceID) + "_" + idBalance

		channels = append(channels,
			openhab.NewChannel(id, openhab.ChannelTypeNumber).
				WithStateTopic(b.config.TopicServiceBalance.Format(account.NNAccount, serviceID)).
				AddItems(
					openhab.NewItem(itemPrefix+id, openhab.ItemTypeNumber).
						WithLabel(service.Service+" [%.2f ₽]").
						WithIcon("price"),
				),
		)
	}

	return openhab.StepsByBind(b, nil, channels...)
}
