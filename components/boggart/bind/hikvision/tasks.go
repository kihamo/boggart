package hikvision

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/mqtt"
	"github.com/kihamo/go-workers"
	"github.com/kihamo/go-workers/task"
)

func (d *Bind) Tasks() []workers.Task {
	taskLiveness := task.NewFunctionTask(d.taskLiveness)
	taskLiveness.SetTimeout(time.Second * 5)
	taskLiveness.SetRepeats(-1)
	taskLiveness.SetRepeatInterval(time.Minute)
	taskLiveness.SetName("bind-hikvision-liveness")

	taskPTZStatus := task.NewFunctionTillStopTask(d.taskPTZStatus)
	taskPTZStatus.SetTimeout(time.Second * 5)
	taskPTZStatus.SetRepeats(-1)
	taskPTZStatus.SetRepeatInterval(time.Minute)
	taskPTZStatus.SetName("bind-hikvision-ptz-status")

	taskState := task.NewFunctionTask(d.taskState)
	taskState.SetTimeout(time.Second * 30)
	taskState.SetRepeats(-1)
	taskState.SetRepeatInterval(time.Minute * 15)
	taskState.SetName("bind-hikvision-state")

	return []workers.Task{
		taskLiveness,
		taskPTZStatus,
		taskState,
	}
}

func (d *Bind) taskLiveness(ctx context.Context) (interface{}, error) {
	deviceInfo, err := d.isapi.SystemDeviceInfo(ctx)
	if err != nil {
		d.UpdateStatus(boggart.DeviceStatusOffline)
		return nil, err
	}

	if deviceInfo.SerialNumber == "" {
		d.UpdateStatus(boggart.DeviceStatusOffline)
		return nil, errors.New("device returns empty serial number")
	}

	if d.SerialNumber() == "" {
		ptzChannels := make(map[uint64]HikVisionPTZChannel, 0)
		if list, err := d.isapi.PTZChannels(ctx); err == nil {
			for _, channel := range list.Channels {
				ptzChannels[channel.ID] = HikVisionPTZChannel{
					Channel: channel,
				}
			}
		}

		d.mutex.Lock()
		d.ptzChannels = ptzChannels
		d.mutex.Unlock()

		//if err := d.startAlertStreaming(); err != nil {
		//	return nil, err, false
		//}

		d.SetSerialNumber(deviceInfo.SerialNumber)

		d.MQTTPublishAsync(ctx, MQTTTopicStateModel.Format(deviceInfo.SerialNumber), 0, true, deviceInfo.Model)
		d.MQTTPublishAsync(ctx, MQTTTopicStateFirmwareVersion.Format(deviceInfo.SerialNumber), 0, true, deviceInfo.FirmwareVersion)
		d.MQTTPublishAsync(ctx, MQTTTopicStateFirmwareReleasedDate.Format(deviceInfo.SerialNumber), 0, true, deviceInfo.FirmwareReleasedDate)
	}

	d.UpdateStatus(boggart.DeviceStatusOnline)

	return nil, nil
}

func (d *Bind) taskPTZStatus(ctx context.Context) (interface{}, error, bool) {
	if d.Status() != boggart.DeviceStatusOnline {
		return nil, nil, false
	}

	d.mutex.RLock()
	defer d.mutex.RUnlock()

	if d.ptzChannels == nil {
		return nil, nil, false
	}

	if len(d.ptzChannels) == 0 {
		return nil, nil, true
	}

	stop := true

	for id, channel := range d.ptzChannels {
		if !channel.Channel.Enabled {
			continue
		}

		if err := d.updateStatusByChannelId(ctx, id); err != nil {
			log.Printf("failed updated status for %s device in channel %d", d.SerialNumber(), id)
			continue
		}

		stop = false
	}

	return nil, nil, stop
}

func (d *Bind) taskState(ctx context.Context) (interface{}, error) {
	if d.Status() != boggart.DeviceStatusOnline {
		return nil, nil
	}

	d.mutex.RLock()
	defer d.mutex.RUnlock()

	status, err := d.isapi.SystemStatus(ctx)
	if err != nil {
		return nil, err
	}

	sn := mqtt.NameReplace(d.SerialNumber())

	d.MQTTPublishAsync(ctx, MQTTTopicStateUpTime.Format(sn), 1, false, status.DeviceUpTime)
	d.MQTTPublishAsync(ctx, MQTTTopicStateMemoryAvailable.Format(sn), 1, false, uint64(status.Memory[0].MemoryAvailable.Float64())*MB)
	d.MQTTPublishAsync(ctx, MQTTTopicStateMemoryUsage.Format(sn), 1, false, uint64(status.Memory[0].MemoryUsage.Float64())*MB)

	storage, err := d.isapi.ContentManagementStorage(ctx)
	if err != nil {
		return nil, err
	}

	for _, hdd := range storage.HDD {
		d.MQTTPublishAsync(ctx, MQTTTopicStateHDDCapacity.Format(sn, hdd.ID), 1, false, hdd.Capacity*MB)
		d.MQTTPublishAsync(ctx, MQTTTopicStateHDDFree.Format(sn, hdd.ID), 1, false, hdd.FreeSpace*MB)
		d.MQTTPublishAsync(ctx, MQTTTopicStateHDDUsage.Format(sn, hdd.ID), 1, false, (hdd.Capacity-hdd.FreeSpace)*MB)
	}

	return nil, nil
}
