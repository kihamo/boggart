package esphome

import (
	"context"

	"github.com/kihamo/go-workers"
	"github.com/kihamo/go-workers/task"
)

func (b *Bind) Tasks() []workers.Task {
	taskSyncState := task.NewFunctionTask(b.taskSyncState)
	taskSyncState.SetRepeats(-1)
	taskSyncState.SetRepeatInterval(b.config.SyncStateInterval)
	taskSyncState.SetName("sync-state-" + b.config.Address)

	return []workers.Task{
		taskSyncState,
	}
}

func (b *Bind) taskSyncState(ctx context.Context) (interface{}, error) {
	if !b.IsStatusOnline() {
		return nil, nil
	}

	if b.SerialNumber() == "" {
		info, err := b.provider.DeviceInfo(ctx)
		if err != nil {
			return nil, err
		}

		b.SetSerialNumber(info.MacAddress)
	}

	entities, err := b.provider.ListEntities(ctx)
	if err != nil {
		return nil, err
	}

	return nil, b.syncState(ctx, entities...)
}
