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
		b.UpdateStatus(boggart.DeviceStatusOffline)
		return nil, err
	}

	if deviceInfo.SerialNumber == "" {
		b.UpdateStatus(boggart.DeviceStatusOffline)
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

		b.MQTTPublishAsync(ctx, MQTTTopicStateModel.Format(deviceInfo.SerialNumber), 0, true, deviceInfo.Model)
		b.MQTTPublishAsync(ctx, MQTTTopicStateFirmwareVersion.Format(deviceInfo.SerialNumber), 0, true, deviceInfo.FirmwareVersion)
		b.MQTTPublishAsync(ctx, MQTTTopicStateFirmwareReleasedDate.Format(deviceInfo.SerialNumber), 0, true, deviceInfo.FirmwareReleasedDate)
	}

	b.UpdateStatus(boggart.DeviceStatusOnline)

	return nil, nil
}

func (b *Bind) taskPTZ(ctx context.Context) (interface{}, error, bool) {
	if b.Status() != boggart.DeviceStatusOnline {
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
	if b.Status() != boggart.DeviceStatusOnline {
		return nil, nil
	}

	b.mutex.RLock()
	defer b.mutex.RUnlock()

	status, err := b.isapi.SystemStatus(ctx)
	if err != nil {
		return nil, err
	}

	sn := mqtt.NameReplace(b.SerialNumber())

	b.MQTTPublishAsync(ctx, MQTTTopicStateUpTime.Format(sn), 1, false, status.DeviceUpTime)
	b.MQTTPublishAsync(ctx, MQTTTopicStateMemoryAvailable.Format(sn), 1, false, uint64(status.Memory[0].MemoryAvailable.Float64())*MB)
	b.MQTTPublishAsync(ctx, MQTTTopicStateMemoryUsage.Format(sn), 1, false, uint64(status.Memory[0].MemoryUsage.Float64())*MB)

	storage, err := b.isapi.ContentManagementStorage(ctx)
	if err != nil {
		return nil, err
	}

	for _, hdd := range storage.HDD {
		b.MQTTPublishAsync(ctx, MQTTTopicStateHDDCapacity.Format(sn, hdd.ID), 1, false, hdd.Capacity*MB)
		b.MQTTPublishAsync(ctx, MQTTTopicStateHDDFree.Format(sn, hdd.ID), 1, false, hdd.FreeSpace*MB)
		b.MQTTPublishAsync(ctx, MQTTTopicStateHDDUsage.Format(sn, hdd.ID), 1, false, (hdd.Capacity-hdd.FreeSpace)*MB)
	}

	return nil, nil
}
