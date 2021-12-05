package elektroset

import (
	"context"
	"strconv"
	"time"

	"github.com/kihamo/boggart/components/boggart/installer"
	"github.com/kihamo/boggart/components/boggart/installer/openhab"
	"github.com/kihamo/boggart/providers/integratorit/elektroset"
)

func (b *Bind) InstallersSupport() []installer.System {
	return []installer.System{
		installer.SystemOpenHab,
	}
}

func (b *Bind) InstallerSteps(ctx context.Context, _ installer.System) ([]installer.Step, error) {
	houses, err := b.Houses(ctx)
	if err != nil {
		return nil, err
	}

	itemPrefix := openhab.ItemPrefixFromBindMeta(b.Meta())
	cfg := b.config()

	const (
		idBalance      = "Balance"
		idBill         = "Bill"
		idMeterValue   = "Value"
		idMeterDate    = "Date"
		idMeterCheckup = "Checkup"
	)

	channels := make([]*openhab.Channel, 0)
	var (
		idPrefix      string
		id            string
		meterTariff   string
		meterIDPrefix string
	)

	for _, house := range houses {
		idPrefix = "House" + strconv.FormatUint(house.ID, 10) + "_"
		id = idPrefix + idBalance

		channels = append(channels,
			openhab.NewChannel(id, openhab.ChannelTypeNumber).
				WithStateTopic(cfg.TopicBalance.Format(house.ID)).
				AddItems(
					openhab.NewItem(itemPrefix+id, openhab.ItemTypeNumber).
						WithLabel("House balance [%.2f ₽]").
						WithIcon("price"),
				),
		)

		for _, service := range house.Services {
			serviceID := strconv.FormatUint(service.ID, 10)
			id = idPrefix + "Service" + serviceID + "_"

			channels = append(channels,
				openhab.NewChannel(id+idBalance, openhab.ChannelTypeNumber).
					WithStateTopic(cfg.TopicServiceBalance.Format(house.ID, serviceID)).
					AddItems(
						openhab.NewItem(itemPrefix+id+idBalance, openhab.ItemTypeNumber).
							WithLabel("Balance "+service.ServiceTypeName+" [%.2f ₽]").
							WithIcon("price"),
					),
				openhab.NewChannel(id+idBill, openhab.ChannelTypeString).
					WithStateTopic(cfg.TopicLastBill.Format(house.ID, serviceID)).
					AddItems(
						openhab.NewItem(itemPrefix+id+idBill, openhab.ItemTypeString).
							WithLabel("Bill "+service.ServiceTypeName+" [%s]").
							WithIcon("returnpipe"),
					),
			)

			details, err := b.client.BalanceDetails(ctx, service.AccountID, service.Provider, time.Now().Add(-1*time.Hour*24*31*2), time.Now())
			if err != nil {
				return nil, err
			}

			meterIDPrefix = idPrefix + "Service" + serviceID + "_Meter"

			channels = append(channels,
				openhab.NewChannel(id+idMeterCheckup, openhab.ChannelTypeDateTime).
					WithStateTopic(cfg.TopicMeterCheckupDate.Format(house.ID, serviceID)).
					AddItems(
						openhab.NewItem(itemPrefix+id+idMeterCheckup, openhab.ItemTypeDateTime).
							WithLabel("Checkup date [%1$td.%1$tm.%1$tY]").
							WithIcon("time"),
					),
			)

			meterIDPrefix += "T"

			channelsExist := make(map[string]bool, 3)

			for _, balance := range details {
				if balance.TariffPlanEntityID != elektroset.TariffPlanEntityValue {
					continue
				}

				meterTariff = ""

				switch {
				case balance.T1Value != nil:
					meterTariff = "1"
				case balance.T2Value != nil:
					meterTariff = "2"
				case balance.T3Value != nil:
					meterTariff = "3"
				}

				if meterTariff == "" || channelsExist[meterTariff] {
					continue
				}

				channelsExist[meterTariff] = true
				id = meterIDPrefix + meterTariff

				channels = append(channels,
					openhab.NewChannel(id+idMeterValue, openhab.ChannelTypeNumber).
						WithStateTopic(cfg.TopicMeterValue.Format(house.ID, serviceID, meterTariff)).
						AddItems(
							openhab.NewItem(itemPrefix+id+idMeterValue, openhab.ItemTypeNumber).
								WithLabel("Tariff "+meterTariff+" value [JS("+openhab.StepDefaultTransformHumanWatts.Base()+"):%s]").
								WithIcon("pressure"),
						),
					openhab.NewChannel(id+idMeterDate, openhab.ChannelTypeDateTime).
						WithStateTopic(cfg.TopicMeterDate.Format(house.ID, serviceID, meterTariff)).
						AddItems(
							openhab.NewItem(itemPrefix+id+idMeterDate, openhab.ItemTypeDateTime).
								WithLabel("Tariff "+meterTariff+" date [%1$td.%1$tm.%1$tY]").
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
