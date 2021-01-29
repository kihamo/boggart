package owntracks

import (
	"context"
	"time"

	"github.com/kihamo/boggart/components/boggart/tasks"
)

func (b *Bind) Tasks() []tasks.Task {
	if !b.config.RegionsSyncEnabled {
		return nil
	}

	return []tasks.Task{
		tasks.NewTask().
			WithName("waypoints").
			WithHandler(
				b.Workers().WrapTaskHandlerIsOnline(
					tasks.HandlerFuncFromShortToLong(b.taskWayPointsHandler),
				),
			).
			WithSchedule(
				tasks.ScheduleWithSuccessLimit(
					tasks.ScheduleWithDuration(tasks.ScheduleNow(), time.Second*10),
					1,
				),
			),
	}
}

func (b *Bind) taskWayPointsHandler(context.Context) error {
	return b.CommandWayPoints()
}
