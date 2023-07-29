package timelapse

import (
	"context"
	"os"
	"path/filepath"
	"time"

	"github.com/kihamo/boggart/components/boggart/tasks"
)

func (b *Bind) Tasks() []tasks.Task {
	list := []tasks.Task{
		tasks.NewTask().
			WithName("updater").
			WithHandler(
				b.Workers().WrapTaskHandlerIsOnline(
					tasks.HandlerFuncFromShortToLong(b.taskUpdaterHandler),
				),
			).
			WithSchedule(tasks.ScheduleWithDuration(tasks.ScheduleNow(), b.config().UpdaterInterval)),
	}

	if b.config().EnableMigrationV1ToV2 {
		list = append(list,
			tasks.NewTask().
				WithName("migration-v1-to-v2").
				WithHandler(
					tasks.HandlerFuncFromShortToLong(b.taskMigrationV1ToV2),
				).
				WithSchedule(
					tasks.ScheduleWithSuccessLimit(
						tasks.ScheduleWithDuration(tasks.ScheduleNow(), time.Second*30),
						1,
					),
				))
	}

	return list
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

func (b *Bind) taskMigrationV1ToV2(_ context.Context) error {
	dir, err := os.Open(b.config().SaveDirectory)
	if err != nil {
		return err
	}
	defer dir.Close()

	files, err := dir.Readdir(-1)
	if err != nil {
		return err
	}

	var fileCount, dirCount int

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		subDirectory := filepath.Join(dir.Name(), file.ModTime().Format(SubDirectoryNameLayout))

		if _, err = os.Stat(subDirectory); os.IsNotExist(err) {
			if err = os.Mkdir(subDirectory, b.config().DirectoryPerm.FileMode); err != nil {
				return err
			}

			dirCount++
		}

		err = os.Rename(filepath.Join(dir.Name(), file.Name()), filepath.Join(subDirectory, file.Name()))
		if err != nil {
			return err
		}

		fileCount++
	}

	b.Logger().Infof("Migration %d files and created %d directories", fileCount, dirCount)

	return nil
}
