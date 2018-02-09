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
	isapi        *hikvision.ISAPI
	id           string
	description  string
	serialNumber string
	enabled      bool
}

func NewVideoRecorderHikVision(isapi *hikvision.ISAPI) (*VideoRecorderHikVision, error) {
	device := &VideoRecorderHikVision{
		isapi: isapi,
	}

	ctx, ctxCancel := context.WithTimeout(context.Background(), time.Millisecond*300)
	defer ctxCancel()

	deviceInfo, err := isapi.SystemDeviceInfo(ctx)
	if err != nil {
		return nil, err
	}

	device.enabled = true
	device.id = deviceInfo.DeviceID
	device.description = deviceInfo.SerialNumber
	device.serialNumber = deviceInfo.SerialNumber

	return device, nil
}

func (d *VideoRecorderHikVision) Id() string {
	return d.id
}

func (d *VideoRecorderHikVision) Position() (int64, int64) {
	return 0, 0
}

func (d *VideoRecorderHikVision) Description() string {
	return d.description
}

func (d *VideoRecorderHikVision) IsEnabled() bool {
	return d.enabled
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
