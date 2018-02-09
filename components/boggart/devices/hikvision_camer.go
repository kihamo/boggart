package devices

import (
	"context"
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/providers/hikvision"
	"github.com/kihamo/snitch"
	"github.com/pborman/uuid"
)

var (
	metricHikVisionCameraMemoryUsage     = snitch.NewGauge(boggart.ComponentName+"_device_camera_hikvision_memory_usage_bytes", "Memory usage in HikVision camera")
	metricHikVisionCameraMemoryAvailable = snitch.NewGauge(boggart.ComponentName+"_device_camera_hikvision_memory_available_bytes", "Memory available in HikVision camera")
)

type HikVisionCamera struct {
	isapi        *hikvision.ISAPI
	channel      uint64
	id           string
	description  string
	serialNumber string
	enabled      bool
}

func NewHikVisionCamera(isapi *hikvision.ISAPI, channel uint64) *HikVisionCamera {
	device := &HikVisionCamera{
		isapi:   isapi,
		channel: channel,
	}

	ctx, ctxCancel := context.WithTimeout(context.Background(), time.Millisecond*300)
	defer ctxCancel()

	if deviceInfo, err := isapi.SystemDeviceInfo(ctx); err == nil {
		device.enabled = true
		device.id = deviceInfo.DeviceID
		device.description = deviceInfo.SerialNumber
		device.serialNumber = deviceInfo.SerialNumber
	} else {
		device.enabled = false
		device.id = uuid.New()
	}

	return device
}

func (c *HikVisionCamera) Id() string {
	return c.id
}

func (c *HikVisionCamera) Position() (int64, int64) {
	return 0, 0
}

func (c *HikVisionCamera) Description() string {
	return c.description
}

func (c *HikVisionCamera) IsEnabled() bool {
	return c.enabled
}

func (c *HikVisionCamera) Snapshot(ctx context.Context) ([]byte, error) {
	return c.isapi.StreamingPicture(ctx, c.channel)
}

func (c *HikVisionCamera) Describe(ch chan<- *snitch.Description) {
	if c.serialNumber == "" {
		return
	}

	metricHikVisionCameraMemoryUsage.With("serial_number", c.serialNumber).Describe(ch)
	metricHikVisionCameraMemoryAvailable.With("serial_number", c.serialNumber).Describe(ch)
}

func (c *HikVisionCamera) Collect(ch chan<- snitch.Metric) {
	if c.serialNumber == "" {
		return
	}

	ctx, ctxCancel := context.WithTimeout(context.Background(), time.Millisecond*300)
	defer ctxCancel()

	status, err := c.isapi.SystemStatus(ctx)
	if err != nil || len(status.Memory) == 0 {
		// FIXME: log
		return
	}

	memoryUsage := metricHikVisionCameraMemoryUsage.With("serial_number", c.serialNumber)
	memoryUsage.Set(status.Memory[0].MemoryUsage * 1024 * 1024)
	memoryUsage.Collect(ch)

	memoryAvailable := metricHikVisionCameraMemoryAvailable.With("serial_number", c.serialNumber)
	memoryAvailable.Set(status.Memory[0].MemoryAvailable * 1024 * 1024)
	memoryAvailable.Collect(ch)
}
