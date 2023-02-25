package myheat

import (
	"context"
	"errors"
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

func (b *Bind) taskUpdaterHandler(ctx context.Context) error {
	response, err := b.client.Sensors.GetSensors(sensors.NewGetSensorsParamsWithContext(ctx), nil)
	if err != nil {
		return err
	}

	cfg := b.config()
	sn := b.Meta().SerialNumber()

	for _, sensor := range response.GetPayload() {
		if e := b.MQTT().PublishAsync(ctx, cfg.TopicSensorValue.Format(sn, sensor.ID), sensor.Value); e != nil {
			err = multierr.Append(err, e)
		}
	}

	return nil
}
