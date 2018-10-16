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
	CameraHikVisionIgnoreInterval = time.Second * 5
	CameraMQTTTopicPrefix         = boggart.ComponentName + "/camera/"
)

type cameraHikVisionPTZChannel struct {
	Channel hikvision.PTZChannel
	Status  *hikvision.PTZStatus
}

type CameraHikVision struct {
	boggart.DeviceBase
	boggart.DeviceSerialNumber

	isapi                 *hikvision.ISAPI
	interval              time.Duration
	mutex                 sync.RWMutex
	alertStreamingHistory map[string]time.Time
	mqtt                  mqtt.Component

	ptzChannels map[uint64]cameraHikVisionPTZChannel
}

func NewCameraHikVision(isapi *hikvision.ISAPI, interval time.Duration, m mqtt.Component) *CameraHikVision {
	device := &CameraHikVision{
		isapi:                 isapi,
		interval:              interval,
		alertStreamingHistory: make(map[string]time.Time),
		mqtt: m,
	}
	device.Init()
	device.SetDescription("HikVision camera")

	return device
}

func (d *CameraHikVision) Types() []boggart.DeviceType {
	return []boggart.DeviceType{
		boggart.DeviceTypeCamera,
	}
}

func (d *CameraHikVision) Ping(ctx context.Context) bool {
	_, err := d.isapi.SystemStatus(ctx)
	return err == nil
}

func (d *CameraHikVision) Tasks() []workers.Task {
	taskSerialNumber := task.NewFunctionTillStopTask(d.taskSerialNumber)
	taskSerialNumber.SetTimeout(time.Second * 5)
	taskSerialNumber.SetRepeats(-1)
	taskSerialNumber.SetRepeatInterval(time.Minute)
	taskSerialNumber.SetName("device-camera-hikvision-serial-number")

	taskPTZStatus := task.NewFunctionTillStopTask(d.taskPTZStatus)
	taskPTZStatus.SetTimeout(time.Second * 5)
	taskPTZStatus.SetRepeats(-1)
	taskPTZStatus.SetRepeatInterval(time.Minute)
	taskPTZStatus.SetName("device-camera-hikvision-ptz-status")

	return []workers.Task{
		taskSerialNumber,
		taskPTZStatus,
	}
}

func (d *CameraHikVision) taskSerialNumber(ctx context.Context) (interface{}, error, bool) {
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

	ptzChannels := make(map[uint64]cameraHikVisionPTZChannel, 0)
	if list, err := d.isapi.PTZChannels(ctx); err == nil {
		for _, channel := range list.Channels {
			ptzChannels[channel.ID] = cameraHikVisionPTZChannel{
				Channel: channel,
			}
		}
	}

	d.mutex.Lock()
	d.ptzChannels = ptzChannels
	d.mutex.Unlock()

	if err := d.startAlertStreaming(); err != nil {
		return nil, err, false
	}

	d.mqtt.Subscribe(d)
	return nil, nil, true
}

func (d *CameraHikVision) taskPTZStatus(ctx context.Context) (interface{}, error, bool) {
	d.mutex.RLock()
	defer d.mutex.RUnlock()

	if !d.IsEnabled() || d.ptzChannels == nil {
		return nil, nil, false
	}

	if len(d.ptzChannels) == 0 {
		return nil, nil, true
	}

	stop := true

	for id, channel := range d.ptzChannels {
		if !channel.Channel.Enabled {
			continue
		}

		status, err := d.isapi.PTZStatus(ctx, channel.Channel.ID)
		if err != nil {
			// TODO: log
			continue
		}

		stop = false
		topicPrefix := CameraMQTTTopicPrefix + d.SerialNumber() + "/ptz/" + strconv.FormatUint(channel.Channel.ID, 10) + "/status/"

		if channel.Status == nil || channel.Status.AbsoluteHigh.Elevation != status.AbsoluteHigh.Elevation {
			d.mqtt.Publish(topicPrefix+"elevation", 1, false, strconv.FormatInt(status.AbsoluteHigh.Elevation, 10))
		}

		if channel.Status == nil || channel.Status.AbsoluteHigh.Azimuth != status.AbsoluteHigh.Azimuth {
			d.mqtt.Publish(topicPrefix+"azimuth", 1, false, strconv.FormatUint(status.AbsoluteHigh.Azimuth, 10))
		}

		if channel.Status == nil || channel.Status.AbsoluteHigh.AbsoluteZoom != status.AbsoluteHigh.AbsoluteZoom {
			d.mqtt.Publish(topicPrefix+"zoom", 1, false, strconv.FormatUint(status.AbsoluteHigh.AbsoluteZoom, 10))
		}

		channel.Status = &status
		d.ptzChannels[id] = channel
	}

	return nil, nil, stop
}

