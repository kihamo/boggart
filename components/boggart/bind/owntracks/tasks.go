package owntracks

import (
	"context"
	"time"

	"github.com/kihamo/go-workers"
	"github.com/kihamo/go-workers/task"
)

func (b *Bind) Tasks() []workers.Task {
	if !b.config.WayPointsSyncEnabled {
		return nil
	}

	taskWayPoints := task.NewFunctionTillSuccessTask(b.taskWayPoints)
	taskWayPoints.SetRepeats(-1)
	taskWayPoints.SetRepeatInterval(time.Second * 10)
	taskWayPoints.SetName("bind-owntracks-waypoints-" + b.config.User + "-" + b.config.Device)

	return []workers.Task{
		taskWayPoints,
	}
}

func (b *Bind) taskWayPoints(context.Context) (interface{}, error) {
	return nil, b.CommandWayPoints()
}
