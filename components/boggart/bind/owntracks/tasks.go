package owntracks

import (
	"context"
	"time"

	"github.com/kihamo/go-workers"
)

func (b *Bind) Tasks() []workers.Task {
	if !b.config.RegionsSyncEnabled {
		return nil
	}

	taskWayPoints := b.Workers().WrapTaskOnceSuccess(b.taskWayPoints)
	taskWayPoints.SetRepeats(-1)
	taskWayPoints.SetRepeatInterval(time.Second * 10)
	taskWayPoints.SetName("waypoints")

	return []workers.Task{
		taskWayPoints,
	}
}

func (b *Bind) taskWayPoints(context.Context) error {
	return b.CommandWayPoints()
}
