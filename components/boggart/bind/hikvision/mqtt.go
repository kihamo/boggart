package hikvision

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/mqtt"
)

const (
	MQTTTopicEvent                     mqtt.Topic = boggart.ComponentName + "/cctv/+/+/+"
	MQTTTopicPTZMove                   mqtt.Topic = boggart.ComponentName + "/cctv/+/ptz/+/move"
	MQTTTopicPTZAbsolute               mqtt.Topic = boggart.ComponentName + "/cctv/+/ptz/+/absolute"
	MQTTTopicPTZContinuous             mqtt.Topic = boggart.ComponentName + "/cctv/+/ptz/+/continuous"
	MQTTTopicPTZRelative               mqtt.Topic = boggart.ComponentName + "/cctv/+/ptz/+/relative"
	MQTTTopicPTZPreset                 mqtt.Topic = boggart.ComponentName + "/cctv/+/ptz/+/preset"
	MQTTTopicPTZMomentary              mqtt.Topic = boggart.ComponentName + "/cctv/+/ptz/+/momentary"
	MQTTTopicPTZStatusElevation        mqtt.Topic = boggart.ComponentName + "/cctv/+/ptz/+/status/elevation"
	MQTTTopicPTZStatusAzimuth          mqtt.Topic = boggart.ComponentName + "/cctv/+/ptz/+/status/azimuth"
	MQTTTopicPTZStatusZoom             mqtt.Topic = boggart.ComponentName + "/cctv/+/ptz/+/status/zoom"
	MQTTTopicStateModel                mqtt.Topic = boggart.ComponentName + "/cctv/+/state/model"
	MQTTTopicStateFirmwareVersion      mqtt.Topic = boggart.ComponentName + "/cctv/+/state/firmware/version"
	MQTTTopicStateFirmwareReleasedDate mqtt.Topic = boggart.ComponentName + "/cctv/+/state/firmware/release-date"
	MQTTTopicStateUpTime               mqtt.Topic = boggart.ComponentName + "/cctv/+/state/uptime"
	MQTTTopicStateMemoryUsage          mqtt.Topic = boggart.ComponentName + "/cctv/+/state/memory/usage"
	MQTTTopicStateMemoryAvailable      mqtt.Topic = boggart.ComponentName + "/cctv/+/state/memory/available"
	MQTTTopicStateHDDCapacity          mqtt.Topic = boggart.ComponentName + "/cctv/+/state/hdd/+/capacity"
	MQTTTopicStateHDDFree              mqtt.Topic = boggart.ComponentName + "/cctv/+/state/hdd/+/free"
	MQTTTopicStateHDDUsage             mqtt.Topic = boggart.ComponentName + "/cctv/+/state/hdd/+/usage"
)

func (d *Bind) MQTTTopics() []mqtt.Topic {
	return []mqtt.Topic{
		MQTTTopicEvent,
		MQTTTopicPTZMove,
		MQTTTopicPTZAbsolute,
		MQTTTopicPTZContinuous,
		MQTTTopicPTZRelative,
		MQTTTopicPTZPreset,
		MQTTTopicPTZMomentary,
		MQTTTopicPTZStatusElevation,
		MQTTTopicPTZStatusAzimuth,
		MQTTTopicPTZStatusZoom,
		MQTTTopicStateModel,
		MQTTTopicStateFirmwareVersion,
		MQTTTopicStateFirmwareReleasedDate,
		MQTTTopicStateUpTime,
		MQTTTopicStateMemoryUsage,
		MQTTTopicStateMemoryAvailable,
		MQTTTopicStateHDDCapacity,
		MQTTTopicStateHDDFree,
		MQTTTopicStateHDDUsage,
	}
}

