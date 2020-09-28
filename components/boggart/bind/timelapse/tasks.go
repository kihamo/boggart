package timelapse

import (
	"context"

	"github.com/kihamo/go-workers"
)

func (b *Bind) Tasks() []workers.Task {
	taskStateUpdater := b.Workers().WrapTaskIsOnline(b.taskUpdater)
	taskStateUpdater.SetRepeats(-1)
	taskStateUpdater.SetRepeatInterval(b.config.UpdaterInterval)
	taskStateUpdater.SetName("updater")

	return []workers.Task{
		taskStateUpdater,
	}
}

func (b *Bind) taskUpdater(ctx context.Context) error {
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
