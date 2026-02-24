package device

import (
	"context"
	"fmt"
	"strconv"

	"github.com/kihamo/boggart/providers/myheat"

	"github.com/kihamo/boggart/components/boggart/tasks"
	"github.com/kihamo/boggart/providers/myheat/device/client/sensors"
	"github.com/kihamo/boggart/providers/myheat/device/client/state"
	"go.uber.org/multierr"
)

const (
	TaskNameUpdater = "updater"
)

func (b *Bind) Tasks() []tasks.Task {
	cfg := b.config()

	return []tasks.Task{
		tasks.NewTask().
			WithName(TaskNameUpdater).
			WithHandler(
				b.Workers().WrapTaskHandlerIsOnline(
					tasks.HandlerWithTimeout(
						tasks.HandlerFuncFromShortToLong(b.taskUpdaterHandler),
						cfg.UpdaterTimeout,
					),
				),
			).
			WithSchedule(tasks.ScheduleWithDuration(tasks.ScheduleNow(), cfg.UpdaterInterval)),
	}
}

func (b *Bind) taskUpdaterHandler(ctx context.Context) (err error) {
	cfg := b.config()
	sn := b.Meta().SerialNumber()

	// устройство
	deviceResponse, e := b.client.State.GetState(state.NewGetStateParamsWithContext(ctx))
	if e == nil {
		if sn == "" && deviceResponse.Payload.Serial != "" {
			b.Meta().SetSerialNumber(deviceResponse.Payload.Serial)
			sn = deviceResponse.Payload.Serial
		}

		if e := b.MQTT().PublishAsync(ctx, cfg.TopicInternetConnected.Format(sn), deviceResponse.Payload.Inet); e != nil {
			err = multierr.Append(err, fmt.Errorf("publish internet connected return error: %w", e))
		}
	} else {
		err = multierr.Append(err, fmt.Errorf("get sensor state failed: %w", e))
	}

	// сенсоры
	sensorsResponse, e := b.client.Sensors.GetSensors(sensors.NewGetSensorsParamsWithContext(ctx))
	if e == nil {
		for _, sensor := range sensorsResponse.Payload {
			metricSensorValue.With("serial_number", sn).With("id", strconv.FormatInt(sensor.ID, 10)).Set(sensor.Value)

			if e := b.MQTT().PublishAsync(ctx, cfg.TopicSensorValue.Format(sn, sensor.ID), sensor.Value); e != nil {
				err = multierr.Append(err, fmt.Errorf("publish value for sensor %d return error: %w", sensor.ID, e))
			}
		}
	} else {
		err = multierr.Append(err, fmt.Errorf("get sensor state failed: %w", e))
	}

	// основные статусы
	stateObjResponse, e := b.client.State.GetObjState(state.NewGetObjStateParamsWithContext(ctx))
	if e == nil {
		// heaters
		for _, heater := range stateObjResponse.Payload.Heaters {
			heaterID := strconv.FormatInt(heater.ID, 10)

			if v, ok := heater.State[myheat.HeaterHeatingFeedTemperatureCelsius]; ok {
				metricHeaterHeatingFeedTemperatureCelsius.With("serial_number", sn).With("id", heaterID).Set(v)

				if e := b.MQTT().PublishAsync(ctx, cfg.TopicHeaterHeatingFeedTemperature.Format(sn, heaterID), v); e != nil {
					err = multierr.Append(err, fmt.Errorf("publish value for heater %d return error: %w", heaterID, e))
				}
			}

			if v, ok := heater.State[myheat.HeaterHeatingReturnTemperatureCelsius]; ok {
				metricHeaterHeatingReturnTemperatureCelsius.With("serial_number", sn).With("id", heaterID).Set(v)

				if e := b.MQTT().PublishAsync(ctx, cfg.TopicHeaterHeatingReturnTemperature.Format(sn, heaterID), v); e != nil {
					err = multierr.Append(err, fmt.Errorf("publish value for heater %d return error: %w", heaterID, e))
				}
			}

			if v, ok := heater.State[myheat.HeaterHeatingCircuitPressureBar]; ok {
				metricHeaterHeatingCircuitPressureBar.With("serial_number", sn).With("id", heaterID).Set(v)

				if e := b.MQTT().PublishAsync(ctx, cfg.TopicHeaterHeatingCircuitPressure.Format(sn, heaterID), v); e != nil {
					err = multierr.Append(err, fmt.Errorf("publish value for heater %d return error: %w", heaterID, e))
				}
			}
		}

		for _, env := range stateObjResponse.Payload.Environments {
			v, ok := env.State[myheat.EnvironmentStateTemperatureValue]
			if !ok {
				continue
			}

			envID := strconv.FormatInt(env.ID, 10)

			metricEnvironmentStateTemperatureCelsius.With("serial_number", sn).With("id", envID).Set(v)

			if e := b.MQTT().PublishAsync(ctx, cfg.TopicEnvironmentStateTemperature.Format(sn, envID), v); e != nil {
				err = multierr.Append(err, fmt.Errorf("publish value for environment %d return error: %w", envID, e))
			}

			// TODO: mqtt

			// зона должна контроллироваться, иначе не возможно установить таргетное значение
			//if v, ok := env.State[myheat.SettingsZoneControlled]; !ok || v != 1 {
			//	continue
			//}

			// таргетное значение установлено явно, а не выключено, хоть и контроллируется
			//if v, ok := env.State[myheat.EnvironmentSetTargetTemperature]; !ok || v != 1 {
			//	continue
			//}
		}

		if v := stateObjResponse.Payload.SecurityArmed; v != nil {
			if e := b.MQTT().PublishAsync(ctx, cfg.TopicSecurityArmedState.Format(sn), v); e != nil {
				err = multierr.Append(err, fmt.Errorf("publish security armed state return error: %w", e))
			}
		}

		if e := b.MQTT().PublishAsync(ctx, cfg.TopicAlarmPowerSupply.Format(sn), 2&stateObjResponse.Payload.DeviceFlags != 0); e != nil {
			err = multierr.Append(err, fmt.Errorf("publish alarm power supply return error: %w", e))
		}
		if e := b.MQTT().PublishAsync(ctx, cfg.TopicAlarmReplaceBattery.Format(sn), 4&stateObjResponse.Payload.DeviceFlags != 0); e != nil {
			err = multierr.Append(err, fmt.Errorf("publish alarm replace battery return error: %w", e))
		}
		if e := b.MQTT().PublishAsync(ctx, cfg.TopicAlarmGSMBalance.Format(sn), 16&stateObjResponse.Payload.DeviceFlags != 0); e != nil {
			err = multierr.Append(err, fmt.Errorf("publish alarm GSM balance return error: %w", e))
		}

		if e := b.MQTT().PublishAsync(ctx, cfg.TopicDeviceSeverity.Format(sn), stateObjResponse.Payload.DeviceSeverity); e != nil {
			err = multierr.Append(err, fmt.Errorf("publish device severity return error: %w", e))
		}

		if e := b.MQTT().PublishAsync(ctx, cfg.TopicGSMSignalLevel.Format(sn), stateObjResponse.Payload.SimSignal); e != nil {
			err = multierr.Append(err, fmt.Errorf("publish GSM signal return error: %w", e))
		}

		if e := b.MQTT().PublishAsync(ctx, cfg.TopicGSMBalance.Format(sn), stateObjResponse.Payload.SimBalance); e != nil {
			err = multierr.Append(err, fmt.Errorf("publish GSM balance return error: %w", e))
		}
	} else {
		err = multierr.Append(err, fmt.Errorf("get object state failed: %w", e))
	}

	return err
}
