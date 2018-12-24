package devices

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/providers/hikvision"
	"github.com/kihamo/boggart/components/mqtt"
	"github.com/kihamo/go-workers"
	"github.com/kihamo/go-workers/task"
)

const (
	MB                            uint64 = 1024 * 1024
	CameraHikVisionIgnoreInterval        = time.Second * 5

	CameraHikVisionMQTTTopicEvent                mqtt.Topic = boggart.ComponentName + "/cctv/+/+/+"
	CameraHikVisionMQTTTopicPTZMove              mqtt.Topic = boggart.ComponentName + "/cctv/+/ptz/+/move"
	CameraHikVisionMQTTTopicPTZAbsolute          mqtt.Topic = boggart.ComponentName + "/cctv/+/ptz/+/absolute"
	CameraHikVisionMQTTTopicPTZContinuous        mqtt.Topic = boggart.ComponentName + "/cctv/+/ptz/+/continuous"
	CameraHikVisionMQTTTopicPTZRelative          mqtt.Topic = boggart.ComponentName + "/cctv/+/ptz/+/relative"
	CameraHikVisionMQTTTopicPTZPreset            mqtt.Topic = boggart.ComponentName + "/cctv/+/ptz/+/preset"
	CameraHikVisionMQTTTopicPTZMomentary         mqtt.Topic = boggart.ComponentName + "/cctv/+/ptz/+/momentary"
	CameraHikVisionMQTTTopicPTZStatusElevation   mqtt.Topic = boggart.ComponentName + "/cctv/+/ptz/+/status/elevation"
	CameraHikVisionMQTTTopicPTZStatusAzimuth     mqtt.Topic = boggart.ComponentName + "/cctv/+/ptz/+/status/azimuth"
	CameraHikVisionMQTTTopicPTZStatusZoom        mqtt.Topic = boggart.ComponentName + "/cctv/+/ptz/+/status/zoom"
	CameraHikVisionMQTTTopicStateUpTime          mqtt.Topic = boggart.ComponentName + "/cctv/+/state/uptime"
	CameraHikVisionMQTTTopicStateMemoryUsage     mqtt.Topic = boggart.ComponentName + "/cctv/+/state/memory/usage"
	CameraHikVisionMQTTTopicStateMemoryAvailable mqtt.Topic = boggart.ComponentName + "/cctv/+/state/memory/available"
	CameraHikVisionMQTTTopicStateHDDCapacity     mqtt.Topic = boggart.ComponentName + "/cctv/+/state/hdd/+/capacity"
	CameraHikVisionMQTTTopicStateHDDFree         mqtt.Topic = boggart.ComponentName + "/cctv/+/state/hdd/+/free"
	CameraHikVisionMQTTTopicStateHDDUsage        mqtt.Topic = boggart.ComponentName + "/cctv/+/state/hdd/+/usage"
)

type cameraHikVisionPTZChannel struct {
	Channel hikvision.PTZChannel
	Status  *hikvision.PTZStatus
}

type CameraHikVision struct {
	boggart.DeviceBase
	boggart.DeviceSerialNumber
	boggart.DeviceMQTT

	mutex    sync.RWMutex
	initOnce sync.Once

	isapi                 *hikvision.ISAPI
	interval              time.Duration
	alertStreamingHistory map[string]time.Time

	ptzChannels map[uint64]cameraHikVisionPTZChannel
}

