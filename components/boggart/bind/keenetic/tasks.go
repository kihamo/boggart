package keenetic

import (
	"context"
	"fmt"
	"time"

	"github.com/kihamo/boggart/components/boggart/tasks"
	"github.com/kihamo/boggart/providers/keenetic/client/show"
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

	return nil
}
