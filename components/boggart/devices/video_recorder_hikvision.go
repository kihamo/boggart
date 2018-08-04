package devices

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/providers/hikvision"
	"github.com/kihamo/go-workers"
	"github.com/kihamo/go-workers/task"
)

/*
var (
	metricVideoRecorderHikVisionMemoryUsage      = snitch.NewGauge(boggart.ComponentName+"_device_video_recorder_hikvision_memory_usage_bytes", "Memory usage in HikVision video recorder")
	metricVideoRecorderHikVisionMemoryAvailable  = snitch.NewGauge(boggart.ComponentName+"_device_video_recorder_hikvision_memory_available_bytes", "Memory available in HikVision video recorder")
	metricVideoRecorderHikVisionStorageUsage     = snitch.NewGauge(boggart.ComponentName+"_device_video_recorder_hikvision_storage_usage_bytes", "Storage usage in HikVision video recorder")
	metricVideoRecorderHikVisionStorageAvailable = snitch.NewGauge(boggart.ComponentName+"_device_video_recorder_hikvision_storage_available_bytes", "Storage available in HikVision video recorder")
)
*/
const (
	VideoRecorderHikVisionIgnoreInterval = time.Second * 5
)

type VideoRecorderHikVision struct {
	boggart.DeviceBase
	boggart.DeviceSerialNumber

	isapi                 *hikvision.ISAPI
	interval              time.Duration
	mutex                 sync.Mutex
	alertStreamingHistory map[string]time.Time
}

func NewVideoRecorderHikVision(isapi *hikvision.ISAPI, interval time.Duration) *VideoRecorderHikVision {
	device := &VideoRecorderHikVision{
		isapi:                 isapi,
		interval:              interval,
		alertStreamingHistory: make(map[string]time.Time),
	}
	device.Init()
	device.SetDescription("HikVision video recorder")

	return device
}

func (d *VideoRecorderHikVision) Types() []boggart.DeviceType {
	return []boggart.DeviceType{
		boggart.DeviceTypeVideoRecorder,
	}
}

/*
func (d *VideoRecorderHikVision) Describe(ch chan<- *snitch.Description) {
	serialNumber := d.SerialNumber()
	if serialNumber == "" {
		return
	}

	metricVideoRecorderHikVisionMemoryUsage.With("serial_number", serialNumber).Describe(ch)
	metricVideoRecorderHikVisionMemoryAvailable.With("serial_number", serialNumber).Describe(ch)
	metricVideoRecorderHikVisionStorageUsage.With("serial_number", serialNumber).Describe(ch)
	metricVideoRecorderHikVisionStorageAvailable.With("serial_number", serialNumber).Describe(ch)
}

func (d *VideoRecorderHikVision) Collect(ch chan<- snitch.Metric) {
	serialNumber := d.SerialNumber()
	if serialNumber == "" {
		return
	}

	metricVideoRecorderHikVisionMemoryUsage.With("serial_number", serialNumber).Collect(ch)
	metricVideoRecorderHikVisionMemoryAvailable.With("serial_number", serialNumber).Collect(ch)
	metricVideoRecorderHikVisionStorageUsage.With("serial_number", serialNumber).Collect(ch)
	metricVideoRecorderHikVisionStorageAvailable.With("serial_number", serialNumber).Collect(ch)
}
*/
func (d *VideoRecorderHikVision) Ping(ctx context.Context) bool {
	_, err := d.isapi.SystemStatus(ctx)
	return err == nil
}

func (d *VideoRecorderHikVision) Tasks() []workers.Task {
	taskSerialNumber := task.NewFunctionTillStopTask(d.taskSerialNumber)
	taskSerialNumber.SetTimeout(time.Second * 5)
	taskSerialNumber.SetRepeats(-1)
	taskSerialNumber.SetRepeatInterval(time.Minute)
	taskSerialNumber.SetName("device-video-recorder-hikvision-serial-number")
	/*
		taskUpdater := task.NewFunctionTask(d.taskUpdater)
		taskUpdater.SetRepeats(-1)
		taskUpdater.SetRepeatInterval(d.interval)
		taskUpdater.SetName("device-video-recorder-hikvision-updater-" + d.Id())
	*/
	return []workers.Task{
		taskSerialNumber,
		//taskUpdater,
	}
}

func (d *VideoRecorderHikVision) taskSerialNumber(ctx context.Context) (interface{}, error, bool) {
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
	d.SetDescription(d.Description() + " with serial number " + deviceInfo.SerialNumber)

	if err := d.startAlertStreaming(); err != nil {
		return nil, err, false
	}

	return nil, nil, true
}

/*
func (d *VideoRecorderHikVision) taskUpdater(ctx context.Context) (interface{}, error) {
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

	metricVideoRecorderHikVisionMemoryUsage.With("serial_number", serialNumber).Set(status.Memory[0].MemoryUsage.Float64() * 1024 * 1024)
	metricVideoRecorderHikVisionMemoryAvailable.With("serial_number", serialNumber).Set(status.Memory[0].MemoryAvailable.Float64() * 1024 * 1024)

	storage, err := d.isapi.ContentManagementStorage(ctx)
	if err != nil {
		return nil, err
	}

	if len(storage.HDD) > 0 {
		for _, hdd := range storage.HDD {
			// TODO: send event if status isn't OK

			metricVideoRecorderHikVisionStorageUsage.With("serial_number", serialNumber).With(
				"id", strconv.FormatUint(hdd.ID, 10),
				"name", hdd.Name,
				"type", "hdd",
			).Set(float64((hdd.Capacity - hdd.FreeSpace) * 1024 * 1024))

			metricVideoRecorderHikVisionStorageAvailable.With("serial_number", serialNumber).With(
				"id", strconv.FormatUint(hdd.ID, 10),
				"name", hdd.Name,
				"type", "hdd",
			).Set(float64(hdd.FreeSpace * 1024 * 1024))
		}
	}

	return nil, nil
}
*/
func (d *VideoRecorderHikVision) startAlertStreaming() error {
	stream, err := d.isapi.EventNotificationAlertStream(context.Background())
	if err != nil {
		return err
	}

	go func() {
		for {
			event, err := stream.NextAlert()
			if err != nil || !d.IsEnabled() || event.EventState != hikvision.EventEventStateActive {
				continue
			}

			id := fmt.Sprintf("%d-%s", event.DynChannelID, event.EventType)

			d.mutex.Lock()
			lastFire, ok := d.alertStreamingHistory[id]
			d.alertStreamingHistory[id] = event.DateTime
			d.mutex.Unlock()

			if !ok || event.DateTime.Sub(lastFire) > VideoRecorderHikVisionIgnoreInterval {
				d.TriggerEvent(boggart.DeviceEventHikvisionEventNotificationAlert, event, d.SerialNumber())
			}
		}
	}()

	return nil
}

/*
func (d *VideoRecorderHikVision) Snapshot(ctx context.Context, streaming uint64, input uint64) ([]byte, error) {
	return d.isapi.StreamingPicture(ctx, streaming*100+input)
}
*/
