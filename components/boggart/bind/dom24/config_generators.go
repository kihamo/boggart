package dom24

import (
	"fmt"

	"github.com/kihamo/boggart/components/boggart/config_generators"
	"github.com/kihamo/boggart/components/boggart/config_generators/openhab"
	"github.com/kihamo/boggart/providers/dom24/client/accounting"
)

func (b *Bind) GenerateConfigOpenHab() ([]generators.Step, error) {
	accountingResponse, err := b.provider.Accounting.AccountingInfo(accounting.NewAccountingInfoParams())
	if err != nil {
		return nil, fmt.Errorf("get accounting info failed: %w", err)
	}

	itemPrefix := openhab.ItemPrefixFromBindMeta(b.Meta())

	const (
		idBalance = "Balance"
		idBill    = "Bill"
	)

	channels := make([]*openhab.Channel, 0, len(accountingResponse.Payload.Data)*2)
	var id string

	for _, account := range accountingResponse.Payload.Data {
		id = "Account" + openhab.IDNormalizeCamelCase(account.Ident) + "_"

		channels = append(channels,
			openhab.NewChannel(id+idBalance, openhab.ChannelTypeNumber).
				WithStateTopic(b.config.TopicAccountBalance.Format(account.Ident)).
				AddItems(
					openhab.NewItem(itemPrefix+id+idBalance, openhab.ItemTypeNumber).
						WithLabel(account.AccountType+" #"+account.Ident+" [%.2f ₽]").
						WithIcon("price"),
				),
			openhab.NewChannel(id+idBill, openhab.ChannelTypeString).
				WithStateTopic(b.config.TopicAccountBill.Format(account.Ident)).
				AddItems(
					openhab.NewItem(itemPrefix+id+idBill, openhab.ItemTypeString).
						WithLabel("Bill #"+account.Ident+" [%s]").
						WithIcon("returnpipe"),
				),
		)
	}

	return openhab.StepsByBind(b, nil, channels...)
}
