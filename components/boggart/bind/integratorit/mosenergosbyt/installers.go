package mosenergosbyt

import (
	"context"
	"strings"
	"unicode"

	"github.com/kihamo/boggart/components/boggart/installer"
	"github.com/kihamo/boggart/components/boggart/installer/openhab"
)

func (b *Bind) InstallersSupport() []installer.System {
	return []installer.System{
		installer.SystemOpenHab,
	}
}

func (b *Bind) InstallerSteps(ctx context.Context, _ installer.System) ([]installer.Step, error) {
	account, err := b.Account(ctx)
	if err != nil {
		return nil, err
	}

	itemPrefix := openhab.ItemPrefixFromBindMeta(b.Meta())
	cfg := b.config()
	channels := make([]*openhab.Channel, 0)

	if b.client.IsSupportCurrentBalance(account) {
		const (
			idAccountBalance = "AccountBalance"
			idServiceBalance = "ServiceBalance"
		)

		balance, err := b.client.CurrentBalance(ctx, account)
		if err != nil {
			return nil, err
		}

		channels = append(channels, openhab.NewChannel(idAccountBalance, openhab.ChannelTypeNumber).
			WithStateTopic(cfg.TopicBalance.Format(account.AccountID)).
			AddItems(
				openhab.NewItem(itemPrefix+idAccountBalance, openhab.ItemTypeNumber).
					WithLabel("Account balance [%.2f ₽]").
					WithIcon("price"),
			))

		for _, service := range balance.Services {
			id := idServiceBalance + "_" + openhab.IDNormalizeCamelCase(service.ID)
			serviceName := strings.ToLower(service.Name)
			serviceName_ := []rune(serviceName)
			serviceName_[0] = unicode.ToUpper(serviceName_[0])
			serviceName = string(serviceName_)

			channels = append(channels,
				openhab.NewChannel(id, openhab.ChannelTypeNumber).
					WithStateTopic(cfg.TopicServiceBalance.Format(account.AccountID, service.ID)).
					AddItems(
						openhab.NewItem(itemPrefix+id, openhab.ItemTypeNumber).
							WithLabel(serviceName+" [%.2f ₽]").
							WithIcon("price"),
					),
			)
		}
	}

	if b.client.IsSupportBills(account) {
		const idBill = "Bill"

		channels = append(channels, openhab.NewChannel(idBill, openhab.ChannelTypeString).
			WithStateTopic(cfg.TopicLastBill.Format(account.AccountID)).
			AddItems(
				openhab.NewItem(itemPrefix+idBill, openhab.ItemTypeString).
					WithLabel("Bill [%s]").
					WithIcon("returnpipe"),
			))
	}

	//const idBalance        = "Balance"
	//var serviceID string
	//
	//for i, service := range balance.Services {
	//	if id, ok := services[strings.ToLower(service.Service)]; ok {
	//		serviceID = id
	//	} else {
	//		serviceID = strconv.FormatInt(int64(i), 10)
	//	}
	//
	//	id := openhab.IDNormalizeCamelCase("Service "+serviceID) + "_" + idBalance
	//

	//}

	return openhab.StepsByBind(b, nil, channels...)
}
