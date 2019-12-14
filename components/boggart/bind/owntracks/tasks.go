package owntracks

import (
	"context"
	"time"

	"github.com/kihamo/go-workers"
	"github.com/kihamo/go-workers/task"
)

func (b *Bind) Tasks() []workers.Task {
	if !b.config.RegionsSyncEnabled {
		return nil
	}

	taskWayPoints := task.NewFunctionTillSuccessTask(b.taskWayPoints)
	taskWayPoints.SetRepeats(-1)
	taskWayPoints.SetRepeatInterval(time.Second * 10)
	taskWayPoints.SetName("waypoints")

	return []workers.Task{
		taskWayPoints,
	}
}

func (b *Bind) taskWayPoints(context.Context) (interface{}, error) {
	return nil, b.CommandWayPoints()
}
