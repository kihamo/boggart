package devices

import (
	"context"
	"errors"
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/providers/hikvision"
	"github.com/kihamo/go-workers"
	"github.com/kihamo/go-workers/task"
	"github.com/kihamo/snitch"
)

var (
	metricCameraHikVisionMemoryUsage     = snitch.NewGauge(boggart.ComponentName+"_device_camera_hikvision_memory_usage_bytes", "Memory usage in HikVision camera")
	metricCameraHikVisionMemoryAvailable = snitch.NewGauge(boggart.ComponentName+"_device_camera_hikvision_memory_available_bytes", "Memory available in HikVision camera")
)

type HikVisionCamera struct {
	boggart.DeviceWithSerialNumber

	isapi    *hikvision.ISAPI
	channel  uint64
	interval time.Duration
}

func NewCameraHikVision(isapi *hikvision.ISAPI, channel uint64, interval time.Duration) *HikVisionCamera {
	device := &HikVisionCamera{
		isapi:    isapi,
		channel:  channel,
		interval: interval,
	}
	device.Init()
	device.SetDescription("HikVision camera")

	return device
}

func (d *HikVisionCamera) Types() []boggart.DeviceType {
	return []boggart.DeviceType{
		boggart.DeviceTypeCamera,
	}
}

func (d *HikVisionCamera) Snapshot(ctx context.Context) ([]byte, error) {
	return d.isapi.StreamingPicture(ctx, d.channel)
}

func (d *HikVisionCamera) Describe(ch chan<- *snitch.Description) {
	serialNumber := d.SerialNumber()
	if serialNumber == "" {
		return
	}

	metricCameraHikVisionMemoryUsage.With("serial_number", serialNumber).Describe(ch)
	metricCameraHikVisionMemoryAvailable.With("serial_number", serialNumber).Describe(ch)
}

func (d *HikVisionCamera) Collect(ch chan<- snitch.Metric) {
	serialNumber := d.SerialNumber()
	if serialNumber == "" {
		return
	}

	metricCameraHikVisionMemoryUsage.With("serial_number", serialNumber).Collect(ch)
	metricCameraHikVisionMemoryAvailable.With("serial_number", serialNumber).Collect(ch)
}

func (d *HikVisionCamera) Ping(ctx context.Context) bool {
	_, err := d.isapi.SystemStatus(ctx)
	return err == nil
}

func (d *HikVisionCamera) Tasks() []workers.Task {
	taskSerialNumber := task.NewFunctionTillStopTask(d.taskSerialNumber)
	taskSerialNumber.SetTimeout(time.Second * 5)
	taskSerialNumber.SetRepeats(-1)
	taskSerialNumber.SetRepeatInterval(time.Minute)
	taskSerialNumber.SetName("device-camera-hikvision-serial-number")

	taskUpdater := task.NewFunctionTask(d.taskUpdater)
	taskUpdater.SetRepeats(-1)
	taskUpdater.SetRepeatInterval(d.interval)
	taskUpdater.SetName("device-camera-hikvision-updater-" + d.Id())

	return []workers.Task{
		taskSerialNumber,
		taskUpdater,
	}
}

func (d *HikVisionCamera) taskSerialNumber(ctx context.Context) (interface{}, error, bool) {
	if !d.IsEnabled() {
		return nil, nil, false
	}

	deviceInfo, err := d.isapi.SystemDeviceInfo(ctx)
	if err != nil {
		return nil, err, false
	}

	if deviceInfo.SerialNumber == "" {
		return nil, errors.New("Device returns empty serial number"), false
	}

	d.SetSerialNumber(deviceInfo.SerialNumber)
	d.SetDescription("HikVision camera with serial number " + deviceInfo.SerialNumber)

	return nil, nil, true
}

func (d *HikVisionCamera) taskUpdater(ctx context.Context) (interface{}, error) {
	if !d.IsEnabled() {
		return nil, nil
	}

	serialNumber := d.SerialNumber()
	if serialNumber == "" {
		return nil, nil
	}

	status, err := d.isapi.SystemStatus(ctx)
	if err != nil {
		return nil, err
	}

	if len(status.Memory) == 0 {
		return nil, errors.New("Try to get system status for video recorder failed because response is wrong")
	}

	metricCameraHikVisionMemoryUsage.With("serial_number", serialNumber).Set(status.Memory[0].MemoryUsage.Float64() * 1024 * 1024)
	metricCameraHikVisionMemoryAvailable.With("serial_number", serialNumber).Set(status.Memory[0].MemoryAvailable.Float64() * 1024 * 1024)

	return nil, nil
}
