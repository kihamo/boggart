package owntracks

import (
	"context"
	"time"

	"github.com/kihamo/go-workers"
	"github.com/kihamo/go-workers/task"
)

func (b *Bind) Tasks() []workers.Task {
	tasks := make([]workers.Task, 0, len(b.devices))
	for user, device := range b.devices {
		task := task.NewFunctionTillSuccessTask(b.taskWayPoints(user, device))
		task.SetRepeats(-1)
		task.SetRepeatInterval(time.Second * 10)
		task.SetName("bind-owntracks-waypoints-" + user + "-" + device)

		tasks = append(tasks, task)
	}

	return tasks
}

func (b *Bind) taskWayPoints(user, device string) func(context.Context) (interface{}, error) {
	return func(context.Context) (interface{}, error) {
		return nil, b.CommandWayPoints(user, device)
	}
}
