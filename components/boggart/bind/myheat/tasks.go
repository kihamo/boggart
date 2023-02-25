package myheat

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/kihamo/boggart/components/boggart/tasks"
	"github.com/kihamo/boggart/providers/myheat/device/client/sensors"
	"github.com/kihamo/boggart/providers/myheat/device/client/state"
	"go.uber.org/multierr"
)

const (
	TaskNameSerialNumber = "serial-number"
	TaskNameUpdater      = "updater"
)

func (b *Bind) Tasks() []tasks.Task {
	return []tasks.Task{
		tasks.NewTask().
			WithName(TaskNameSerialNumber).
			WithHandler(
				b.Workers().WrapTaskHandlerIsOnline(
					tasks.HandlerFuncFromShortToLong(b.taskSerialNumberHandler),
				),
			).
			WithSchedule(
				tasks.ScheduleWithSuccessLimit(
					tasks.ScheduleWithDuration(tasks.ScheduleNow(), time.Second*30),
					1,
				),
			),
	}
}

func (b *Bind) taskSerialNumberHandler(ctx context.Context) error {
	response, err := b.client.State.GetState(state.NewGetStateParamsWithContext(ctx), nil)
	if err != nil {
		return err
	}

	if sn := response.GetPayload().Serial; sn == "" {
		return errors.New("device returns empty serial number")
	} else {
		b.Meta().SetSerialNumber(sn)
	}

	cfg := b.config()

	_, err = b.Workers().RegisterTask(
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
	)

	return err
}

func (b *Bind) taskUpdaterHandler(ctx context.Context) (err error) {
	cfg := b.config()
	sn := b.Meta().SerialNumber()

	sensorsResponse, e := b.client.Sensors.GetSensors(sensors.NewGetSensorsParamsWithContext(ctx), nil)
	if e == nil {
		for _, sensor := range sensorsResponse.GetPayload() {
			if e := b.MQTT().PublishAsync(ctx, cfg.TopicSensorValue.Format(sn, sensor.ID), sensor.Value); e != nil {
				err = multierr.Append(err, fmt.Errorf("publish value for sensor %d return error: %w", sensor.ID, e))
			}
		}
	} else {
		err = multierr.Append(err, fmt.Errorf("get sensor state failed: %w", e))
	}

	stateObjResponse, e := b.client.State.GetObjState(state.NewGetObjStateParamsWithContext(ctx), nil)
	if e == nil {
		if v := stateObjResponse.Payload.SecurityArmed; v != nil {
			if e := b.MQTT().PublishAsync(ctx, cfg.TopicSecurityArmedState.Format(sn), v); e != nil {
				err = multierr.Append(err, fmt.Errorf("publish security armed state return error: %w", e))
			}
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
