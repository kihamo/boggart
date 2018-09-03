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

	return []workers.Task{
		taskSerialNumber,
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

func (d *VideoRecorderHikVision) startAlertStreaming() error {
	ctx := context.Background()

	stream, err := d.isapi.EventNotificationAlertStream(ctx)
	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case event := <-stream.NextAlert():
				if !d.IsEnabled() || event.EventState != hikvision.EventEventStateActive {
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

			case _ = <-stream.NextError():
				// TODO: log errors

			case <-ctx.Done():
				return
			}
		}
	}()

	return nil
}
