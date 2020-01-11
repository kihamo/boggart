package native_api

import (
	"context"

	"github.com/kihamo/go-workers"
)

func (b *Bind) Tasks() []workers.Task {
	taskSyncState := b.WrapTaskIsOnline(b.taskSyncState)
	taskSyncState.SetRepeats(-1)
	taskSyncState.SetRepeatInterval(b.config.SyncStateInterval)
	taskSyncState.SetName("sync-state")

	return []workers.Task{
		taskSyncState,
	}
}

func (b *Bind) taskSyncState(ctx context.Context) error {
	if b.SerialNumber() == "" {
		info, err := b.provider.DeviceInfo(ctx)
		if err != nil {
			return err
		}

		b.SetSerialNumber(info.MacAddress)
	}

	entities, err := b.provider.ListEntities(ctx)
	if err != nil {
		return err
	}

	return b.syncState(ctx, entities...)
}
