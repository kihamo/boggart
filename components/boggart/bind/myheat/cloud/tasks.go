package cloud

import (
	"context"
	"fmt"
	"strconv"

	"github.com/kihamo/boggart/components/boggart/tasks"
	"github.com/kihamo/boggart/providers/myheat/cloud/client/operations"
	"github.com/kihamo/boggart/providers/myheat/cloud/models"
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
			WithSchedule(tasks.ScheduleWithDuration(tasks.ScheduleNow(), b.config().UpdaterInterval)),
	}
}

func (b *Bind) taskUpdaterHandler(ctx context.Context) (err error) {
	deviceIDs := make([]int64, 0)
	cfg := b.config()

	if cfg.DeviceID > 0 {
		deviceIDs = append(deviceIDs, cfg.DeviceID)
	} else {
		devicesResponse, e := b.client.Operations.GetDevices(operations.NewGetDevicesParamsWithContext(ctx))
		if e == nil {
			for _, device := range devicesResponse.GetPayload().Data.Devices {
				deviceIDs = append(deviceIDs, device.ID)
			}
		} else {
			err = fmt.Errorf("get devices list failed: %w", e)
		}
	}

	params := operations.NewGetDeviceInfoParamsWithContext(ctx)
	var (
		infoResponse                 *operations.GetDeviceInfoOK
		data                         *models.ResponseDeviceInfo
		deviceIdAsString, idAsString string
		e                            error
	)

	for _, deviceId := range deviceIDs {
		params.Request.DeviceID = deviceId
		deviceIdAsString = strconv.FormatInt(deviceId, 10)

		infoResponse, e = b.client.Operations.GetDeviceInfo(params)
		if e != nil {
			err = fmt.Errorf("get devices info %d failed: %w", deviceId, e)
			continue
		}

		data = infoResponse.GetPayload().Data

		if !data.DataActual {
			b.Logger().Warnf("Device info %d isn't actual", deviceId)
			continue
		}

		metricWeatherTemperature.With("device", deviceIdAsString).Set(data.WeatherTemp.Value())

		for _, heater := range data.Heaters {
			idAsString = strconv.FormatInt(heater.ID, 10)

			metricHeaterHeatingFeedTemperatureCelsius.With("device", deviceIdAsString, "id", idAsString).Set(heater.FlowTemp)
			metricHeaterHeatingReturnTemperatureCelsius.With("device", deviceIdAsString, "id", idAsString).Set(heater.ReturnTemp)
			metricHeaterHeatingTargetTemperatureCelsius.With("device", deviceIdAsString, "id", idAsString).Set(heater.TargetTemp)
			metricHeaterHeatingCircuitPressureBar.With("device", deviceIdAsString, "id", idAsString).Set(heater.Pressure)
			metricHeaterModulationPercent.With("device", deviceIdAsString, "id", idAsString).Set(float64(heater.Modulation))
		}

		for _, env := range data.Environments {
			idAsString = strconv.FormatInt(env.ID, 10)

			metricEnvironmentStateTemperatureCelsius.With("device", deviceIdAsString, "id", idAsString).Set(env.Value)
			// TODO: сделать корректную обработку null значения (когда контур отключен)
			metricEnvironmentTargetTemperatureCelsius.With("device", deviceIdAsString, "id", idAsString).Set(env.Target)
		}
	}

	return err
}
