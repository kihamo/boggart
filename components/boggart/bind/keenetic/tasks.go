package keenetic

import (
	"context"
	"fmt"
	"go.uber.org/multierr"
	"time"

	"github.com/kihamo/boggart/components/boggart/tasks"
	"github.com/kihamo/boggart/providers/keenetic/client/show"
)

const (
	TaskNameUpdater = "updater"
)

func (b *Bind) Tasks() []tasks.Task {
	return []tasks.Task{
		tasks.NewTask().
			WithName("serial-number").
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
	defaults, err := b.client.Show.ShowDefaults(show.NewShowDefaultsParamsWithContext(ctx))
	if err != nil {
		return fmt.Errorf("get defaults value failed: %w", err)
	}

	b.Meta().SetSerialNumber(defaults.Payload.Serial)

	cfg := b.config()

	_, err = b.Workers().RegisterTask(
		tasks.NewTask().
			WithName(TaskNameUpdater).
			WithHandler(
				b.Workers().WrapTaskHandlerIsOnline(
					tasks.HandlerFuncFromShortToLong(b.taskUpdaterHandler),
				),
			).
			WithSchedule(tasks.ScheduleWithDuration(tasks.ScheduleNow(), cfg.UpdaterInterval)),
	)
	if err != nil {
		return fmt.Errorf("register task "+TaskNameUpdater+" failed: %w", err)
	}

	return nil
}

func (b *Bind) taskUpdaterHandler(ctx context.Context) (err error) {
	sn := b.Meta().SerialNumber()

	system, err := b.client.Show.ShowSystem(show.NewShowSystemParamsWithContext(ctx))
	if err == nil {
		metricUpTime.With("serial_number", sn).Set(float64(system.Payload.Uptime))
		metricCPULoad.With("serial_number", sn).Set(float64(system.Payload.Cpuload))
		metricMemoryAvailable.With("serial_number", sn).Set(float64(system.Payload.Memfree * 1024))
		metricMemoryUsage.With("serial_number", sn).Set(float64(system.Payload.Memtotal-system.Payload.Memfree) * 1024)
	} else {
		err = multierr.Append(err, fmt.Errorf("get system info failed: %w", err))
	}

	return err
}
