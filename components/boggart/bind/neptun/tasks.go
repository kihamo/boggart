package neptun

import (
	"context"
	"errors"
	"strconv"

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
			WithSchedule(tasks.ScheduleWithDuration(tasks.ScheduleNow(), b.config().UpdaterInterval)),
	}
}

func (b *Bind) taskUpdaterHandler(ctx context.Context) error {
	provider := b.Provider()

	values, err := provider.CountersValues()
	if err != nil {
		return err
	}

	configs, err := provider.CountersConfigurations()
	if err != nil {
		return err
	}

	if len(values) != len(configs) {
		return errors.New("Counter values and configs slices sizes do not match")
	}

	for i, value := range values {
		cfg := configs[i]

		if cfg.Disabled() {
			continue
		}

		metricCounterValue.With("slot", strconv.Itoa(value.Slot()), "number", strconv.Itoa(value.Number())).Set(value.Value())
	}

	return err
}
