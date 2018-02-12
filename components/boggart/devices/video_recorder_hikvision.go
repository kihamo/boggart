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
	metricVideoRecorderHikVisionMemoryUsage     = snitch.NewGauge(boggart.ComponentName+"_device_video_recorder_hikvision_memory_usage_bytes", "Memory usage in HikVision video recorder")
	metricVideoRecorderHikVisionMemoryAvailable = snitch.NewGauge(boggart.ComponentName+"_device_video_recorder_hikvision_memory_available_bytes", "Memory available in HikVision video recorder")
)

type VideoRecorderHikVision struct {
	boggart.DeviceBase

	isapi        *hikvision.ISAPI
	serialNumber string
}

func NewVideoRecorderHikVision(isapi *hikvision.ISAPI) (*VideoRecorderHikVision, error) {
	device := &VideoRecorderHikVision{
		isapi: isapi,
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
}

func (d *VideoRecorderHikVision) Collect(ch chan<- snitch.Metric) {
	if d.serialNumber == "" {
		return
	}

	ctx, ctxCancel := context.WithTimeout(context.Background(), time.Millisecond*300)
	defer ctxCancel()

	status, err := d.isapi.SystemStatus(ctx)
	if err != nil {
		log.Printf("Try to get system status for video recorder failed with error %s", err.Error())
		return
	}

	if len(status.Memory) == 0 {
		log.Print("Try to get system status for video recorder failed because response is wrong")
		return
	}

	memoryUsage := metricVideoRecorderHikVisionMemoryUsage.With("serial_number", d.serialNumber)
	memoryUsage.Set(status.Memory[0].MemoryUsage.Float64() * 1024 * 1024)
	memoryUsage.Collect(ch)

	memoryAvailable := metricVideoRecorderHikVisionMemoryAvailable.With("serial_number", d.serialNumber)
	memoryAvailable.Set(status.Memory[0].MemoryAvailable.Float64() * 1024 * 1024)
	memoryAvailable.Collect(ch)
}
