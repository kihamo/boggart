package apcupsd

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

func (b *Bind) InstallerSteps(ctx context.Context, system installer.System) ([]installer.Step, error) {
	status, err := b.client.Status(ctx)
	if err != nil {
		return nil, err
	}

	meta := b.Meta()
	sn := meta.SerialNumber()
	if sn == "" {
		return nil, errors.New("serial number is empty")
	}

	itemPrefix := openhab.ItemPrefixFromBindMeta(meta)
	cfg := b.config()
	channels := make([]*openhab.Channel, 0, 0)

	const (
		idVersion              = "Version"
		idModel                = "Model"
		idStatus               = "Status"
		idLineVoltage          = "LineVoltage"
		idLoadPercent          = "LoadPercent"
		idBatteryChargePercent = "BatteryChargePercent"
	)

	channels = append(channels,
		openhab.NewChannel(idVersion, openhab.ChannelTypeString).
			WithStateTopic(cfg.TopicVariable.Format(sn, VariableVersion)).
			AddItems(
				openhab.NewItem(itemPrefix+idVersion, openhab.ItemTypeString).
					WithLabel("Version"),
			),
		openhab.NewChannel(idModel, openhab.ChannelTypeString).
			WithStateTopic(cfg.TopicVariable.Format(sn, VariableModel)).
			AddItems(
				openhab.NewItem(itemPrefix+idModel, openhab.ItemTypeString).
					WithLabel("Model"),
			),
		openhab.NewChannel(idStatus, openhab.ChannelTypeString).
			WithStateTopic(cfg.TopicVariable.Format(sn, VariableStatus)).
			AddItems(
				openhab.NewItem(itemPrefix+idStatus, openhab.ItemTypeString).
					WithLabel("Status"),
			),
	)

	if status.LineVoltage != nil {
		channels = append(channels,
			openhab.NewChannel(idLineVoltage, openhab.ChannelTypeNumber).
				WithStateTopic(cfg.TopicVariable.Format(sn, VariableLineVoltage)).
				AddItems(
					openhab.NewItem(itemPrefix+idLineVoltage, openhab.ItemTypeNumber).
						WithLabel("Line voltage [%d V]").
						WithIcon("energy"),
				),
		)
	}

	if status.LoadPercent != nil {
		channels = append(channels,
			openhab.NewChannel(idLoadPercent, openhab.ChannelTypeNumber).
				WithStateTopic(cfg.TopicVariable.Format(sn, VariableLoadPercent)).
				AddItems(
					openhab.NewItem(itemPrefix+idLoadPercent, openhab.ItemTypeNumber).
						WithLabel("Load percent [%d %%]").
						WithIcon("batterylevel"),
				),
		)
	}

	if status.BatteryChargePercent != nil {
		channels = append(channels,
			openhab.NewChannel(idBatteryChargePercent, openhab.ChannelTypeNumber).
				WithStateTopic(cfg.TopicVariable.Format(sn, VariableBatteryChargePercent)).
				AddItems(
					openhab.NewItem(itemPrefix+idBatteryChargePercent, openhab.ItemTypeNumber).
						WithLabel("Battery charge percent [%d %%]").
						WithIcon("batterylevel"),
				),
		)
	}

	return openhab.StepsByBind(b, nil, channels...)
}
