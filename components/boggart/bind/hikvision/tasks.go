package hikvision

import (
	"context"
	"errors"
	"log"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/providers/hikvision/client/content_manager"
	"github.com/kihamo/boggart/components/boggart/providers/hikvision/client/ptz"
	"github.com/kihamo/boggart/components/boggart/providers/hikvision/client/system"
	"github.com/kihamo/boggart/components/mqtt"
	"github.com/kihamo/go-workers"
	"github.com/kihamo/go-workers/task"
	"go.uber.org/multierr"
)

func (b *Bind) Tasks() []workers.Task {
	taskLiveness := task.NewFunctionTask(b.taskLiveness)
	taskLiveness.SetTimeout(b.config.LivenessTimeout)
	taskLiveness.SetRepeats(-1)
	taskLiveness.SetRepeatInterval(b.config.LivenessInterval)
	taskLiveness.SetName("bind-hikvision-liveness-" + b.address.Host)

	taskState := task.NewFunctionTask(b.taskUpdater)
	taskState.SetTimeout(b.config.UpdaterTimeout)
	taskState.SetRepeats(-1)
	taskState.SetRepeatInterval(b.config.UpdaterInterval)
	taskState.SetName("bind-hikvision-updater-" + b.address.Host)

	tasks := []workers.Task{
		taskLiveness,
		taskState,
	}

	if b.config.PTZEnabled {
		taskPTZStatus := task.NewFunctionTillStopTask(b.taskPTZ)
		taskPTZStatus.SetTimeout(b.config.PTZTimeout)
		taskPTZStatus.SetRepeats(-1)
		taskPTZStatus.SetRepeatInterval(b.config.PTZInterval)
		taskPTZStatus.SetName("bind-hikvision-ptz-" + b.address.Host)

		tasks = append(tasks, taskPTZStatus)
	}

	return tasks
}

func (b *Bind) taskLiveness(ctx context.Context) (interface{}, error) {
	deviceInfo, err := b.client.System.GetSystemDeviceInfo(system.NewGetSystemDeviceInfoParamsWithContext(ctx), nil)

	if err != nil {
		b.UpdateStatus(boggart.BindStatusOffline)
		return nil, nil
	}

	if deviceInfo.Payload.SerialNumber == "" {
		b.UpdateStatus(boggart.BindStatusOffline)
		return nil, errors.New("device returns empty serial number")
	}

	if b.SerialNumber() == "" {
		ptzChannels := make(map[uint64]PTZChannel)
		if list, err := b.client.Ptz.GetPtzChannels(ptz.NewGetPtzChannelsParamsWithContext(ctx), nil); err == nil {
			for _, channel := range list.Payload {
				ptzChannels[channel.ID] = PTZChannel{
					Channel: channel,
				}
			}
		}

		b.mutex.Lock()
		b.ptzChannels = ptzChannels
		b.mutex.Unlock()

		if b.config.EventsEnabled && b.config.EventsStreamingEnabled {
			b.startAlertStreaming()
		}

		b.SetSerialNumber(deviceInfo.Payload.SerialNumber)

		if e := b.MQTTPublishAsync(ctx, MQTTPublishTopicStateModel.Format(deviceInfo.Payload.SerialNumber), deviceInfo.Payload.Model); e != nil {
			err = multierr.Append(err, e)
		}

		if e := b.MQTTPublishAsync(ctx, MQTTPublishTopicStateFirmwareVersion.Format(deviceInfo.Payload.SerialNumber), deviceInfo.Payload.FirmwareVersion); e != nil {
			err = multierr.Append(err, e)
		}

		if e := b.MQTTPublishAsync(ctx, MQTTPublishTopicStateFirmwareReleasedDate.Format(deviceInfo.Payload.SerialNumber), deviceInfo.Payload.FirmwareReleasedDate); e != nil {
			err = multierr.Append(err, e)
		}
	}

	b.UpdateStatus(boggart.BindStatusOnline)

	return nil, err
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

	status, err := b.client.System.GetStatus(system.NewGetStatusParamsWithContext(ctx), nil)
	if err == nil {
		// TODO:
		_ = b.MQTTPublishAsync(ctx, MQTTPublishTopicStateUpTime.Format(snMQTT), status.Payload.DeviceUpTime)
		metricUpTime.With("serial_number", sn).Set(float64(status.Payload.DeviceUpTime))

		memoryUsage := int64(status.Payload.MemoryList[0].MemoryUsage) * MB
		// TODO:
		_ = b.MQTTPublishAsync(ctx, MQTTPublishTopicStateMemoryUsage.Format(snMQTT), memoryUsage)
		metricMemoryUsage.With("serial_number", sn).Set(float64(memoryUsage))

		memoryAvailable := int64(status.Payload.MemoryList[0].MemoryAvailable) * MB
		// TODO:
		_ = b.MQTTPublishAsync(ctx, MQTTPublishTopicStateMemoryAvailable.Format(snMQTT), memoryAvailable)
		metricMemoryAvailable.With("serial_number", sn).Set(float64(memoryAvailable))
	} else {
		b.Logger().Error("Request SystemStatus failed", "error", err.Error())
	}

	storage, err := b.client.ContentManager.GetStorage(content_manager.NewGetStorageParamsWithContext(ctx), nil)
	if err == nil {
		for _, hdd := range storage.Payload.HddList {
			if hdd.Name == "" {
				continue
			}

			// TODO:
			_ = b.MQTTPublishAsync(ctx, MQTTPublishTopicStateHDDCapacity.Format(snMQTT, hdd.ID), hdd.Capacity*MB)

			// TODO:
			_ = b.MQTTPublishAsync(ctx, MQTTPublishTopicStateHDDUsage.Format(snMQTT, hdd.ID), (hdd.Capacity-hdd.FreeSpace)*MB)
			metricStorageUsage.With("serial_number", sn).With("name", hdd.Name).Set(float64((hdd.Capacity - hdd.FreeSpace) * MB))

			// TODO:
			_ = b.MQTTPublishAsync(ctx, MQTTPublishTopicStateHDDFree.Format(snMQTT, hdd.ID), hdd.FreeSpace*MB)
			metricStorageAvailable.With("serial_number", sn).With("name", hdd.Name).Set(float64(hdd.FreeSpace * MB))
		}
	} else {
		b.Logger().Error("Request ContentManagementStorage failed", "error", err.Error())
	}

	return nil, nil
}
