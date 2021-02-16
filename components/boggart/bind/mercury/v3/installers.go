package v3

import (
	"context"
	"errors"

	"github.com/kihamo/boggart/components/boggart/installer"
	"github.com/kihamo/boggart/components/boggart/installer/openhab"
)

func (b *Bind) InstallersSupport() []installer.System {
	return []installer.System{
		installer.SystemOpenHab,
	}
}

func (b *Bind) InstallerSteps(context.Context, installer.System) ([]installer.Step, error) {
	meta := b.Meta()
	sn := meta.SerialNumber()
	if sn == "" {
		return nil, errors.New("serial number is empty")
	}

	itemPrefix := openhab.ItemPrefixFromBindMeta(meta)
	cfg := b.config()
	channels := make([]*openhab.Channel, 0, 14)

	const (
		idTariff          = "Tariff"
		idVoltagePhase1   = "VoltagePhase1"
		idVoltagePhase2   = "VoltagePhase2"
		idVoltagePhase3   = "VoltagePhase3"
		idAmperagePhase1  = "AmperagePhase1"
		idAmperagePhase2  = "AmperagePhase2"
		idAmperagePhase3  = "AmperagePhase3"
		idAmperageTotal   = "AmperageTotal"
		idPowerPhase1     = "PowerPhase1"
		idPowerPhase2     = "PowerPhase2"
		idPowerPhase3     = "PowerPhase3"
		idPowerTotal      = "PowerTotal"
		idMakeDate        = "MakeDate"
		idFirmwareVersion = "FirmwareVersion"
	)

	channels = append(channels,
		openhab.NewChannel(idTariff, openhab.ChannelTypeNumber).
			WithStateTopic(cfg.TopicTariff.Format(sn)).
			AddItems(
				openhab.NewItem(itemPrefix+idTariff, openhab.ItemTypeNumber).
					WithLabel("Tariff [JS(human_watts.js):%s]").
					WithIcon("energy"),
			),
		openhab.NewChannel(idVoltagePhase1, openhab.ChannelTypeNumber).
			WithStateTopic(cfg.TopicVoltagePhase1.Format(sn)).
			AddItems(
				openhab.NewItem(itemPrefix+idVoltagePhase1, openhab.ItemTypeNumber).
					WithLabel("Voltage phase 1 [%d V]").
					WithIcon("energy"),
			),
		openhab.NewChannel(idVoltagePhase2, openhab.ChannelTypeNumber).
			WithStateTopic(cfg.TopicVoltagePhase2.Format(sn)).
			AddItems(
				openhab.NewItem(itemPrefix+idVoltagePhase2, openhab.ItemTypeNumber).
					WithLabel("Voltage phase 2 [%d V]").
					WithIcon("energy"),
			),
		openhab.NewChannel(idVoltagePhase3, openhab.ChannelTypeNumber).
			WithStateTopic(cfg.TopicVoltagePhase3.Format(sn)).
			AddItems(
				openhab.NewItem(itemPrefix+idVoltagePhase3, openhab.ItemTypeNumber).
					WithLabel("Voltage phase 3 [%d V]").
					WithIcon("energy"),
			),
		openhab.NewChannel(idAmperagePhase1, openhab.ChannelTypeNumber).
			WithStateTopic(cfg.TopicAmperagePhase1.Format(sn)).
			AddItems(
				openhab.NewItem(itemPrefix+idAmperagePhase1, openhab.ItemTypeNumber).
					WithLabel("Amperage phase 1 [%.2f A]").
					WithIcon("energy"),
			),
		openhab.NewChannel(idAmperagePhase2, openhab.ChannelTypeNumber).
			WithStateTopic(cfg.TopicAmperagePhase2.Format(sn)).
			AddItems(
				openhab.NewItem(itemPrefix+idAmperagePhase2, openhab.ItemTypeNumber).
					WithLabel("Amperage phase 2 [%.2f A]").
					WithIcon("energy"),
			),
		openhab.NewChannel(idAmperagePhase3, openhab.ChannelTypeNumber).
			WithStateTopic(cfg.TopicAmperagePhase3.Format(sn)).
			AddItems(
				openhab.NewItem(itemPrefix+idAmperagePhase3, openhab.ItemTypeNumber).
					WithLabel("Amperage phase 3 [%.2f A]").
					WithIcon("energy"),
			),
		openhab.NewChannel(idAmperageTotal, openhab.ChannelTypeNumber).
			WithStateTopic(cfg.TopicAmperageTotal.Format(sn)).
			AddItems(
				openhab.NewItem(itemPrefix+idAmperageTotal, openhab.ItemTypeNumber).
					WithLabel("Amperage total [%.2f A]").
					WithIcon("energy"),
			),
		openhab.NewChannel(idPowerPhase1, openhab.ChannelTypeNumber).
			WithStateTopic(cfg.TopicPowerPhase1.Format(sn)).
			AddItems(
				openhab.NewItem(itemPrefix+idPowerPhase1, openhab.ItemTypeNumber).
					WithLabel("Power phase 1 [JS(human_watts.js):%s]").
					WithIcon("energy"),
			),
		openhab.NewChannel(idPowerPhase2, openhab.ChannelTypeNumber).
			WithStateTopic(cfg.TopicPowerPhase2.Format(sn)).
			AddItems(
				openhab.NewItem(itemPrefix+idPowerPhase2, openhab.ItemTypeNumber).
					WithLabel("Power phase 2 [JS(human_watts.js):%s]").
					WithIcon("energy"),
			),
		openhab.NewChannel(idPowerPhase3, openhab.ChannelTypeNumber).
			WithStateTopic(cfg.TopicPowerPhase3.Format(sn)).
			AddItems(
				openhab.NewItem(itemPrefix+idPowerPhase3, openhab.ItemTypeNumber).
					WithLabel("Power phase 3 [JS(human_watts.js):%s]").
					WithIcon("energy"),
			),
		openhab.NewChannel(idPowerTotal, openhab.ChannelTypeNumber).
			WithStateTopic(cfg.TopicPowerTotal.Format(sn)).
			AddItems(
				openhab.NewItem(itemPrefix+idPowerTotal, openhab.ItemTypeNumber).
					WithLabel("Power total [JS(human_watts.js):%s]").
					WithIcon("energy"),
			),
		openhab.NewChannel(idMakeDate, openhab.ChannelTypeDateTime).
			WithStateTopic(cfg.TopicMakeDate).
			AddItems(
				openhab.NewItem(itemPrefix+idMakeDate, openhab.ItemTypeDateTime).
					WithLabel("Make date [%1$td.%1$tm.%1$tY]").
					WithIcon("time"),
			),
		openhab.NewChannel(idFirmwareVersion, openhab.ChannelTypeString).
			WithStateTopic(cfg.TopicFirmwareVersion).
			AddItems(
				openhab.NewItem(itemPrefix+idFirmwareVersion, openhab.ItemTypeString).
					WithLabel("Firmware version [%s]"),
			),
	)

	return openhab.StepsByBind(b, []installer.Step{
		openhab.StepDefault(openhab.StepDefaultTransformHumanWatts),
	}, channels...)
}