func (d *Bind) MQTTSubscribers() []mqtt.Subscriber {
	return []mqtt.Subscriber{
		mqtt.NewSubscriber(MQTTTopicPTZMove.String(), 0, boggart.WrapMQTTSubscribeDeviceIsOnline(d, d.callbackMQTTAbsolute)),
		mqtt.NewSubscriber(MQTTTopicPTZAbsolute.String(), 0, boggart.WrapMQTTSubscribeDeviceIsOnline(d, d.callbackMQTTAbsolute)),
		mqtt.NewSubscriber(MQTTTopicPTZContinuous.String(), 0, boggart.WrapMQTTSubscribeDeviceIsOnline(d, d.callbackMQTTContinuous)),
		mqtt.NewSubscriber(MQTTTopicPTZRelative.String(), 0, boggart.WrapMQTTSubscribeDeviceIsOnline(d, d.callbackMQTTRelative)),
		mqtt.NewSubscriber(MQTTTopicPTZPreset.String(), 0, boggart.WrapMQTTSubscribeDeviceIsOnline(d, d.callbackMQTTPreset)),
		mqtt.NewSubscriber(MQTTTopicPTZMomentary.String(), 0, boggart.WrapMQTTSubscribeDeviceIsOnline(d, d.callbackMQTTMomentary)),
	}
}

func (d *Bind) updateStatusByChannelId(ctx context.Context, channelId uint64) error {
	channel, ok := d.ptzChannels[channelId]
	if !ok {
		return fmt.Errorf("channel %d not found", channelId)
	}

	status, err := d.isapi.PTZStatus(ctx, channelId)
	if err != nil {
		return err
	}

	sn := mqtt.NameReplace(d.SerialNumber())

	if channel.Status == nil || channel.Status.AbsoluteHigh.Elevation != status.AbsoluteHigh.Elevation {
		d.MQTTPublishAsync(ctx, MQTTTopicPTZStatusElevation.Format(sn, channelId), 1, false, status.AbsoluteHigh.Elevation)
	}

	if channel.Status == nil || channel.Status.AbsoluteHigh.Azimuth != status.AbsoluteHigh.Azimuth {
		d.MQTTPublishAsync(ctx, MQTTTopicPTZStatusAzimuth.Format(sn, channelId), 1, false, status.AbsoluteHigh.Azimuth)
	}

	if channel.Status == nil || channel.Status.AbsoluteHigh.AbsoluteZoom != status.AbsoluteHigh.AbsoluteZoom {
		d.MQTTPublishAsync(ctx, MQTTTopicPTZStatusZoom.Format(sn, channelId), 1, false, status.AbsoluteHigh.AbsoluteZoom)
	}

	channel.Status = &status
	d.ptzChannels[channelId] = channel

	return nil
}

func (d *Bind) checkTopic(topic string) (uint64, error) {
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

func (d *Bind) callbackMQTTAbsolute(ctx context.Context, client mqtt.Component, message mqtt.Message) {
	if !boggart.CheckSerialNumberInMQTTTopic(d, message.Topic(), 4) {
		return
	}

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

func (d *Bind) callbackMQTTContinuous(ctx context.Context, client mqtt.Component, message mqtt.Message) {
	if !boggart.CheckSerialNumberInMQTTTopic(d, message.Topic(), 4) {
		return
	}

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

func (d *Bind) callbackMQTTRelative(ctx context.Context, client mqtt.Component, message mqtt.Message) {
	if !boggart.CheckSerialNumberInMQTTTopic(d, message.Topic(), 4) {
		return
	}

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

func (d *Bind) callbackMQTTPreset(ctx context.Context, client mqtt.Component, message mqtt.Message) {
	if !boggart.CheckSerialNumberInMQTTTopic(d, message.Topic(), 4) {
		return
	}

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

func (d *Bind) callbackMQTTMomentary(ctx context.Context, client mqtt.Component, message mqtt.Message) {
	if !boggart.CheckSerialNumberInMQTTTopic(d, message.Topic(), 4) {
		return
	}

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

func (d *Bind) callbackMQTTMove(ctx context.Context, client mqtt.Component, message mqtt.Message) {
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
