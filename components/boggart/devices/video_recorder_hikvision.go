package devices

import (
	"context"
	"errors"
	"strconv"
	"sync"
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
	metricVideoRecorderHikVisionStorageUsage     = snitch.NewGauge(boggart.ComponentName+"_device_video_recorder_hikvision_storage_usage_bytes", "Storage usage in HikVision video recorder")
	metricVideoRecorderHikVisionStorageAvailable = snitch.NewGauge(boggart.ComponentName+"_device_video_recorder_hikvision_storage_available_bytes", "Storage available in HikVision video recorder")
)

type VideoRecorderHikVision struct {
	boggart.DeviceBase

	mutex        sync.RWMutex
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
	device.SetDescription("HikVision video recorder")

	return device, nil
}

func (d *VideoRecorderHikVision) Types() []boggart.DeviceType {
	return []boggart.DeviceType{
		boggart.DeviceTypeVideoRecorder,
	}
}

func (d *VideoRecorderHikVision) Describe(ch chan<- *snitch.Description) {
	if d.SerialNumber() == "" {
		return
	}

	metricVideoRecorderHikVisionMemoryUsage.With("serial_number", d.SerialNumber()).Describe(ch)
	metricVideoRecorderHikVisionMemoryAvailable.With("serial_number", d.SerialNumber()).Describe(ch)
	metricVideoRecorderHikVisionStorageUsage.With("serial_number", d.SerialNumber()).Describe(ch)
	metricVideoRecorderHikVisionStorageAvailable.With("serial_number", d.SerialNumber()).Describe(ch)
}

func (d *VideoRecorderHikVision) Collect(ch chan<- snitch.Metric) {
	if d.SerialNumber() == "" {
		return
	}

	metricVideoRecorderHikVisionMemoryUsage.With("serial_number", d.SerialNumber()).Collect(ch)
	metricVideoRecorderHikVisionMemoryAvailable.With("serial_number", d.SerialNumber()).Collect(ch)
	metricVideoRecorderHikVisionStorageUsage.With("serial_number", d.SerialNumber()).Collect(ch)
	metricVideoRecorderHikVisionStorageAvailable.With("serial_number", d.SerialNumber()).Collect(ch)
}

func (d *VideoRecorderHikVision) Ping(ctx context.Context) bool {
	_, err := d.isapi.SystemStatus(ctx)
	return err == nil
}

func (d *VideoRecorderHikVision) SerialNumber() string {
	d.mutex.RLock()
	defer d.mutex.RUnlock()

	return d.serialNumber
}

func (d *VideoRecorderHikVision) Tasks() []workers.Task {
	taskSerialNumber := task.NewFunctionTillSuccessTask(d.taskSerialNumber)
	taskSerialNumber.SetTimeout(time.Second * 5)
	taskSerialNumber.SetRepeats(-1)
	taskSerialNumber.SetRepeatInterval(time.Minute)
	taskSerialNumber.SetName("device-video-recorder-hikvision-serial-number")

	taskUpdater := task.NewFunctionTask(d.updater)
	taskUpdater.SetRepeats(-1)
	taskUpdater.SetRepeatInterval(d.interval)
	taskUpdater.SetName("device-video-recorder-hikvision-updater-" + d.Id())

	return []workers.Task{
		taskSerialNumber,
		taskUpdater,
	}
}

func (d *VideoRecorderHikVision) taskSerialNumber(ctx context.Context) (interface{}, error) {
	if !d.IsEnabled() {
		return nil, errors.New("Device is disabled")
	}

	deviceInfo, err := d.isapi.SystemDeviceInfo(ctx)
	if err != nil {
		return nil, err
	}

	if deviceInfo.SerialNumber == "" {
		return nil, errors.New("Device returns empty serial number")
	}

	d.SetDescription("HikVision video recorder with serial number " + deviceInfo.SerialNumber)

	d.mutex.Lock()
	d.serialNumber = deviceInfo.SerialNumber
	d.mutex.Unlock()

	return nil, nil
}

func (d *VideoRecorderHikVision) updater(ctx context.Context) (interface{}, error) {
	if !d.IsEnabled() || d.SerialNumber() == "" {
		return nil, nil
	}

	status, err := d.isapi.SystemStatus(ctx)
	if err != nil {
		return nil, err
	}

	if len(status.Memory) == 0 {
		return nil, errors.New("Try to get system status for video recorder failed because response is wrong")
	}

	metricVideoRecorderHikVisionMemoryUsage.With("serial_number", d.SerialNumber()).Set(status.Memory[0].MemoryUsage.Float64() * 1024 * 1024)
	metricVideoRecorderHikVisionMemoryAvailable.With("serial_number", d.SerialNumber()).Set(status.Memory[0].MemoryAvailable.Float64() * 1024 * 1024)

	storage, err := d.isapi.ContentManagementStorage(ctx)
	if err != nil {
		return nil, err
	}

	if len(storage.HDD) > 0 {
		for _, hdd := range storage.HDD {
			// TODO: send event if status isn't OK

			metricVideoRecorderHikVisionStorageUsage.With("serial_number", d.SerialNumber()).With(
				"id", strconv.FormatUint(hdd.ID, 10),
				"name", hdd.Name,
				"type", "hdd",
			).Set(float64((hdd.Capacity - hdd.FreeSpace) * 1024 * 1024))

			metricVideoRecorderHikVisionStorageAvailable.With("serial_number", d.SerialNumber()).With(
				"id", strconv.FormatUint(hdd.ID, 10),
				"name", hdd.Name,
				"type", "hdd",
			).Set(float64(hdd.FreeSpace * 1024 * 1024))
		}
	}

	return nil, nil
}
