package smcenter

import (
	"context"
	"fmt"
	"strconv"
	"unicode"

	"github.com/kihamo/boggart/components/boggart/installer"
	"github.com/kihamo/boggart/components/boggart/installer/openhab"
	"github.com/kihamo/boggart/providers/smcenter/client/accounting"
	"github.com/kihamo/boggart/providers/smcenter/client/meters"
)

func (b *Bind) InstallersSupport() []installer.System {
	return []installer.System{
		installer.SystemOpenHab,
	}
}

func (b *Bind) InstallerSteps(context.Context, installer.System) ([]installer.Step, error) {
	accountingResponse, err := b.provider.Accounting.AccountingInfo(accounting.NewAccountingInfoParams())
	if err != nil {
		return nil, fmt.Errorf("get accounting info failed: %w", err)
	}

	metersResponse, err := b.provider.Meters.List(meters.NewListParams())
	if err != nil {
		return nil, fmt.Errorf("get meters list failed: %w", err)
	}

	itemPrefix := openhab.ItemPrefixFromBindMeta(b.Meta())
	cfg := b.config()

	const (
		idBalance      = "Balance"
		idBill         = "Bill"
		idMeterCheckup = "Checkup"
		idMeterValue   = "Value"
	)

	channels := make([]*openhab.Channel, 0, len(accountingResponse.Payload.Data)*2)
	var id string

	for _, account := range accountingResponse.GetPayload().Data {
		id = "Account" + openhab.IDNormalizeCamelCase(account.Ident) + "_"

		channels = append(channels,
			openhab.NewChannel(id+idBalance, openhab.ChannelTypeNumber).
				WithStateTopic(cfg.TopicAccountBalance.Format(account.Ident)).
				AddItems(
					openhab.NewItem(itemPrefix+id+idBalance, openhab.ItemTypeNumber).
						WithLabel(account.AccountType+" #"+account.Ident+" [%.2f ₽]").
						WithIcon("price"),
				),
			openhab.NewChannel(id+idBill, openhab.ChannelTypeString).
				WithStateTopic(cfg.TopicAccountBill.Format(account.Ident)).
				AddItems(
					openhab.NewItem(itemPrefix+id+idBill, openhab.ItemTypeString).
						WithLabel("Bill #"+account.Ident+" [%s]").
						WithIcon("returnpipe"),
				),
		)
	}

	for _, meter := range metersResponse.GetPayload().Data {
		if meter.IsDisabled {
			continue
		}

		id = "Account" + openhab.IDNormalizeCamelCase(meter.Ident) + "_Meter" + openhab.IDNormalizeCamelCase(meter.FactoryNumber) + "_"

		name := meter.Name
		if name == "" && meter.CustomName != "" {
			name = meter.CustomName
		}
		if name == "" && meter.Resource != "" {
			name = meter.Resource
		}

		name_ := []rune(name)
		name_[0] = unicode.ToUpper(name_[0])
		name = string(name_)

		channels = append(channels,
			openhab.NewChannel(id+idMeterCheckup, openhab.ChannelTypeDateTime).
				WithStateTopic(cfg.TopicMeterCheckupDate.Format(meter.Ident, meter.FactoryNumber)).
				AddItems(
					openhab.NewItem(itemPrefix+id+idMeterCheckup, openhab.ItemTypeDateTime).
						WithLabel(name+" checkup date [%1$td.%1$tm.%1$tY]").
						WithIcon("time"),
				),
			openhab.NewChannel(id+idMeterValue, openhab.ChannelTypeNumber).
				WithStateTopic(cfg.TopicMeterValue.Format(meter.Ident, meter.FactoryNumber)).
				AddItems(
					openhab.NewItem(itemPrefix+id+idMeterValue, openhab.ItemTypeNumber).
						WithLabel(name+" [%."+strconv.FormatUint(meter.NumberOfDecimalPlaces, 10)+"f "+openhab.LabelEscape(meter.Units)+"]").
						WithIcon("heating-60"),
				),
		)
	}

	return openhab.StepsByBind(b, nil, channels...)
}
