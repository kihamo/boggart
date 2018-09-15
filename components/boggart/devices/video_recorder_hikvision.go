package devices

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	m "github.com/eclipse/paho.mqtt.golang"
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/providers/hikvision"
	"github.com/kihamo/boggart/components/mqtt"
	"github.com/kihamo/go-workers"
	"github.com/kihamo/go-workers/task"
)

const (
	VideoRecorderHikVisionIgnoreInterval = time.Second * 5
	VideoRecorderMQTTTopicPrefix         = boggart.ComponentName + "/video-recorder/"
)

type VideoRecorderHikVision struct {
	boggart.DeviceBase
	boggart.DeviceSerialNumber

	isapi                 *hikvision.ISAPI
	interval              time.Duration
	mutex                 sync.Mutex
	alertStreamingHistory map[string]time.Time
	mqtt                  mqtt.Component
}

func NewVideoRecorderHikVision(isapi *hikvision.ISAPI, interval time.Duration, m mqtt.Component) *VideoRecorderHikVision {
	device := &VideoRecorderHikVision{
		isapi:                 isapi,
		interval:              interval,
		alertStreamingHistory: make(map[string]time.Time),
		mqtt: m,
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

	d.mqtt.Subscribe(d)
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

func (d *VideoRecorderHikVision) Filters() map[string]byte {
	sn := strings.Replace(d.SerialNumber(), "/", "_", -1)
	if sn == "" {
		return map[string]byte{}
	}

	return map[string]byte{
		VideoRecorderMQTTTopicPrefix + sn + "/ptz/#": 0,
	}
}

func (d *VideoRecorderHikVision) Callback(client mqtt.Component, message m.Message) {
	if !d.IsEnabled() {
		return
	}

	receivedChannelId := message.Topic()[len(VideoRecorderMQTTTopicPrefix+d.SerialNumber()+"/ptz/"):]
	parts := strings.Split(receivedChannelId, "/")

	if len(parts) < 2 {
		return
	}

	channelId, err := strconv.ParseUint(parts[0], 10, 64)
	if err != nil {
		return
	}

	switch strings.ToLower(parts[1]) {
	case "presets":
		presetId, err := strconv.ParseUint(string(message.Payload()), 10, 64)
		if err != nil {
			return
		}

		d.isapi.PTZPresetGoTo(context.Background(), channelId, presetId)

	case "relative":
		var request struct {
			X    int64 `xml:"x,omitempty"`
			Y    int64 `xml:"y,omitempty"`
			Zoom int64 `xml:"zoom,omitempty"`
		}

		if err := json.Unmarshal(message.Payload(), &request); err == nil {
			d.isapi.PTZRelative(context.Background(), channelId, request.X, request.Y, request.Zoom)
		}

	case "absolute":
		var request struct {
			Elevation int64  `json:"elevation,omitempty"`
			Azimuth   uint64 `json:"azimuth,omitempty"`
			Zoom      uint64 `json:"zoom,omitempty"`
		}

		if err := json.Unmarshal(message.Payload(), &request); err == nil {
			d.isapi.PTZAbsolute(context.Background(), channelId, request.Elevation, request.Azimuth, request.Zoom)
		}
	}
}
