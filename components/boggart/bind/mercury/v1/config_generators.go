package v1

import (
	"fmt"
	"strconv"

	"github.com/kihamo/boggart/components/boggart/config_generators"
	"github.com/kihamo/boggart/components/boggart/config_generators/openhab"
)

func (b *Bind) GenerateConfigOpenHab() ([]generators.Step, error) {
	tariffCount, err := b.TariffCount()

	if err != nil {
		return nil, fmt.Errorf("get tariff count failed: %w", err)
	}

	meta := b.Meta()
	itemPrefix := openhab.ItemPrefixFromBindMeta(meta)

	channels := make([]*openhab.Channel, 0, tariffCount+9)

	for i := uint64(1); i <= uint64(tariffCount); i++ {
		index := strconv.FormatUint(i, 10)
		id := openhab.IDNormalizeCamelCase("Tariff " + index)

		channels = append(channels,
			openhab.NewChannel(id, openhab.ChannelTypeNumber).
				WithStateTopic(b.config.TopicTariff.Format(index)).
				AddItems(
					openhab.NewItem(itemPrefix+id, openhab.ItemTypeNumber).
						WithLabel("Tariff "+index+" [JS(human_watts.js):%s]").
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
			WithStateTopic(b.config.TopicVoltage).
			AddItems(
				openhab.NewItem(itemPrefix+idVoltage, openhab.ItemTypeNumber).
					WithLabel("Voltage [%d V]").
					WithIcon("energy"),
			),
		openhab.NewChannel(idAmperage, openhab.ChannelTypeNumber).
			WithStateTopic(b.config.TopicAmperage).
			AddItems(
				openhab.NewItem(itemPrefix+idAmperage, openhab.ItemTypeNumber).
					WithLabel("Amperage [%.2f A]").
					WithIcon("energy"),
			),
		openhab.NewChannel(idPower, openhab.ChannelTypeNumber).
			WithStateTopic(b.config.TopicPower).
			AddItems(
				openhab.NewItem(itemPrefix+idPower, openhab.ItemTypeNumber).
					WithLabel("Power [JS(human_watts.js):%s]").
					WithIcon("energy"),
			),
		openhab.NewChannel(idBatteryVoltage, openhab.ChannelTypeNumber).
			WithStateTopic(b.config.TopicBatteryVoltage).
			AddItems(
				openhab.NewItem(itemPrefix+idBatteryVoltage, openhab.ItemTypeNumber).
					WithLabel("Battery Voltage [%.2f V]").
					WithIcon("battery"),
			),
		openhab.NewChannel(idLastPowerOff, openhab.ChannelTypeDateTime).
			WithStateTopic(b.config.TopicLastPowerOff).
			AddItems(
				openhab.NewItem(itemPrefix+idLastPowerOff, openhab.ItemTypeDateTime).
					WithLabel("Last power off [%1$td.%1$tm.%1$tY at %1$tH:%1$tM:%1$tS]").
					WithIcon("time"),
			),
		openhab.NewChannel(idLastPowerOn, openhab.ChannelTypeDateTime).
			WithStateTopic(b.config.TopicLastPowerOn).
			AddItems(
				openhab.NewItem(itemPrefix+idLastPowerOn, openhab.ItemTypeDateTime).
					WithLabel("Last power on [%1$td.%1$tm.%1$tY at %1$tH:%1$tM:%1$tS]").
					WithIcon("time"),
			),
		openhab.NewChannel(idMakeDate, openhab.ChannelTypeDateTime).
			WithStateTopic(b.config.TopicMakeDate).
			AddItems(
				openhab.NewItem(itemPrefix+idMakeDate, openhab.ItemTypeDateTime).
					WithLabel("Make date [%1$td.%1$tm.%1$tY]").
					WithIcon("time"),
			),
		openhab.NewChannel(idFirmwareDate, openhab.ChannelTypeDateTime).
			WithStateTopic(b.config.TopicFirmwareDate).
			AddItems(
				openhab.NewItem(itemPrefix+idFirmwareDate, openhab.ItemTypeDateTime).
					WithLabel("Firmware date [%1$td.%1$tm.%1$tY]").
					WithIcon("time"),
			),
		openhab.NewChannel(idFirmwareVersion, openhab.ChannelTypeString).
			WithStateTopic(b.config.TopicFirmwareVersion).
			AddItems(
				openhab.NewItem(itemPrefix+idFirmwareVersion, openhab.ItemTypeString).
					WithLabel("Firmware version [%s]"),
			),
	)

	return openhab.StepsByBind(b, []generators.Step{
		openhab.StepDefault(openhab.StepDefaultTransformHumanWatts),
	}, channels...)
}
