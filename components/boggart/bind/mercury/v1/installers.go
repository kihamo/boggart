package v1

import (
	"context"
	"fmt"
	"strconv"

	"github.com/kihamo/boggart/components/boggart/installer"
	"github.com/kihamo/boggart/components/boggart/installer/openhab"
)

func (b *Bind) InstallersSupport() []installer.System {
	return []installer.System{
		installer.SystemOpenHab,
	}
}

func (b *Bind) InstallerSteps(context.Context, installer.System) ([]installer.Step, error) {
	tariffCount, err := b.TariffCount()

	if err != nil {
		return nil, fmt.Errorf("get tariff count failed: %w", err)
	}

	meta := b.Meta()
	itemPrefix := openhab.ItemPrefixFromBindMeta(meta)
	cfg := b.config()

	channels := make([]*openhab.Channel, 0, tariffCount+9)
	transformHumanWatts := openhab.StepDefaultTransformHumanWatts.Base()

	for i := uint64(1); i <= uint64(tariffCount); i++ {
		index := strconv.FormatUint(i, 10)
		id := openhab.IDNormalizeCamelCase("Tariff " + index)

		channels = append(channels,
			openhab.NewChannel(id, openhab.ChannelTypeNumber).
				WithStateTopic(cfg.TopicTariff.Format(cfg.Address, index)).
				AddItems(
					openhab.NewItem(itemPrefix+id, openhab.ItemTypeNumber).
						WithLabel("Tariff "+index+" [JS("+transformHumanWatts+"):%s]").
						WithIcon("energy"),
				),
		)
	}

	const (
		idVoltage         = "Voltage"
		idAmperage        = "Amperage"
		idPower           = "Power"
		idBatteryVoltage  = "BatteryVoltage"
		idLastPowerOff    = "LastPowerOff"
		idLastPowerOn     = "LastPowerOn"
		idMakeDate        = "MakeDate"
		idFirmwareDate    = "FirmwareDate"
		idFirmwareVersion = "FirmwareVersion"
	)

	channels = append(channels,
		openhab.NewChannel(idVoltage, openhab.ChannelTypeNumber).
			WithStateTopic(cfg.TopicVoltage.Format(cfg.Address)).
			AddItems(
				openhab.NewItem(itemPrefix+idVoltage, openhab.ItemTypeNumber).
					WithLabel("Voltage [%d V]").
					WithIcon("energy"),
			),
		openhab.NewChannel(idAmperage, openhab.ChannelTypeNumber).
			WithStateTopic(cfg.TopicAmperage.Format(cfg.Address)).
			AddItems(
				openhab.NewItem(itemPrefix+idAmperage, openhab.ItemTypeNumber).
					WithLabel("Amperage [%.2f A]").
					WithIcon("energy"),
			),
		openhab.NewChannel(idPower, openhab.ChannelTypeNumber).
			WithStateTopic(cfg.TopicPower.Format(cfg.Address)).
			AddItems(
				openhab.NewItem(itemPrefix+idPower, openhab.ItemTypeNumber).
					WithLabel("Power [JS("+transformHumanWatts+"):%s]").
					WithIcon("energy"),
			),
		openhab.NewChannel(idBatteryVoltage, openhab.ChannelTypeNumber).
			WithStateTopic(cfg.TopicBatteryVoltage.Format(cfg.Address)).
			AddItems(
				openhab.NewItem(itemPrefix+idBatteryVoltage, openhab.ItemTypeNumber).
					WithLabel("Battery Voltage [%.2f V]").
					WithIcon("battery"),
			),
		openhab.NewChannel(idLastPowerOff, openhab.ChannelTypeDateTime).
			WithStateTopic(cfg.TopicLastPowerOff.Format(cfg.Address)).
			AddItems(
				openhab.NewItem(itemPrefix+idLastPowerOff, openhab.ItemTypeDateTime).
					WithLabel("Last power off [%1$td.%1$tm.%1$tY at %1$tH:%1$tM:%1$tS]").
					WithIcon("time"),
			),
		openhab.NewChannel(idLastPowerOn, openhab.ChannelTypeDateTime).
			WithStateTopic(cfg.TopicLastPowerOn.Format(cfg.Address)).
			AddItems(
				openhab.NewItem(itemPrefix+idLastPowerOn, openhab.ItemTypeDateTime).
					WithLabel("Last power on [%1$td.%1$tm.%1$tY at %1$tH:%1$tM:%1$tS]").
					WithIcon("time"),
			),
		openhab.NewChannel(idMakeDate, openhab.ChannelTypeDateTime).
			WithStateTopic(cfg.TopicMakeDate.Format(cfg.Address)).
			AddItems(
				openhab.NewItem(itemPrefix+idMakeDate, openhab.ItemTypeDateTime).
					WithLabel("Make date [%1$td.%1$tm.%1$tY]").
					WithIcon("time"),
			),
		openhab.NewChannel(idFirmwareDate, openhab.ChannelTypeDateTime).
			WithStateTopic(cfg.TopicFirmwareDate.Format(cfg.Address)).
			AddItems(
				openhab.NewItem(itemPrefix+idFirmwareDate, openhab.ItemTypeDateTime).
					WithLabel("Firmware date [%1$td.%1$tm.%1$tY]").
					WithIcon("time"),
			),
		openhab.NewChannel(idFirmwareVersion, openhab.ChannelTypeString).
			WithStateTopic(cfg.TopicFirmwareVersion.Format(cfg.Address)).
			AddItems(
				openhab.NewItem(itemPrefix+idFirmwareVersion, openhab.ItemTypeString).
					WithLabel("Firmware version [%s]"),
			),
	)

	return openhab.StepsByBind(b, []installer.Step{
		openhab.StepDefault(openhab.StepDefaultTransformHumanWatts),
	}, channels...)
}
