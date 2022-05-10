package apcupsd

import (
	"context"
	"errors"
	"time"

	"github.com/kihamo/boggart/components/boggart/installer"
	"github.com/kihamo/boggart/components/boggart/installer/openhab"
)

func (b *Bind) InstallersSupport() []installer.System {
	return []installer.System{
		installer.SystemOpenHab,
		installer.SystemDevice,
	}
}

func (b *Bind) InstallerSteps(ctx context.Context, system installer.System) ([]installer.Step, error) {
	if system == installer.SystemDevice {
		return []installer.Step{{
			Description: "Self testing",
			Content: `1. Stop service if running. service apcupsd stop
2. Run apctest
3. Select item "10) Perform battery calibration"
4. Waiting
5. Quit program Q
6. service apcupsd start`,
		}, {
			Description: "Change battery date",
			Content: `1. Stop service if running. service apcupsd stop
2. Run apctest
3. Select item "4)  View/Change battery date"
4. Enter new date in format MM/DD/YYYY, example ` + time.Now().Format("01/02/2006") + `
5. Quit program Q
6. service apcupsd start`,
		}}, nil
	}

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
		idVersion                     = "Version"
		idModel                       = "Model"
		idStatus                      = "Status"
		idLineVoltage                 = "LineVoltage"
		idLoadPercent                 = "LoadPercent"
		idBatteryChargePercent        = "BatteryChargePercent"
		idTimeLeft                    = "TimeLeft"
		idMinimumBatteryChargePercent = "MinimumBatteryChargePercent"
		idMaxLineVoltage              = "MaxLineVoltage"
		idMinLineVoltage              = "MinLineVoltage"
		idOutputVoltage               = "OutputVoltage"
		idBatteryVoltage              = "BatteryVoltage"
		idLastTransfer                = "LastTransfer"
		idNominalBatteryVoltage       = "NominalBatteryVoltage"
		idNominalPower                = "NominalPower"
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

	if status.TimeLeft != nil {
		channels = append(channels,
			openhab.NewChannel(idTimeLeft, openhab.ChannelTypeNumber).
				WithStateTopic(cfg.TopicVariable.Format(sn, VariableTimeLeft)).
				AddItems(
					openhab.NewItem(itemPrefix+idTimeLeft, openhab.ItemTypeNumber).
						WithLabel("Time left [JS("+openhab.StepDefaultTransformHumanSeconds.Base()+"):%s]").
						WithIcon("time"),
				),
		)
	}

	if status.MinimumBatteryChargePercent != nil {
		channels = append(channels,
			openhab.NewChannel(idMinimumBatteryChargePercent, openhab.ChannelTypeNumber).
				WithStateTopic(cfg.TopicVariable.Format(sn, VariableMinimumBatteryChargePercent)).
				AddItems(
					openhab.NewItem(itemPrefix+idMinimumBatteryChargePercent, openhab.ItemTypeNumber).
						WithLabel("Minimum battery charge percent [%d %%]").
						WithIcon("text"),
				),
		)
	}

	if status.MaxLineVoltage != nil {
		channels = append(channels,
			openhab.NewChannel(idMaxLineVoltage, openhab.ChannelTypeNumber).
				WithStateTopic(cfg.TopicVariable.Format(sn, VariableMaxLineVoltage)).
				AddItems(
					openhab.NewItem(itemPrefix+idMaxLineVoltage, openhab.ItemTypeNumber).
						WithLabel("Maximum line voltage [%d V]"),
				),
		)
	}

	if status.MinLineVoltage != nil {
		channels = append(channels,
			openhab.NewChannel(idMinLineVoltage, openhab.ChannelTypeNumber).
				WithStateTopic(cfg.TopicVariable.Format(sn, VariableMinLineVoltage)).
				AddItems(
					openhab.NewItem(itemPrefix+idMinLineVoltage, openhab.ItemTypeNumber).
						WithLabel("Minimum line voltage [%d V]"),
				),
		)
	}

	if status.OutputVoltage != nil {
		channels = append(channels,
			openhab.NewChannel(idOutputVoltage, openhab.ChannelTypeNumber).
				WithStateTopic(cfg.TopicVariable.Format(sn, VariableOutputVoltage)).
				AddItems(
					openhab.NewItem(itemPrefix+idOutputVoltage, openhab.ItemTypeNumber).
						WithLabel("Output voltage [%d V]"),
				),
		)
	}

	if status.Sense != nil {
		// TODO:
	}

	if status.DelayShutdown != nil {
		// TODO:
	}

	if status.DelayLowBattery != nil {
		// TODO:
	}

	if status.LowTransferVoltage != nil {
		// TODO:
	}

	if status.HighTransferVoltage != nil {
		// TODO:
	}

	if status.InternalTemp != nil {
		// TODO:
	}

	if status.BatteryVoltage != nil {
		channels = append(channels,
			openhab.NewChannel(idBatteryVoltage, openhab.ChannelTypeNumber).
				WithStateTopic(cfg.TopicVariable.Format(sn, VariableBatteryVoltage)).
				AddItems(
					openhab.NewItem(itemPrefix+idBatteryVoltage, openhab.ItemTypeNumber).
						WithLabel("Battery voltage [%.1f V]").
						WithIcon("energy"),
				),
		)
	}

	if status.LineFrequency != nil {
		// TODO:
	}

	if status.LastTransfer != nil {
		channels = append(channels,
			openhab.NewChannel(idLastTransfer, openhab.ChannelTypeString).
				WithStateTopic(cfg.TopicVariable.Format(sn, VariableLastTransfer)).
				AddItems(
					openhab.NewItem(itemPrefix+idLastTransfer, openhab.ItemTypeString).
						WithLabel("Reason transfer"),
				),
		)
	}

	if status.SelfTest != nil {
		// TODO:
	}

	if status.SelfTestInterval != nil {
		// TODO:
	}

	if status.ManufacturedDate != nil {
		// TODO:
	}

	if status.SerialNumber != nil {
		// TODO:
	}

	if status.BatteryDate != nil {
		// TODO:
	}

	if status.NominalOutputVoltage != nil {
		// TODO:
	}

	if status.NominalInputVoltage != nil {
		// TODO:
	}

	if status.NominalBatteryVoltage != nil {
		channels = append(channels,
			openhab.NewChannel(idNominalBatteryVoltage, openhab.ChannelTypeNumber).
				WithStateTopic(cfg.TopicVariable.Format(sn, VariableNominalBatteryVoltage)).
				AddItems(
					openhab.NewItem(itemPrefix+idNominalBatteryVoltage, openhab.ItemTypeNumber).
						WithLabel("Nominal battery voltage [%d V]").
						WithIcon("energy"),
				),
		)
	}

	if status.NominalPower != nil {
		channels = append(channels,
			openhab.NewChannel(idNominalPower, openhab.ChannelTypeNumber).
				WithStateTopic(cfg.TopicVariable.Format(sn, VariableNominalPower)).
				AddItems(
					openhab.NewItem(itemPrefix+idNominalPower, openhab.ItemTypeNumber).
						WithLabel("Nominal power [%d W]").
						WithIcon("energy"),
				),
		)
	}

	if status.Humidity != nil {
		// TODO:
	}

	if status.ExternalBatteries != nil {
		// TODO:
	}

	if status.BadBatteryPacks != nil {
		// TODO:
	}

	if status.Firmware != nil {
		// TODO:
	}

	return openhab.StepsByBind(b, []installer.Step{
		openhab.StepDefault(openhab.StepDefaultTransformHumanSeconds),
	}, channels...)
}
