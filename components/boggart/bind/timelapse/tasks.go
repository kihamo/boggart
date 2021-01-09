package timelapse

import (
	"context"

	"github.com/kihamo/boggart/components/boggart/tasks"
)

func (b *Bind) Tasks() []tasks.Task {
	return []tasks.Task{
		tasks.NewTask().
			WithName("updater").
			WithHandler(
				b.Workers().WrapTaskIsOnline(
					tasks.HandlerFuncFromShortToLong(b.taskUpdaterHandler),
				),
			).
			WithSchedule(tasks.ScheduleWithDuration(tasks.ScheduleNow(), b.config.UpdaterInterval)),
	}
}

func (b *Bind) taskUpdaterHandler(ctx context.Context) error {
	id := b.Meta().ID()
	if id == "" {
		return nil
	}

	files, err := b.Files(nil, nil)
	if err != nil {
		return err
	}

	metricTotalFiles.With("id", id).Set(float64(len(files)))

	var sizeTotal int64
	for _, f := range files {
		sizeTotal += f.Size()
	}

	metricTotalSize.With("id", id).Set(float64(sizeTotal))

	return nil
}
