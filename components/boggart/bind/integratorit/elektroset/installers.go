package elektroset

import (
	"context"
	"errors"
	"strconv"

	"github.com/kihamo/boggart/components/boggart/installer"
	"github.com/kihamo/boggart/components/boggart/installer/openhab"
)

func (b *Bind) InstallersSupport() []installer.System {
	return []installer.System{
		installer.SystemOpenHab,
	}
}

func (b *Bind) InstallerSteps(ctx context.Context, _ installer.System) ([]installer.Step, error) {
	if b.metersCount.IsNil() {
		return nil, errors.New("meters counter is nil")
	}

	accounts, err := b.client.Accounts(ctx)
	if err != nil {
		return nil, err
	}

	itemPrefix := openhab.ItemPrefixFromBindMeta(b.Meta())
	cfg := b.config()

	const (
		idBalance    = "Balance"
		idBill       = "Bill"
		idMeterValue = "Value"
		idMeterDate  = "Date"
	)

	channels := make([]*openhab.Channel, 0)
	var id string

	for _, account := range accounts {
		id = "Account" + openhab.IDNormalizeCamelCase(account.Number) + "_" + idBalance

		channels = append(channels,
			openhab.NewChannel(id, openhab.ChannelTypeNumber).
				WithStateTopic(cfg.TopicBalance.Format(account.Number)).
				AddItems(
					openhab.NewItem(itemPrefix+id, openhab.ItemTypeNumber).
						WithLabel("Account balance [%.2f ₽]").
						WithIcon("price"),
				),
		)

		for _, service := range account.Services {
			serviceID := strconv.FormatUint(service.ID, 10)
			id = "Service" + serviceID + "_"

			channels = append(channels,
				openhab.NewChannel(id+idBalance, openhab.ChannelTypeNumber).
					WithStateTopic(cfg.TopicServiceBalance.Format(account.Number, serviceID)).
					AddItems(
						openhab.NewItem(itemPrefix+id+idBalance, openhab.ItemTypeNumber).
							WithLabel("Balance "+service.NMServiceType+" [%.2f ₽]").
							WithIcon("price"),
					),
				openhab.NewChannel(id+idBill, openhab.ChannelTypeString).
					WithStateTopic(cfg.TopicLastBill.Format(account.Number, serviceID)).
					AddItems(
						openhab.NewItem(itemPrefix+id+idBill, openhab.ItemTypeString).
							WithLabel("Bill "+service.NMServiceType+" [%s]").
							WithIcon("returnpipe"),
					),
			)

			for i := uint32(0); i < b.metersCount.Load(); {
				i++

				meter := strconv.FormatUint(uint64(i), 10)
				id = "Service" + serviceID + "_Meter" + meter

				channels = append(channels,
					openhab.NewChannel(id+idMeterValue, openhab.ChannelTypeNumber).
						WithStateTopic(cfg.TopicMeterValue.Format(account.Number, serviceID, meter)).
						AddItems(
							openhab.NewItem(itemPrefix+id+idMeterValue, openhab.ItemTypeNumber).
								WithLabel("Tariff "+meter+" value [JS(human_watts.js):%s]").
								WithIcon("pressure"),
						),
					openhab.NewChannel(id+idMeterDate, openhab.ChannelTypeDateTime).
						WithStateTopic(cfg.TopicMeterDate.Format(account.Number, serviceID, meter)).
						AddItems(
							openhab.NewItem(itemPrefix+id+idMeterDate, openhab.ItemTypeDateTime).
								WithLabel("Tariff "+meter+" date [%1$td.%1$tm.%1$tY]").
								WithIcon("calendar"),
						),
				)
			}
		}
	}

	return openhab.StepsByBind(b, []installer.Step{
		openhab.StepDefault(openhab.StepDefaultTransformHumanWatts),
	}, channels...)
}
