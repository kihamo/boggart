package nativeapi

import (
	"context"

	"github.com/kihamo/boggart/components/boggart/tasks"
)

func (b *Bind) Tasks() []tasks.Task {
	return []tasks.Task{
		tasks.NewTask().
			WithName("sync-state").
			WithHandler(
				b.Workers().WrapTaskHandlerIsOnline(
					tasks.HandlerFuncFromShortToLong(b.taskSyncStateHandler),
				),
			).
			WithSchedule(tasks.ScheduleWithDuration(tasks.ScheduleNow(), b.config().SyncStateInterval)),
	}
}

func (b *Bind) taskSyncStateHandler(ctx context.Context) error {
	if b.Meta().MAC() == nil {
		info, err := b.provider.DeviceInfo(ctx)
		if err != nil {
			return err
		}

		if err := b.Meta().SetMACAsString(info.MacAddress); err != nil {
			return err
		}
	}

	entities, err := b.provider.ListEntities(ctx)
	if err != nil {
		return err
	}

	return b.syncState(ctx, entities...)
}
