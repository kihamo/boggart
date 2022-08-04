package modbus

import (
	"context"
	"fmt"

	"github.com/kihamo/boggart/components/boggart/tasks"
	"go.uber.org/multierr"
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
	provider := b.Provider()
	id := b.Meta().ID()
	cfg := b.config()

	deviceType, err := provider.DeviceType()
	if err != nil {
		return fmt.Errorf("get device type failed: %w", err)
	}

	if e := b.MQTT().PublishAsync(ctx, cfg.TopicDeviceType.Format(id), deviceType); e != nil {
		err = multierr.Append(err, e)
	}

	if val, e := provider.HeatingOutputStatus(); e == nil {
		if e = b.MQTT().PublishAsync(ctx, cfg.TopicHeatingOutputStatus.Format(id), val); e != nil {
			err = multierr.Append(err, e)
		}
	} else {
		err = multierr.Append(err, fmt.Errorf("get heating output status failed: %w", e))
	}

	if val, e := provider.RoomTemperature(); e == nil {
		metricRoomTemperature.With("id", id).Set(val)

		if e = b.MQTT().PublishAsync(ctx, cfg.TopicRoomTemperature.Format(id), val); e != nil {
			err = multierr.Append(err, fmt.Errorf("get room temperature failed: %w", e))
		}
	} else {
		err = multierr.Append(err, e)
	}

	if val, e := provider.FloorTemperature(); e == nil {
		metricFlourTemperature.With("id", id).Set(val)

		if e = b.MQTT().PublishAsync(ctx, cfg.TopicFloorTemperature.Format(id), val); e != nil {
			err = multierr.Append(err, e)
		}
	} else {
		err = multierr.Append(err, fmt.Errorf("get floor temperature failed: %w", e))
	}

	if val, e := provider.Humidity(); e == nil {
		metricHumidity.With("id", id).Set(float64(val))

		if e = b.MQTT().PublishAsync(ctx, cfg.TopicHumidity.Format(id), val); e != nil {
			err = multierr.Append(err, e)
		}
	} else {
		err = multierr.Append(err, fmt.Errorf("get humidity failed: %w", e))
	}

	return err
}
