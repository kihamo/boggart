package hikvision

import (
	"context"
	"errors"
	"log"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/mqtt"
	"github.com/kihamo/go-workers"
	"github.com/kihamo/go-workers/task"
)

func (b *Bind) Tasks() []workers.Task {
	taskLiveness := task.NewFunctionTask(b.taskLiveness)
	taskLiveness.SetTimeout(b.livenessTimeout)
	taskLiveness.SetRepeats(-1)
	taskLiveness.SetRepeatInterval(b.livenessInterval)
	taskLiveness.SetName("bind-hikvision-liveness-" + b.address.Host)

	taskState := task.NewFunctionTask(b.taskUpdater)
	taskState.SetTimeout(b.updaterTimeout)
	taskState.SetRepeats(-1)
	taskState.SetRepeatInterval(b.updaterInterval)
	taskState.SetName("bind-hikvision-updater-" + b.address.Host)

	tasks := []workers.Task{
		taskLiveness,
		taskState,
	}

	if b.ptzEnabled {
		taskPTZStatus := task.NewFunctionTillStopTask(b.taskPTZ)
		taskPTZStatus.SetTimeout(b.ptzTimeout)
		taskPTZStatus.SetRepeats(-1)
		taskPTZStatus.SetRepeatInterval(b.ptzInterval)
		taskPTZStatus.SetName("bind-hikvision-ptz-" + b.address.Host)

		tasks = append(tasks, taskPTZStatus)
	}

	return tasks
}

func (b *Bind) taskLiveness(ctx context.Context) (interface{}, error) {
	deviceInfo, err := b.isapi.SystemDeviceInfo(ctx)
	if err != nil {
		b.UpdateStatus(boggart.BindStatusOffline)
		return nil, nil
	}

	if deviceInfo.SerialNumber == "" {
		b.UpdateStatus(boggart.BindStatusOffline)
		return nil, errors.New("device returns empty serial number")
	}

	if b.SerialNumber() == "" {
		ptzChannels := make(map[uint64]PTZChannel, 0)
		if list, err := b.isapi.PTZChannels(ctx); err == nil {
			for _, channel := range list.Channels {
				ptzChannels[channel.ID] = PTZChannel{
					Channel: channel,
				}
			}
		}

		b.mutex.Lock()
		b.ptzChannels = ptzChannels
		b.mutex.Unlock()

		if b.eventsEnabled {
			if err := b.startAlertStreaming(); err != nil {
				return nil, err
			}
		}

		b.SetSerialNumber(deviceInfo.SerialNumber)

		// TODO:
		_ = b.MQTTPublishAsync(ctx, MQTTPublishTopicStateModel.Format(deviceInfo.SerialNumber), 0, true, deviceInfo.Model)
		_ = b.MQTTPublishAsync(ctx, MQTTPublishTopicStateFirmwareVersion.Format(deviceInfo.SerialNumber), 0, true, deviceInfo.FirmwareVersion)
		_ = b.MQTTPublishAsync(ctx, MQTTPublishTopicStateFirmwareReleasedDate.Format(deviceInfo.SerialNumber), 0, true, deviceInfo.FirmwareReleasedDate)
	}

	b.UpdateStatus(boggart.BindStatusOnline)

	return nil, nil
}

func (b *Bind) taskPTZ(ctx context.Context) (interface{}, error, bool) {
	if b.Status() != boggart.BindStatusOnline {
		return nil, nil, false
	}

	b.mutex.RLock()
	defer b.mutex.RUnlock()

	if b.ptzChannels == nil {
		return nil, nil, false
	}

	if len(b.ptzChannels) == 0 {
		return nil, nil, true
	}

	stop := true

	for id, channel := range b.ptzChannels {
		if !channel.Channel.Enabled {
			continue
		}

		if err := b.updateStatusByChannelId(ctx, id); err != nil {
			log.Printf("failed updated status for %s device in channel %d", b.SerialNumber(), id)
			continue
		}

		stop = false
	}

	return nil, nil, stop
}

func (b *Bind) taskUpdater(ctx context.Context) (interface{}, error) {
	if b.Status() != boggart.BindStatusOnline {
		return nil, nil
	}

	b.mutex.RLock()
	defer b.mutex.RUnlock()

	sn := b.SerialNumber()
	snMQTT := mqtt.NameReplace(sn)

	status, err := b.isapi.SystemStatus(ctx)
	if err == nil {
		// TODO:
		_ = b.MQTTPublishAsync(ctx, MQTTPublishTopicStateUpTime.Format(snMQTT), 1, false, status.DeviceUpTime)
		metricUpTime.With("serial_number", sn).Set(float64(status.DeviceUpTime))

		memoryUsage := uint64(status.Memory[0].MemoryUsage.Float64()) * MB
		// TODO:
		_ = b.MQTTPublishAsync(ctx, MQTTPublishTopicStateMemoryUsage.Format(snMQTT), 1, false, memoryUsage)
		metricMemoryUsage.With("serial_number", sn).Set(float64(memoryUsage))

		memoryAvailable := uint64(status.Memory[0].MemoryAvailable.Float64()) * MB
		// TODO:
		_ = b.MQTTPublishAsync(ctx, MQTTPublishTopicStateMemoryAvailable.Format(snMQTT), 1, false, memoryAvailable)
		metricMemoryAvailable.With("serial_number", sn).Set(float64(memoryAvailable))
	} else {
		// TODO: log
	}

	storage, err := b.isapi.ContentManagementStorage(ctx)
	if err == nil {
		for _, hdd := range storage.HDD {
			// TODO:
			_ = b.MQTTPublishAsync(ctx, MQTTPublishTopicStateHDDCapacity.Format(snMQTT, hdd.ID), 1, false, hdd.Capacity*MB)

			// TODO:
			_ = b.MQTTPublishAsync(ctx, MQTTPublishTopicStateHDDUsage.Format(snMQTT, hdd.ID), 1, false, (hdd.Capacity-hdd.FreeSpace)*MB)
			metricStorageUsage.With("serial_number", sn).With("name", hdd.Name).Set(float64((hdd.Capacity - hdd.FreeSpace) * MB))

			// TODO:
			_ = b.MQTTPublishAsync(ctx, MQTTPublishTopicStateHDDFree.Format(snMQTT, hdd.ID), 1, false, hdd.FreeSpace*MB)
			metricStorageAvailable.With("serial_number", sn).With("name", hdd.Name).Set(float64(hdd.FreeSpace * MB))
		}
	} else {
		// TODO: log
	}

	return nil, nil
}
