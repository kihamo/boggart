package devices

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/url"
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
	MB                      uint64 = 1024 * 1024
	HikVisionIgnoreInterval        = time.Second * 5

	HikVisionMQTTTopicEvent                     mqtt.Topic = boggart.ComponentName + "/cctv/+/+/+"
	HikVisionMQTTTopicPTZMove                   mqtt.Topic = boggart.ComponentName + "/cctv/+/ptz/+/move"
	HikVisionMQTTTopicPTZAbsolute               mqtt.Topic = boggart.ComponentName + "/cctv/+/ptz/+/absolute"
	HikVisionMQTTTopicPTZContinuous             mqtt.Topic = boggart.ComponentName + "/cctv/+/ptz/+/continuous"
	HikVisionMQTTTopicPTZRelative               mqtt.Topic = boggart.ComponentName + "/cctv/+/ptz/+/relative"
	HikVisionMQTTTopicPTZPreset                 mqtt.Topic = boggart.ComponentName + "/cctv/+/ptz/+/preset"
	HikVisionMQTTTopicPTZMomentary              mqtt.Topic = boggart.ComponentName + "/cctv/+/ptz/+/momentary"
	HikVisionMQTTTopicPTZStatusElevation        mqtt.Topic = boggart.ComponentName + "/cctv/+/ptz/+/status/elevation"
	HikVisionMQTTTopicPTZStatusAzimuth          mqtt.Topic = boggart.ComponentName + "/cctv/+/ptz/+/status/azimuth"
	HikVisionMQTTTopicPTZStatusZoom             mqtt.Topic = boggart.ComponentName + "/cctv/+/ptz/+/status/zoom"
	HikVisionMQTTTopicStateModel                mqtt.Topic = boggart.ComponentName + "/cctv/+/state/model"
	HikVisionMQTTTopicStateFirmwareVersion      mqtt.Topic = boggart.ComponentName + "/cctv/+/state/firmware/version"
	HikVisionMQTTTopicStateFirmwareReleasedDate mqtt.Topic = boggart.ComponentName + "/cctv/+/state/firmware/release-date"
	HikVisionMQTTTopicStateUpTime               mqtt.Topic = boggart.ComponentName + "/cctv/+/state/uptime"
	HikVisionMQTTTopicStateMemoryUsage          mqtt.Topic = boggart.ComponentName + "/cctv/+/state/memory/usage"
	HikVisionMQTTTopicStateMemoryAvailable      mqtt.Topic = boggart.ComponentName + "/cctv/+/state/memory/available"
	HikVisionMQTTTopicStateHDDCapacity          mqtt.Topic = boggart.ComponentName + "/cctv/+/state/hdd/+/capacity"
	HikVisionMQTTTopicStateHDDFree              mqtt.Topic = boggart.ComponentName + "/cctv/+/state/hdd/+/free"
	HikVisionMQTTTopicStateHDDUsage             mqtt.Topic = boggart.ComponentName + "/cctv/+/state/hdd/+/usage"
)

type HikVisionPTZChannel struct {
	Channel hikvision.PTZChannel
	Status  *hikvision.PTZStatus
}

type HikVision struct {
	boggart.DeviceBindBase
	boggart.DeviceBindSerialNumber
	boggart.DeviceBindMQTT

	mutex    sync.RWMutex
	initOnce sync.Once

	isapi                 *hikvision.ISAPI
	alertStreamingHistory map[string]time.Time

	ptzChannels map[uint64]HikVisionPTZChannel
}

func (d HikVision) CreateBind(config map[string]interface{}) (boggart.DeviceBind, error) {
	address, ok := config["address"]
	if !ok {
		return nil, errors.New("config option address isn't set")
	}

	if address == "" {
		return nil, errors.New("config option address is empty")
	}

	val := address.(string)

	u, err := url.Parse(val)
	if err != nil {
		return nil, errors.New("config option address has bad value " + val)
	}

	port, _ := strconv.ParseInt(u.Port(), 10, 64)
	password, _ := u.User.Password()

	device := &HikVision{
		isapi:                 hikvision.NewISAPI(u.Hostname(), port, u.User.Username(), password),
		alertStreamingHistory: make(map[string]time.Time),
	}

	device.Init()

	return device, nil
}