func NewCameraHikVision(isapi *hikvision.ISAPI, interval time.Duration) *CameraHikVision {
	device := &CameraHikVision{
		isapi:                 isapi,
		interval:              interval,
		alertStreamingHistory: make(map[string]time.Time),
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

func (d *CameraHikVision) Tasks() []workers.Task {
	taskLiveness := task.NewFunctionTask(d.taskLiveness)
	taskLiveness.SetTimeout(time.Second * 5)
	taskLiveness.SetRepeats(-1)
	taskLiveness.SetRepeatInterval(time.Minute)
	taskLiveness.SetName("device-camera-hikvision-liveness")

	taskPTZStatus := task.NewFunctionTillStopTask(d.taskPTZStatus)
	taskPTZStatus.SetTimeout(time.Second * 5)
	taskPTZStatus.SetRepeats(-1)
	taskPTZStatus.SetRepeatInterval(time.Minute)
	taskPTZStatus.SetName("device-camera-hikvision-ptz-status")

	taskState := task.NewFunctionTask(d.taskState)
	taskState.SetTimeout(time.Second * 30)
	taskState.SetRepeats(-1)
	taskState.SetRepeatInterval(time.Minute * 15)
	taskState.SetName("device-camera-hikvision-state")

	return []workers.Task{
		taskLiveness,
		taskPTZStatus,
		taskState,
	}
}

func (d *CameraHikVision) taskLiveness(ctx context.Context) (interface{}, error) {
	deviceInfo, err := d.isapi.SystemDeviceInfo(ctx)
	if err != nil {
		d.UpdateStatus(boggart.DeviceStatusOffline)
		return nil, err
	}

	if deviceInfo.SerialNumber == "" {
		d.UpdateStatus(boggart.DeviceStatusOffline)
		return nil, errors.New("device returns empty serial number")
	}

	d.UpdateStatus(boggart.DeviceStatusOnline)
	if d.SerialNumber() == "" {
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

		//if err := d.startAlertStreaming(); err != nil {
		//	return nil, err, false
		//}

		d.SetSerialNumber(deviceInfo.SerialNumber)
		d.initOnce.Do(func() {
			sn := d.SerialNumberMQTTEscaped()

			d.MQTTSubscribe(CameraHikVisionMQTTTopicPTZMove.Format(sn), 0, boggart.WrapMQTTSubscribeDeviceIsOnline(d, d.callbackMQTTAbsolute))
			d.MQTTSubscribe(CameraHikVisionMQTTTopicPTZAbsolute.Format(sn), 0, boggart.WrapMQTTSubscribeDeviceIsOnline(d, d.callbackMQTTAbsolute))
			d.MQTTSubscribe(CameraHikVisionMQTTTopicPTZContinuous.Format(sn), 0, boggart.WrapMQTTSubscribeDeviceIsOnline(d, d.callbackMQTTContinuous))
			d.MQTTSubscribe(CameraHikVisionMQTTTopicPTZRelative.Format(sn), 0, boggart.WrapMQTTSubscribeDeviceIsOnline(d, d.callbackMQTTRelative))
			d.MQTTSubscribe(CameraHikVisionMQTTTopicPTZPreset.Format(sn), 0, boggart.WrapMQTTSubscribeDeviceIsOnline(d, d.callbackMQTTPreset))
			d.MQTTSubscribe(CameraHikVisionMQTTTopicPTZMomentary.Format(sn), 0, boggart.WrapMQTTSubscribeDeviceIsOnline(d, d.callbackMQTTMomentary))
		})
	}

	return nil, nil
}

func (d *CameraHikVision) taskPTZStatus(ctx context.Context) (interface{}, error, bool) {
	if d.Status() != boggart.DeviceStatusOnline {
		return nil, nil, false
	}

	d.mutex.RLock()
	defer d.mutex.RUnlock()

	if d.ptzChannels == nil {
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

		if err := d.updateStatusByChannelId(ctx, id); err != nil {
			log.Printf("failed updated status for %s device in channel %d", d.SerialNumber(), id)
			continue
		}

		stop = false
	}

	return nil, nil, stop
}

func (d *CameraHikVision) taskState(ctx context.Context) (interface{}, error) {
	if d.Status() != boggart.DeviceStatusOnline {
		return nil, nil
	}

	d.mutex.RLock()
	defer d.mutex.RUnlock()

	status, err := d.isapi.SystemStatus(ctx)
	if err != nil {
		return nil, err
	}

	sn := d.SerialNumberMQTTEscaped()

	d.MQTTPublishAsync(ctx, CameraHikVisionMQTTTopicStateUpTime.Format(sn), 1, false, status.DeviceUpTime)
	d.MQTTPublishAsync(ctx, CameraHikVisionMQTTTopicStateMemoryAvailable.Format(sn), 1, false, uint64(status.Memory[0].MemoryAvailable.Float64())*MB)
	d.MQTTPublishAsync(ctx, CameraHikVisionMQTTTopicStateMemoryUsage.Format(sn), 1, false, uint64(status.Memory[0].MemoryUsage.Float64())*MB)

	storage, err := d.isapi.ContentManagementStorage(ctx)
	if err != nil {
		return nil, err
	}

	for _, hdd := range storage.HDD {
		d.MQTTPublishAsync(ctx, CameraHikVisionMQTTTopicStateHDDCapacity.Format(sn, hdd.ID), 1, false, hdd.Capacity*MB)
		d.MQTTPublishAsync(ctx, CameraHikVisionMQTTTopicStateHDDFree.Format(sn, hdd.ID), 1, false, hdd.FreeSpace*MB)
		d.MQTTPublishAsync(ctx, CameraHikVisionMQTTTopicStateHDDUsage.Format(sn, hdd.ID), 1, false, (hdd.Capacity-hdd.FreeSpace)*MB)
	}

	return nil, nil
}

func (d *CameraHikVision) MQTTTopics() []mqtt.Topic {
	return []mqtt.Topic{
		CameraHikVisionMQTTTopicEvent,
		CameraHikVisionMQTTTopicPTZMove,
		CameraHikVisionMQTTTopicPTZAbsolute,
		CameraHikVisionMQTTTopicPTZContinuous,
		CameraHikVisionMQTTTopicPTZRelative,
		CameraHikVisionMQTTTopicPTZPreset,
		CameraHikVisionMQTTTopicPTZMomentary,
		CameraHikVisionMQTTTopicPTZStatusElevation,
		CameraHikVisionMQTTTopicPTZStatusAzimuth,
		CameraHikVisionMQTTTopicPTZStatusZoom,
		CameraHikVisionMQTTTopicStateUpTime,
		CameraHikVisionMQTTTopicStateMemoryUsage,
		CameraHikVisionMQTTTopicStateMemoryAvailable,
		CameraHikVisionMQTTTopicStateHDDCapacity,
		CameraHikVisionMQTTTopicStateHDDFree,
		CameraHikVisionMQTTTopicStateHDDUsage,
	}
}

func (d *CameraHikVision) startAlertStreaming() error {
	ctx := context.Background()

	stream, err := d.isapi.EventNotificationAlertStream(ctx)
	if err != nil {
		return err
	}

	go func() {
		sn := d.SerialNumberMQTTEscaped()

		for {
			select {
			case event := <-stream.NextAlert():
				if event.EventState != hikvision.EventEventStateActive {
					continue
				}

				cacheKey := fmt.Sprintf("%d-%s", event.DynChannelID, event.EventType)

				d.mutex.Lock()
				lastFire, ok := d.alertStreamingHistory[cacheKey]
				d.alertStreamingHistory[cacheKey] = event.DateTime
				d.mutex.Unlock()

				if !ok || event.DateTime.Sub(lastFire) > CameraHikVisionIgnoreInterval {
					d.MQTTPublishAsync(ctx, CameraHikVisionMQTTTopicEvent.Format(sn, event.DynChannelID, event.EventType), 0, false, event.EventDescription)
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

func (d *CameraHikVision) updateStatusByChannelId(ctx context.Context, channelId uint64) error {
	channel, ok := d.ptzChannels[channelId]
	if !ok {
		return fmt.Errorf("channel %d not found", channelId)
	}

	status, err := d.isapi.PTZStatus(ctx, channelId)
	if err != nil {
		return err
	}

	sn := d.SerialNumberMQTTEscaped()

	if channel.Status == nil || channel.Status.AbsoluteHigh.Elevation != status.AbsoluteHigh.Elevation {
		d.MQTTPublishAsync(ctx, CameraHikVisionMQTTTopicPTZStatusElevation.Format(sn, channelId), 1, false, status.AbsoluteHigh.Elevation)
	}

	if channel.Status == nil || channel.Status.AbsoluteHigh.Azimuth != status.AbsoluteHigh.Azimuth {
		d.MQTTPublishAsync(ctx, CameraHikVisionMQTTTopicPTZStatusAzimuth.Format(sn, channelId), 1, false, status.AbsoluteHigh.Azimuth)
	}

	if channel.Status == nil || channel.Status.AbsoluteHigh.AbsoluteZoom != status.AbsoluteHigh.AbsoluteZoom {
		d.MQTTPublishAsync(ctx, CameraHikVisionMQTTTopicPTZStatusZoom.Format(sn, channelId), 1, false, status.AbsoluteHigh.AbsoluteZoom)
	}

	channel.Status = &status
	d.ptzChannels[channelId] = channel

	return nil
}

func (d *CameraHikVision) checkTopic(topic string) (uint64, error) {
	if d.ptzChannels == nil || len(d.ptzChannels) == 0 {
		return 0, errors.New("channels is empty")
	}

	parts := mqtt.RouteSplit(topic)

	channelId, err := strconv.ParseUint(parts[4], 10, 64)
	if err != nil {
		return 0, err
	}

	_, ok := d.ptzChannels[channelId]
	if !ok {
		return 0, fmt.Errorf("channel %d not found", channelId)
	}

	return channelId, nil
}

func (d *CameraHikVision) callbackMQTTAbsolute(ctx context.Context, client mqtt.Component, message mqtt.Message) {
	var err error
	defer func() {
		if err != nil {
			log.Printf("Callback for device %s for topic %s with body %s failed with error %s",
				d.SerialNumber(), err.Error(), message.Topic(), string(message.Payload()))
		}
	}()

	channelId, err := d.checkTopic(message.Topic())
	if err != nil {
		return
	}

	var request struct {
		Elevation int64  `json:"elevation,omitempty"`
		Azimuth   uint64 `json:"azimuth,omitempty"`
		Zoom      uint64 `json:"zoom,omitempty"`
	}

	err = json.Unmarshal(message.Payload(), &request)
	if err != nil {
		return
	}

	err = d.isapi.PTZAbsolute(ctx, channelId, request.Elevation, request.Azimuth, request.Zoom)
	if err != nil {
		return
	}

	err = d.updateStatusByChannelId(ctx, channelId)
}

func (d *CameraHikVision) callbackMQTTContinuous(ctx context.Context, client mqtt.Component, message mqtt.Message) {
	var err error
	defer func() {
		if err != nil {
			log.Printf("Callback for device %s for topic %s with body %s failed with error %s",
				d.SerialNumber(), err.Error(), message.Topic(), string(message.Payload()))
		}
	}()

	channelId, err := d.checkTopic(message.Topic())
	if err != nil {
		return
	}

	var request struct {
		Pan  int64 `json:"pan,omitempty"`
		Tilt int64 `json:"tilt,omitempty"`
		Zoom int64 `json:"zoom,omitempty"`
	}

	err = json.Unmarshal(message.Payload(), &request)
	if err != nil {
		return
	}

	err = d.isapi.PTZContinuous(ctx, channelId, request.Pan, request.Tilt, request.Zoom)
	if err != nil {
		return
	}

	err = d.updateStatusByChannelId(ctx, channelId)
}

func (d *CameraHikVision) callbackMQTTRelative(ctx context.Context, client mqtt.Component, message mqtt.Message) {
	var err error
	defer func() {
		if err != nil {
			log.Printf("Callback for device %s for topic %s with body %s failed with error %s",
				d.SerialNumber(), err.Error(), message.Topic(), string(message.Payload()))
		}
	}()

	channelId, err := d.checkTopic(message.Topic())
	if err != nil {
		return
	}

	var request struct {
		X    int64 `xml:"x,omitempty"`
		Y    int64 `xml:"y,omitempty"`
		Zoom int64 `xml:"zoom,omitempty"`
	}

	err = json.Unmarshal(message.Payload(), &request)
	if err != nil {
		return
	}

	err = d.isapi.PTZRelative(ctx, channelId, request.X, request.Y, request.Zoom)
	if err != nil {
		return
	}

	err = d.updateStatusByChannelId(ctx, channelId)
}

func (d *CameraHikVision) callbackMQTTPreset(ctx context.Context, client mqtt.Component, message mqtt.Message) {
	var err error
	defer func() {
		if err != nil {
			log.Printf("Callback for device %s for topic %s with body %s failed with error %s",
				d.SerialNumber(), err.Error(), message.Topic(), string(message.Payload()))
		}
	}()

	channelId, err := d.checkTopic(message.Topic())
	if err != nil {
		return
	}

	presetId, err := strconv.ParseUint(string(message.Payload()), 10, 64)
	if err != nil {
		return
	}

	err = d.isapi.PTZPresetGoTo(ctx, channelId, presetId)
	if err != nil {
		return
	}

	err = d.updateStatusByChannelId(ctx, channelId)
}

func (d *CameraHikVision) callbackMQTTMomentary(ctx context.Context, client mqtt.Component, message mqtt.Message) {
	var err error
	defer func() {
		if err != nil {
			log.Printf("Callback for device %s for topic %s with body %s failed with error %s",
				d.SerialNumber(), err.Error(), message.Topic(), string(message.Payload()))
		}
	}()

	channelId, err := d.checkTopic(message.Topic())
	if err != nil {
		return
	}

	var request struct {
		Pan      int64         `json:"pan,omitempty"`
		Tilt     int64         `json:"tilt,omitempty"`
		Zoom     int64         `json:"zoom,omitempty"`
		Duration time.Duration `json:"duration,omitempty"`
	}

	err = json.Unmarshal(message.Payload(), &request)
	if err != nil {
		return
	}

	duration := time.Duration(request.Duration) * time.Millisecond
	err = d.isapi.PTZMomentary(ctx, channelId, request.Pan, request.Tilt, request.Zoom, duration)
	if err != nil {
		return
	}

	err = d.updateStatusByChannelId(ctx, channelId)
}

func (d *CameraHikVision) callbackMQTTMove(ctx context.Context, client mqtt.Component, message mqtt.Message) {
	var err error
	defer func() {
		if err != nil {
			log.Printf("Callback for device %s for topic %s with body %s failed with error %s",
				d.SerialNumber(), err.Error(), message.Topic(), string(message.Payload()))
		}
	}()

	channelId, err := d.checkTopic(message.Topic())
	if err != nil {
		return
	}

	var pan, tilt, zoom int64

	switch string(message.Payload()) {
	case "right":
		pan = 1
	case "left":
		pan = -1
	case "up":
		tilt = 1
	case "up-right":
		pan = 1
		tilt = 1
	case "up-left":
		pan = -1
		tilt = 1
	case "down":
		tilt = -1
	case "down-right":
		pan = 1
		tilt = -1
	case "down-left":
		pan = -1
		tilt = -1
	case "narrow":
		zoom = 1
	case "wide":
		zoom = -1
	case "stop":
	default:
		err = fmt.Errorf("unknown operation %s", string(message.Payload()))
	}

	err = d.isapi.PTZContinuous(ctx, channelId, pan, tilt, zoom)
	if err != nil {
		return
	}

	err = d.updateStatusByChannelId(ctx, channelId)
}

func (d *CameraHikVision) Snapshot(ctx context.Context, channel uint64, writer io.Writer) error {
	return d.isapi.StreamingPictureToWriter(ctx, channel, writer)
}
