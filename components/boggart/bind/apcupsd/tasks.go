package apcupsd

import (
	"context"
	"errors"
	"github.com/kihamo/boggart/components/boggart/tasks"
)

func (b *Bind) Tasks() []tasks.Task {
	return []tasks.Task{
		tasks.NewTask().
			WithName("updater").
			WithHandler(
				b.Workers().WrapTaskHandlerIsOnline(
					tasks.HandlerFuncFromShortToLong(b.taskUpdaterHandler),
				),
			).
			WithSchedule(
				tasks.ScheduleWithDuration(
					tasks.ScheduleNow(),
					b.config().UpdaterInterval,
				),
			),
	}
}

func (b *Bind) taskUpdaterHandler(ctx context.Context) error {
	status, err := b.client.Status(ctx)
	if err != nil {
		return err
	}

	if status.Status.IsCommunicationLost {
		return errors.New("communications with UPS lost")
	}

	sn := b.Meta().SerialNumber()

	if sn == "" {
		if status.SerialNumber == nil {
			sn = b.Meta().ID()
		} else {
			sn = *status.SerialNumber
		}

		b.Meta().SetSerialNumber(sn)
	}

	if sn == "" {
		return nil
	}

	_ = b.MQTT().PublishAsync(ctx, b.config().TopicVariable.Format(sn, VariableUPSName), status.UPSName)
	_ = b.MQTT().PublishAsync(ctx, b.config().TopicVariable.Format(sn, VariableVersion), status.Version)
	_ = b.MQTT().PublishAsync(ctx, b.config().TopicVariable.Format(sn, VariableModel), status.Model)
	_ = b.MQTT().PublishAsync(ctx, b.config().TopicVariable.Format(sn, VariableStatus), status.Status)

	if status.LineVoltage != nil {
		metricInputVoltage.With("serial_number", sn).Set(*status.LineVoltage)
		_ = b.MQTT().PublishAsync(ctx, b.config().TopicVariable.Format(sn, VariableLineVoltage), *status.LineVoltage)
	}

	if status.LoadPercent != nil {
		metricLoad.With("serial_number", sn).Set(*status.LoadPercent)
		_ = b.MQTT().PublishAsync(ctx, b.config().TopicVariable.Format(sn, VariableLoadPercent), *status.LoadPercent)
	}

	if status.BatteryChargePercent != nil {
		metricBatteryCharge.With("serial_number", sn).Set(*status.BatteryChargePercent)
		_ = b.MQTT().PublishAsync(ctx, b.config().TopicVariable.Format(sn, VariableBatteryChargePercent), *status.BatteryChargePercent)
	}

	if status.TimeLeft != nil {
		metricBatteryRuntime.With("serial_number", sn).Set((*status.TimeLeft).Seconds())
		_ = b.MQTT().PublishAsync(ctx, b.config().TopicVariable.Format(sn, VariableTimeLeft), (*status.TimeLeft).Seconds())
	}

	if status.MinimumBatteryChargePercent != nil {
		_ = b.MQTT().PublishAsync(ctx, b.config().TopicVariable.Format(sn, VariableMinimumBatteryChargePercent), *status.MinimumBatteryChargePercent)
	}

	if status.MaxLineVoltage != nil {
		_ = b.MQTT().PublishAsync(ctx, b.config().TopicVariable.Format(sn, VariableMaxLineVoltage), *status.MaxLineVoltage)
	}

	if status.MinLineVoltage != nil {
		_ = b.MQTT().PublishAsync(ctx, b.config().TopicVariable.Format(sn, VariableMinLineVoltage), *status.MinLineVoltage)
	}

	if status.OutputVoltage != nil {
		_ = b.MQTT().PublishAsync(ctx, b.config().TopicVariable.Format(sn, VariableOutputVoltage), *status.OutputVoltage)
	}

	if status.Sense != nil {
		_ = b.MQTT().PublishAsync(ctx, b.config().TopicVariable.Format(sn, VariableSense), *status.Sense)
	}

	if status.DelayShutdown != nil {
		_ = b.MQTT().PublishAsync(ctx, b.config().TopicVariable.Format(sn, VariableDelayShutdown), (*status.DelayShutdown).Seconds())
	}

	if status.DelayLowBattery != nil {
		_ = b.MQTT().PublishAsync(ctx, b.config().TopicVariable.Format(sn, VariableDelayLowBattery), (*status.DelayLowBattery).Seconds())
	}

	if status.LowTransferVoltage != nil {
		_ = b.MQTT().PublishAsync(ctx, b.config().TopicVariable.Format(sn, VariableLowTransferVoltage), *status.LowTransferVoltage)
	}

	if status.HighTransferVoltage != nil {
		_ = b.MQTT().PublishAsync(ctx, b.config().TopicVariable.Format(sn, VariableHighTransferVoltage), *status.HighTransferVoltage)
	}

	if status.InternalTemp != nil {
		_ = b.MQTT().PublishAsync(ctx, b.config().TopicVariable.Format(sn, VariableInternalTemp), *status.InternalTemp)
	}

	if status.BatteryVoltage != nil {
		metricBatteryVoltage.With("serial_number", sn).Set(*status.BatteryVoltage)
		_ = b.MQTT().PublishAsync(ctx, b.config().TopicVariable.Format(sn, VariableBatteryVoltage), *status.BatteryVoltage)
	}

	if status.LineFrequency != nil {
		_ = b.MQTT().PublishAsync(ctx, b.config().TopicVariable.Format(sn, VariableLineFrequency), *status.LineFrequency)
	}

	if status.LastTransfer != nil {
		_ = b.MQTT().PublishAsync(ctx, b.config().TopicVariable.Format(sn, VariableLastTransfer), *status.LastTransfer)
	}

	if status.SelfTest != nil {
		_ = b.MQTT().PublishAsync(ctx, b.config().TopicVariable.Format(sn, VariableSelfTest), *status.SelfTest)
	}

	if status.SelfTestInterval != nil {
		_ = b.MQTT().PublishAsync(ctx, b.config().TopicVariable.Format(sn, VariableSelfTestInterval), (*status.SelfTestInterval).Seconds())
	}

	if status.ManufacturedDate != nil {
		_ = b.MQTT().PublishAsync(ctx, b.config().TopicVariable.Format(sn, VariableManufacturedDate), *status.ManufacturedDate)
	}

	if status.SerialNumber != nil {
		_ = b.MQTT().PublishAsync(ctx, b.config().TopicVariable.Format(sn, VariableSerialNumber), *status.SerialNumber)
	}

	if status.BatteryDate != nil {
		_ = b.MQTT().PublishAsync(ctx, b.config().TopicVariable.Format(sn, VariableBatteryDate), *status.BatteryDate)
	}

	if status.NominalOutputVoltage != nil {
		_ = b.MQTT().PublishAsync(ctx, b.config().TopicVariable.Format(sn, VariableNominalOutputVoltage), *status.NominalOutputVoltage)
	}

	if status.NominalInputVoltage != nil {
		_ = b.MQTT().PublishAsync(ctx, b.config().TopicVariable.Format(sn, VariableNominalInputVoltage), *status.NominalInputVoltage)
	}

	if status.NominalBatteryVoltage != nil {
		_ = b.MQTT().PublishAsync(ctx, b.config().TopicVariable.Format(sn, VariableNominalBatteryVoltage), *status.NominalBatteryVoltage)
	}

	if status.NominalPower != nil {
		_ = b.MQTT().PublishAsync(ctx, b.config().TopicVariable.Format(sn, VariableNominalPower), *status.NominalPower)
	}

	if status.Humidity != nil {
		_ = b.MQTT().PublishAsync(ctx, b.config().TopicVariable.Format(sn, VariableHumidity), *status.Humidity)
	}

	if status.ExternalBatteries != nil {
		_ = b.MQTT().PublishAsync(ctx, b.config().TopicVariable.Format(sn, VariableExternalBatteries), *status.ExternalBatteries)
	}

	if status.BadBatteryPacks != nil {
		_ = b.MQTT().PublishAsync(ctx, b.config().TopicVariable.Format(sn, VariableBadBatteryPacks), *status.BadBatteryPacks)
	}

	if status.Firmware != nil {
		_ = b.MQTT().PublishAsync(ctx, b.config().TopicVariable.Format(sn, VariableFirmware), *status.Firmware)
	}

	if status.APCModel != nil {
		_ = b.MQTT().PublishAsync(ctx, b.config().TopicVariable.Format(sn, VariableAPCModel), *status.APCModel)
	}

	return nil
}