func (d *HikVision) Tasks() []workers.Task {
	taskLiveness := task.NewFunctionTask(d.taskLiveness)
	taskLiveness.SetTimeout(time.Second * 5)
	taskLiveness.SetRepeats(-1)
	taskLiveness.SetRepeatInterval(time.Minute)
	taskLiveness.SetName("bind-hikvision-liveness")

	taskPTZStatus := task.NewFunctionTillStopTask(d.taskPTZStatus)
	taskPTZStatus.SetTimeout(time.Second * 5)
	taskPTZStatus.SetRepeats(-1)
	taskPTZStatus.SetRepeatInterval(time.Minute)
	taskPTZStatus.SetName("bind-hikvision-ptz-status")

	taskState := task.NewFunctionTask(d.taskState)
	taskState.SetTimeout(time.Second * 30)
	taskState.SetRepeats(-1)
	taskState.SetRepeatInterval(time.Minute * 15)
	taskState.SetName("bind-hikvision-state")

	return []workers.Task{
		taskLiveness,
		taskPTZStatus,
		taskState,
	}
}

func (d *HikVision) taskLiveness(ctx context.Context) (interface{}, error) {
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
		ptzChannels := make(map[uint64]HikVisionPTZChannel, 0)
		if list, err := d.isapi.PTZChannels(ctx); err == nil {
			for _, channel := range list.Channels {
				ptzChannels[channel.ID] = HikVisionPTZChannel{
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

		d.MQTTPublishAsync(ctx, HikVisionMQTTTopicStateModel.Format(deviceInfo.SerialNumber), 0, true, deviceInfo.Model)
		d.MQTTPublishAsync(ctx, HikVisionMQTTTopicStateFirmwareVersion.Format(deviceInfo.SerialNumber), 0, true, deviceInfo.FirmwareVersion)
		d.MQTTPublishAsync(ctx, HikVisionMQTTTopicStateFirmwareReleasedDate.Format(deviceInfo.SerialNumber), 0, true, deviceInfo.FirmwareReleasedDate)

		d.initOnce.Do(func() {
			sn := d.SerialNumberMQTTEscaped()

			d.MQTTSubscribe(HikVisionMQTTTopicPTZMove.Format(sn), 0, boggart.WrapMQTTSubscribeDeviceIsOnline(d, d.callbackMQTTAbsolute))
			d.MQTTSubscribe(HikVisionMQTTTopicPTZAbsolute.Format(sn), 0, boggart.WrapMQTTSubscribeDeviceIsOnline(d, d.callbackMQTTAbsolute))
			d.MQTTSubscribe(HikVisionMQTTTopicPTZContinuous.Format(sn), 0, boggart.WrapMQTTSubscribeDeviceIsOnline(d, d.callbackMQTTContinuous))
			d.MQTTSubscribe(HikVisionMQTTTopicPTZRelative.Format(sn), 0, boggart.WrapMQTTSubscribeDeviceIsOnline(d, d.callbackMQTTRelative))
			d.MQTTSubscribe(HikVisionMQTTTopicPTZPreset.Format(sn), 0, boggart.WrapMQTTSubscribeDeviceIsOnline(d, d.callbackMQTTPreset))
			d.MQTTSubscribe(HikVisionMQTTTopicPTZMomentary.Format(sn), 0, boggart.WrapMQTTSubscribeDeviceIsOnline(d, d.callbackMQTTMomentary))
		})
	}

	return nil, nil
}

func (d *HikVision) taskPTZStatus(ctx context.Context) (interface{}, error, bool) {
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

func (d *HikVision) taskState(ctx context.Context) (interface{}, error) {
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

	d.MQTTPublishAsync(ctx, HikVisionMQTTTopicStateUpTime.Format(sn), 1, false, status.DeviceUpTime)
	d.MQTTPublishAsync(ctx, HikVisionMQTTTopicStateMemoryAvailable.Format(sn), 1, false, uint64(status.Memory[0].MemoryAvailable.Float64())*MB)
	d.MQTTPublishAsync(ctx, HikVisionMQTTTopicStateMemoryUsage.Format(sn), 1, false, uint64(status.Memory[0].MemoryUsage.Float64())*MB)

	storage, err := d.isapi.ContentManagementStorage(ctx)
	if err != nil {
		return nil, err
	}

	for _, hdd := range storage.HDD {
		d.MQTTPublishAsync(ctx, HikVisionMQTTTopicStateHDDCapacity.Format(sn, hdd.ID), 1, false, hdd.Capacity*MB)
		d.MQTTPublishAsync(ctx, HikVisionMQTTTopicStateHDDFree.Format(sn, hdd.ID), 1, false, hdd.FreeSpace*MB)
		d.MQTTPublishAsync(ctx, HikVisionMQTTTopicStateHDDUsage.Format(sn, hdd.ID), 1, false, (hdd.Capacity-hdd.FreeSpace)*MB)
	}

	return nil, nil
}

func (d *HikVision) MQTTTopics() []mqtt.Topic {
	return []mqtt.Topic{
		HikVisionMQTTTopicEvent,
		HikVisionMQTTTopicPTZMove,
		HikVisionMQTTTopicPTZAbsolute,
		HikVisionMQTTTopicPTZContinuous,
		HikVisionMQTTTopicPTZRelative,
		HikVisionMQTTTopicPTZPreset,
		HikVisionMQTTTopicPTZMomentary,
		HikVisionMQTTTopicPTZStatusElevation,
		HikVisionMQTTTopicPTZStatusAzimuth,
		HikVisionMQTTTopicPTZStatusZoom,
		HikVisionMQTTTopicStateModel,
		HikVisionMQTTTopicStateFirmwareVersion,
		HikVisionMQTTTopicStateFirmwareReleasedDate,
		HikVisionMQTTTopicStateUpTime,
		HikVisionMQTTTopicStateMemoryUsage,
		HikVisionMQTTTopicStateMemoryAvailable,
		HikVisionMQTTTopicStateHDDCapacity,
		HikVisionMQTTTopicStateHDDFree,
		HikVisionMQTTTopicStateHDDUsage,
	}
}

func (d *HikVision) startAlertStreaming() error {
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

				if !ok || event.DateTime.Sub(lastFire) > HikVisionIgnoreInterval {
					d.MQTTPublishAsync(ctx, HikVisionMQTTTopicEvent.Format(sn, event.DynChannelID, event.EventType), 0, false, event.EventDescription)
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

func (d *HikVision) updateStatusByChannelId(ctx context.Context, channelId uint64) error {
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
		d.MQTTPublishAsync(ctx, HikVisionMQTTTopicPTZStatusElevation.Format(sn, channelId), 1, false, status.AbsoluteHigh.Elevation)
	}

	if channel.Status == nil || channel.Status.AbsoluteHigh.Azimuth != status.AbsoluteHigh.Azimuth {
		d.MQTTPublishAsync(ctx, HikVisionMQTTTopicPTZStatusAzimuth.Format(sn, channelId), 1, false, status.AbsoluteHigh.Azimuth)
	}

	if channel.Status == nil || channel.Status.AbsoluteHigh.AbsoluteZoom != status.AbsoluteHigh.AbsoluteZoom {
		d.MQTTPublishAsync(ctx, HikVisionMQTTTopicPTZStatusZoom.Format(sn, channelId), 1, false, status.AbsoluteHigh.AbsoluteZoom)
	}

	channel.Status = &status
	d.ptzChannels[channelId] = channel

	return nil
}

func (d *HikVision) checkTopic(topic string) (uint64, error) {
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

func (d *HikVision) callbackMQTTAbsolute(ctx context.Context, client mqtt.Component, message mqtt.Message) {
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

func (d *HikVision) callbackMQTTContinuous(ctx context.Context, client mqtt.Component, message mqtt.Message) {
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

func (d *HikVision) callbackMQTTRelative(ctx context.Context, client mqtt.Component, message mqtt.Message) {
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

func (d *HikVision) callbackMQTTPreset(ctx context.Context, client mqtt.Component, message mqtt.Message) {
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

func (d *HikVision) callbackMQTTMomentary(ctx context.Context, client mqtt.Component, message mqtt.Message) {
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

func (d *HikVision) callbackMQTTMove(ctx context.Context, client mqtt.Component, message mqtt.Message) {
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

func (d *HikVision) Snapshot(ctx context.Context, channel uint64, writer io.Writer) error {
	return d.isapi.StreamingPictureToWriter(ctx, channel, writer)
}
