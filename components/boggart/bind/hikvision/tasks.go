package hikvision

import (
	"context"
	"errors"

	"github.com/kihamo/boggart/providers/hikvision/client/content_manager"
	"github.com/kihamo/boggart/providers/hikvision/client/ptz"
	"github.com/kihamo/boggart/providers/hikvision/client/system"
	"github.com/kihamo/go-workers"
	"github.com/kihamo/go-workers/task"
	"go.uber.org/multierr"
)

func (b *Bind) Tasks() []workers.Task {
	taskState := b.Workers().WrapTaskIsOnline(b.taskUpdater)
	taskState.SetTimeout(b.config.UpdaterTimeout)
	taskState.SetRepeats(-1)
	taskState.SetRepeatInterval(b.config.UpdaterInterval)
	taskState.SetName("updater")

	tasks := []workers.Task{
		taskState,
	}

	if b.config.PTZEnabled {
		taskPTZStatus := task.NewFunctionTillStopTask(b.taskPTZ)
		taskPTZStatus.SetTimeout(b.config.PTZTimeout)
		taskPTZStatus.SetRepeats(-1)
		taskPTZStatus.SetRepeatInterval(b.config.PTZInterval)
		taskPTZStatus.SetName("ptz")

		tasks = append(tasks, taskPTZStatus)
	}

	return tasks
}

func (b *Bind) taskPTZ(ctx context.Context) (interface{}, error, bool) {
	if !b.Meta().IsStatusOnline() {
		return nil, nil, false
	}

	b.mutex.RLock()
	channels := b.ptzChannels
	b.mutex.RUnlock()

	if channels == nil {
		return nil, nil, false
	}

	if len(channels) == 0 {
		return nil, nil, true
	}

	stop := true

	for id, channel := range channels {
		if !channel.Channel.Enabled {
			continue
		}

		if err := b.updateStatusByChannelId(ctx, id); err != nil {
			b.Logger().Errorf("failed updated status for %s device in channel %d", b.Meta().SerialNumber(), id)
			continue
		}

		stop = false
	}

	return nil, nil, stop
}

func (b *Bind) taskUpdater(ctx context.Context) (err error) {
	sn := b.Meta().SerialNumber()

	// first initialization
	if sn == "" {
		deviceInfo, e := b.client.System.GetSystemDeviceInfo(system.NewGetSystemDeviceInfoParamsWithContext(ctx), nil)
		if e != nil {
			return e
		}

		if deviceInfo.Payload.SerialNumber == "" {
			return errors.New("device returns empty serial number")
		}

		if deviceInfo.Payload.MacAddress == "" {
			return errors.New("device returns empty MAC address")
		}

		b.Meta().SetSerialNumber(deviceInfo.Payload.SerialNumber)
		sn = deviceInfo.Payload.SerialNumber

		err = b.Meta().SetMACAsString(deviceInfo.Payload.MacAddress)
		if err != nil {
			return err
		}

		if b.config.EventsEnabled && b.config.EventsStreamingEnabled {
			b.startAlertStreaming()
		}

		if b.config.PTZEnabled {
			ptzChannels := make(map[uint64]PTZChannel)
			if list, e := b.client.Ptz.GetPtzChannels(ptz.NewGetPtzChannelsParamsWithContext(ctx), nil); e == nil {
				for _, channel := range list.Payload {
					ptzChannels[channel.ID] = PTZChannel{
						Channel: channel,
					}
				}
			} else {
				err = multierr.Append(err, e)
			}

			b.mutex.Lock()
			b.ptzChannels = ptzChannels
			b.mutex.Unlock()
		}

		if e := b.MQTT().PublishAsync(ctx, b.config.TopicStateModel.Format(deviceInfo.Payload.SerialNumber), deviceInfo.Payload.Model); e != nil {
			err = multierr.Append(err, e)
		}

		if e := b.MQTT().PublishAsync(ctx, b.config.TopicStateFirmwareVersion.Format(deviceInfo.Payload.SerialNumber), deviceInfo.Payload.FirmwareVersion); e != nil {
			err = multierr.Append(err, e)
		}

		if e := b.MQTT().PublishAsync(ctx, b.config.TopicStateFirmwareReleasedDate.Format(deviceInfo.Payload.SerialNumber), deviceInfo.Payload.FirmwareReleasedDate); e != nil {
			err = multierr.Append(err, e)
		}
	}

	b.mutex.RLock()
	defer b.mutex.RUnlock()

	if status, e := b.client.System.GetStatus(system.NewGetStatusParamsWithContext(ctx), nil); e == nil {
		if e := b.MQTT().PublishAsync(ctx, b.config.TopicStateUpTime.Format(sn), status.Payload.DeviceUpTime); e != nil {
			err = multierr.Append(err, e)
		}
		metricUpTime.With("serial_number", sn).Set(float64(status.Payload.DeviceUpTime))

		memoryUsage := int64(status.Payload.MemoryList[0].MemoryUsage) * MB
		if e := b.MQTT().PublishAsync(ctx, b.config.TopicStateMemoryUsage.Format(sn), memoryUsage); e != nil {
			err = multierr.Append(err, e)
		}
		metricMemoryUsage.With("serial_number", sn).Set(float64(memoryUsage))

		memoryAvailable := int64(status.Payload.MemoryList[0].MemoryAvailable) * MB
		if e := b.MQTT().PublishAsync(ctx, b.config.TopicStateMemoryAvailable.Format(sn), memoryAvailable); e != nil {
			err = multierr.Append(err, e)
		}
		metricMemoryAvailable.With("serial_number", sn).Set(float64(memoryAvailable))
	} else {
		err = multierr.Append(err, e)
	}

	if storage, e := b.client.ContentManager.GetStorage(content_manager.NewGetStorageParamsWithContext(ctx), nil); e == nil {
		for _, hdd := range storage.Payload.HddList {
			if hdd.Name == "" {
				continue
			}

			if e := b.MQTT().PublishAsync(ctx, b.config.TopicStateHDDCapacity.Format(sn, hdd.ID), hdd.Capacity*MB); e != nil {
				err = multierr.Append(err, e)
			}

			if e := b.MQTT().PublishAsync(ctx, b.config.TopicStateHDDUsage.Format(sn, hdd.ID), (hdd.Capacity-hdd.FreeSpace)*MB); e != nil {
				err = multierr.Append(err, e)
			}
			metricStorageUsage.With("serial_number", sn).With("name", hdd.Name).Set(float64((hdd.Capacity - hdd.FreeSpace) * MB))

			if e := b.MQTT().PublishAsync(ctx, b.config.TopicStateHDDFree.Format(sn, hdd.ID), hdd.FreeSpace*MB); e != nil {
				err = multierr.Append(err, e)
			}
			metricStorageAvailable.With("serial_number", sn).With("name", hdd.Name).Set(float64(hdd.FreeSpace * MB))
		}
	} else {
		err = multierr.Append(err, e)
	}

	return err
}
