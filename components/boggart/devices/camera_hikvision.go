package devices

import (
	"context"
	"log"
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/providers/hikvision"
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
}

func NewCameraHikVision(isapi *hikvision.ISAPI, channel uint64) (*HikVisionCamera, error) {
	device := &HikVisionCamera{
		isapi:   isapi,
		channel: channel,
	}
	device.DeviceBase.Init()

	ctx, ctxCancel := context.WithTimeout(context.Background(), time.Millisecond*300)
	defer ctxCancel()

	deviceInfo, err := isapi.SystemDeviceInfo(ctx)
	if err != nil {
		return nil, err
	}

	device.SetId(deviceInfo.DeviceID)
	device.SetDescription(deviceInfo.SerialNumber)
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

	ctx, ctxCancel := context.WithTimeout(context.Background(), time.Millisecond*300)
	defer ctxCancel()

	status, err := d.isapi.SystemStatus(ctx)
	if err != nil {
		log.Printf("Try to get system status for camera failed with error %s", err.Error())
		return
	}

	if len(status.Memory) == 0 {
		log.Print("Try to get system status for video recorder failed because response is wrong")
		return
	}

	memoryUsage := metricCameraHikVisionMemoryUsage.With("serial_number", d.serialNumber)
	memoryUsage.Set(status.Memory[0].MemoryUsage.Float64() * 1024 * 1024)
	memoryUsage.Collect(ch)

	memoryAvailable := metricCameraHikVisionMemoryAvailable.With("serial_number", d.serialNumber)
	memoryAvailable.Set(status.Memory[0].MemoryAvailable.Float64() * 1024 * 1024)
	memoryAvailable.Collect(ch)
}
