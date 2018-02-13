package devices

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/providers/hikvision"
	"github.com/kihamo/go-workers"
	"github.com/kihamo/go-workers/task"
	"github.com/kihamo/snitch"
)

var (
	metricVideoRecorderHikVisionMemoryUsage      = snitch.NewGauge(boggart.ComponentName+"_device_video_recorder_hikvision_memory_usage_bytes", "Memory usage in HikVision video recorder")
	metricVideoRecorderHikVisionMemoryAvailable  = snitch.NewGauge(boggart.ComponentName+"_device_video_recorder_hikvision_memory_available_bytes", "Memory available in HikVision video recorder")
	metricVideoRecorderHikVisionStorageUsage     = snitch.NewGauge(boggart.ComponentName+"_device_video_recorder_hikvision_storage_usage_bytes", "Memory usage in HikVision video recorder")
	metricVideoRecorderHikVisionStorageAvailable = snitch.NewGauge(boggart.ComponentName+"_device_video_recorder_hikvision_storage_available_bytes", "Memory available in HikVision video recorder")
)

type VideoRecorderHikVision struct {
	boggart.DeviceBase

	isapi        *hikvision.ISAPI
	serialNumber string
	interval     time.Duration
}

func NewVideoRecorderHikVision(isapi *hikvision.ISAPI, interval time.Duration) (*VideoRecorderHikVision, error) {
	device := &VideoRecorderHikVision{
		isapi:    isapi,
		interval: interval,
	}
	device.DeviceBase.Init()

	ctx, ctxCancel := context.WithTimeout(context.Background(), time.Millisecond*300)
	defer ctxCancel()

	deviceInfo, err := isapi.SystemDeviceInfo(ctx)
	if err != nil {
		return nil, err
	}

	device.SetId(deviceInfo.DeviceID)
	device.SetDescription("HikVision video recorder with serial number " + deviceInfo.SerialNumber)
	device.serialNumber = deviceInfo.SerialNumber

	return device, nil
}

func (d *VideoRecorderHikVision) Types() []boggart.DeviceType {
	return []boggart.DeviceType{
		boggart.DeviceTypeVideoRecorder,
	}
}

func (d *VideoRecorderHikVision) Describe(ch chan<- *snitch.Description) {
	if d.serialNumber == "" {
		return
	}

	metricVideoRecorderHikVisionMemoryUsage.With("serial_number", d.serialNumber).Describe(ch)
	metricVideoRecorderHikVisionMemoryAvailable.With("serial_number", d.serialNumber).Describe(ch)
	metricVideoRecorderHikVisionStorageUsage.With("serial_number", d.serialNumber).Describe(ch)
	metricVideoRecorderHikVisionStorageAvailable.With("serial_number", d.serialNumber).Describe(ch)
}

func (d *VideoRecorderHikVision) Collect(ch chan<- snitch.Metric) {
	if d.serialNumber == "" {
		return
	}

	metricVideoRecorderHikVisionMemoryUsage.With("serial_number", d.serialNumber).Collect(ch)
	metricVideoRecorderHikVisionMemoryAvailable.With("serial_number", d.serialNumber).Collect(ch)
	metricVideoRecorderHikVisionStorageUsage.With("serial_number", d.serialNumber).Collect(ch)
	metricVideoRecorderHikVisionStorageAvailable.With("serial_number", d.serialNumber).Collect(ch)
}

func (d *VideoRecorderHikVision) Tasks() []workers.Task {
	taskUpdater := task.NewFunctionTask(d.updater)
	taskUpdater.SetRepeats(-1)
	taskUpdater.SetRepeatInterval(d.interval)
	taskUpdater.SetName("device-video-recorder-hikvision-updater-" + d.serialNumber)

	return []workers.Task{
		taskUpdater,
	}
}

func (d *VideoRecorderHikVision) updater(ctx context.Context) (interface{}, error) {
	if !d.IsEnabled() {
		return nil, nil
	}

	status, err := d.isapi.SystemStatus(ctx)
	if err != nil {
		return nil, err
	}

	if len(status.Memory) == 0 {
		return nil, errors.New("Try to get system status for video recorder failed because response is wrong")
	}

	metricVideoRecorderHikVisionMemoryUsage.With("serial_number", d.serialNumber).Set(status.Memory[0].MemoryUsage.Float64() * 1024 * 1024)
	metricVideoRecorderHikVisionMemoryAvailable.With("serial_number", d.serialNumber).Set(status.Memory[0].MemoryAvailable.Float64() * 1024 * 1024)

	storage, err := d.isapi.ContentManagementStorage(ctx)
	if err != nil {
		return nil, err
	}

	if len(storage.HDD) > 0 {
		for _, hdd := range storage.HDD {
			metricVideoRecorderHikVisionStorageUsage.With("serial_number", d.serialNumber).With(
				"id", strconv.FormatUint(hdd.ID, 10),
				"name", hdd.Name,
				"type", "hdd",
			).Set(float64((hdd.Capacity - hdd.FreeSpace) * 1024 * 1024))

			metricVideoRecorderHikVisionStorageAvailable.With("serial_number", d.serialNumber).With(
				"id", strconv.FormatUint(hdd.ID, 10),
				"name", hdd.Name,
				"type", "hdd",
			).Set(float64(hdd.FreeSpace * 1024 * 1024))
		}
	}

	return nil, nil
}
