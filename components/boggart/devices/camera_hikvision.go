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
	boggart.DeviceBase

	isapi        *hikvision.ISAPI
	channel      uint64
	serialNumber string
	interval     time.Duration
}

func NewCameraHikVision(isapi *hikvision.ISAPI, channel uint64, interval time.Duration) (*HikVisionCamera, error) {
	device := &HikVisionCamera{
		isapi:    isapi,
		channel:  channel,
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
	device.SetDescription("HikVision camera with serial number " + deviceInfo.SerialNumber)
	device.serialNumber = deviceInfo.SerialNumber

	return device, nil
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
	if d.serialNumber == "" {
		return
	}

	metricCameraHikVisionMemoryUsage.With("serial_number", d.serialNumber).Describe(ch)
	metricCameraHikVisionMemoryAvailable.With("serial_number", d.serialNumber).Describe(ch)
}

func (d *HikVisionCamera) Collect(ch chan<- snitch.Metric) {
	if d.serialNumber == "" {
		return
	}

	metricCameraHikVisionMemoryUsage.With("serial_number", d.serialNumber).Collect(ch)
	metricCameraHikVisionMemoryAvailable.With("serial_number", d.serialNumber).Collect(ch)
}

func (d *HikVisionCamera) Ping(ctx context.Context) bool {
	_, err := d.isapi.SystemStatus(ctx)
	return err == nil
}

func (d *HikVisionCamera) Tasks() []workers.Task {
	taskUpdater := task.NewFunctionTask(d.updater)
	taskUpdater.SetRepeats(-1)
	taskUpdater.SetRepeatInterval(d.interval)
	taskUpdater.SetName("device-camera-hikvision-updater-" + d.serialNumber)

	return []workers.Task{
		taskUpdater,
	}
}

func (d *HikVisionCamera) updater(ctx context.Context) (interface{}, error) {
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

	metricCameraHikVisionMemoryUsage.With("serial_number", d.serialNumber).Set(status.Memory[0].MemoryUsage.Float64() * 1024 * 1024)
	metricCameraHikVisionMemoryAvailable.With("serial_number", d.serialNumber).Set(status.Memory[0].MemoryAvailable.Float64() * 1024 * 1024)

	return nil, nil
}