func (d *CameraHikVision) startAlertStreaming() error {
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

				if !ok || event.DateTime.Sub(lastFire) > CameraHikVisionIgnoreInterval {
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

func (d *CameraHikVision) Filters() map[string]byte {
	sn := strings.Replace(d.SerialNumber(), "/", "_", -1)
	if sn == "" {
		return map[string]byte{}
	}

	return map[string]byte{
		CameraMQTTTopicPrefix + sn + "/ptz/#": 0,
	}
}

func (d *CameraHikVision) Callback(ctx context.Context, client mqtt.Component, message m.Message) {
	if !d.IsEnabled() || d.ptzChannels == nil || len(d.ptzChannels) == 0 {
		return
	}

	receivedChannelId := message.Topic()[len(CameraMQTTTopicPrefix+d.SerialNumber()+"/ptz/"):]
	parts := strings.Split(receivedChannelId, "/")

	if len(parts) < 2 {
		return
	}

	var err error

	channelId, err := strconv.ParseUint(parts[0], 10, 64)
	if err != nil {
		return
	}

	_, ok := d.ptzChannels[channelId]
	if !ok {
		return
	}

	switch strings.ToLower(parts[1]) {
	case "move":
		switch string(message.Payload()) {
		case "stop":
			err = d.move(ctx, channelId, 0, 0, 0, 0)
		case "right":
			err = d.move(ctx, channelId, 1, 0, 0, 0)
		case "left":
			err = d.move(ctx, channelId, -1, 0, 0, 0)
		case "up":
			err = d.move(ctx, channelId, 0, 1, 0, 0)
		case "up-right":
			err = d.move(ctx, channelId, 1, 1, 0, 0)
		case "up-left":
			err = d.move(ctx, channelId, -1, 1, 0, 0)
		case "down":
			err = d.move(ctx, channelId, 0, -1, 0, 0)
		case "down-right":
			err = d.move(ctx, channelId, 1, -1, 0, 0)
		case "down-left":
			err = d.move(ctx, channelId, -1, -1, 0, 0)
		case "narrow":
			err = d.move(ctx, channelId, 0, 0, 1, 0)
		case "wide":
			err = d.move(ctx, channelId, 0, 0, -1, 0)
		}

	case "preset":
		presetId, err := strconv.ParseUint(string(message.Payload()), 10, 64)
		if err == nil {
			err = d.isapi.PTZPresetGoTo(ctx, channelId, presetId)
		}

	case "relative":
		var request struct {
			X    int64 `xml:"x,omitempty"`
			Y    int64 `xml:"y,omitempty"`
			Zoom int64 `xml:"zoom,omitempty"`
		}

		err = json.Unmarshal(message.Payload(), &request)
		if err == nil {
			err = d.isapi.PTZRelative(ctx, channelId, request.X, request.Y, request.Zoom)
		}

	case "absolute":
		var request struct {
			Elevation int64  `json:"elevation,omitempty"`
			Azimuth   uint64 `json:"azimuth,omitempty"`
			Zoom      uint64 `json:"zoom,omitempty"`
		}

		err = json.Unmarshal(message.Payload(), &request)
		if err == nil {
			err = d.isapi.PTZAbsolute(ctx, channelId, request.Elevation, request.Azimuth, request.Zoom)
		}

	case "continuous":
		var request struct {
			Pan  int64 `json:"pan,omitempty"`
			Tilt int64 `json:"tilt,omitempty"`
			Zoom int64 `json:"zoom,omitempty"`
		}

		err = json.Unmarshal(message.Payload(), &request)
		if err == nil {
			err = d.isapi.PTZContinuous(ctx, channelId, request.Pan, request.Tilt, request.Zoom)
		}

	case "momentary":
		var request struct {
			Pan      int64         `json:"pan,omitempty"`
			Tilt     int64         `json:"tilt,omitempty"`
			Zoom     int64         `json:"zoom,omitempty"`
			Duration time.Duration `json:"duration,omitempty"`
		}

		err = json.Unmarshal(message.Payload(), &request)
		if err == nil {
			duration := time.Duration(request.Duration) * time.Millisecond
			err = d.isapi.PTZMomentary(ctx, channelId, request.Pan, request.Tilt, request.Zoom, duration)
		}
	}

	if err == nil {
		_, err = d.Tasks()[1].Run(ctx)
	}

	if err != nil {
		fmt.Println(err.Error(), parts, string(message.Payload()))
	}
}

func (d *CameraHikVision) move(ctx context.Context, channelId uint64, panDirection, tiltDirection, zoomDirection int64, duration time.Duration) error {
	pan := panDirection
	tilt := tiltDirection
	zoom := zoomDirection

	if duration > 0 {
		return d.isapi.PTZMomentary(ctx, channelId, pan, tilt, zoom, duration)
	}

	return d.isapi.PTZContinuous(ctx, channelId, pan, tilt, zoom)
}
