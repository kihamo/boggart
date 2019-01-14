package hikvision

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
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

func (b *Bind) MQTTPublishes() []mqtt.Topic {
	topics := []mqtt.Topic{
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

	if b.eventsEnabled {
		topics = append(topics, MQTTTopicEvent)
	}

	if b.ptzEnabled {
		topics = append(topics,
			MQTTTopicPTZStatusElevation,
			MQTTTopicPTZStatusAzimuth,
			MQTTTopicPTZStatusZoom,
		)
	}

	return topics
}

func (b *Bind) MQTTSubscribers() []mqtt.Subscriber {
	if !b.ptzEnabled {
		return nil
	}

	return []mqtt.Subscriber{
		mqtt.NewSubscriber(MQTTTopicPTZMove.String(), 0, boggart.WrapMQTTSubscribeDeviceIsOnline(b, b.callbackMQTTAbsolute)),
		mqtt.NewSubscriber(MQTTTopicPTZAbsolute.String(), 0, boggart.WrapMQTTSubscribeDeviceIsOnline(b, b.callbackMQTTAbsolute)),
		mqtt.NewSubscriber(MQTTTopicPTZContinuous.String(), 0, boggart.WrapMQTTSubscribeDeviceIsOnline(b, b.callbackMQTTContinuous)),
		mqtt.NewSubscriber(MQTTTopicPTZRelative.String(), 0, boggart.WrapMQTTSubscribeDeviceIsOnline(b, b.callbackMQTTRelative)),
		mqtt.NewSubscriber(MQTTTopicPTZPreset.String(), 0, boggart.WrapMQTTSubscribeDeviceIsOnline(b, b.callbackMQTTPreset)),
		mqtt.NewSubscriber(MQTTTopicPTZMomentary.String(), 0, boggart.WrapMQTTSubscribeDeviceIsOnline(b, b.callbackMQTTMomentary)),
	}
}

func (b *Bind) updateStatusByChannelId(ctx context.Context, channelId uint64) error {
	channel, ok := b.ptzChannels[channelId]
	if !ok {
		return fmt.Errorf("channel %d not found", channelId)
	}

	status, err := b.isapi.PTZStatus(ctx, channelId)
	if err != nil {
		return err
	}

	sn := mqtt.NameReplace(b.SerialNumber())

	if channel.Status == nil || channel.Status.AbsoluteHigh.Elevation != status.AbsoluteHigh.Elevation {
		// TODO:
		_ = b.MQTTPublishAsync(ctx, MQTTTopicPTZStatusElevation.Format(sn, channelId), 1, false, status.AbsoluteHigh.Elevation)
	}

	if channel.Status == nil || channel.Status.AbsoluteHigh.Azimuth != status.AbsoluteHigh.Azimuth {
		// TODO:
		_ = b.MQTTPublishAsync(ctx, MQTTTopicPTZStatusAzimuth.Format(sn, channelId), 1, false, status.AbsoluteHigh.Azimuth)
	}

	if channel.Status == nil || channel.Status.AbsoluteHigh.AbsoluteZoom != status.AbsoluteHigh.AbsoluteZoom {
		// TODO:
		_ = b.MQTTPublishAsync(ctx, MQTTTopicPTZStatusZoom.Format(sn, channelId), 1, false, status.AbsoluteHigh.AbsoluteZoom)
	}

	channel.Status = &status
	b.ptzChannels[channelId] = channel

	return nil
}

func (b *Bind) checkTopic(topic string) (uint64, error) {
	if b.ptzChannels == nil || len(b.ptzChannels) == 0 {
		return 0, errors.New("channels is empty")
	}

	parts := mqtt.RouteSplit(topic)

	channelId, err := strconv.ParseUint(parts[4], 10, 64)
	if err != nil {
		return 0, err
	}

	_, ok := b.ptzChannels[channelId]
	if !ok {
		return 0, fmt.Errorf("channel %d not found", channelId)
	}

	return channelId, nil
}

func (b *Bind) callbackMQTTAbsolute(ctx context.Context, client mqtt.Component, message mqtt.Message) error {
	if !boggart.CheckSerialNumberInMQTTTopic(b, message.Topic(), 4) {
		return nil
	}

	channelId, err := b.checkTopic(message.Topic())
	if err != nil {
		return err
	}

	var request struct {
		Elevation int64  `json:"elevation,omitempty"`
		Azimuth   uint64 `json:"azimuth,omitempty"`
		Zoom      uint64 `json:"zoom,omitempty"`
	}

	if err := json.Unmarshal(message.Payload(), &request); err != nil {
		return err
	}

	if err := b.isapi.PTZAbsolute(ctx, channelId, request.Elevation, request.Azimuth, request.Zoom); err != nil {
		return err
	}

	return b.updateStatusByChannelId(ctx, channelId)
}

func (b *Bind) callbackMQTTContinuous(ctx context.Context, client mqtt.Component, message mqtt.Message) error {
	if !boggart.CheckSerialNumberInMQTTTopic(b, message.Topic(), 4) {
		return nil
	}

	channelId, err := b.checkTopic(message.Topic())
	if err != nil {
		return err
	}

	var request struct {
		Pan  int64 `json:"pan,omitempty"`
		Tilt int64 `json:"tilt,omitempty"`
		Zoom int64 `json:"zoom,omitempty"`
	}

	if err := json.Unmarshal(message.Payload(), &request); err != nil {
		return err
	}

	if err = b.isapi.PTZContinuous(ctx, channelId, request.Pan, request.Tilt, request.Zoom); err != nil {
		return err
	}

	return b.updateStatusByChannelId(ctx, channelId)
}

func (b *Bind) callbackMQTTRelative(ctx context.Context, client mqtt.Component, message mqtt.Message) error {
	if !boggart.CheckSerialNumberInMQTTTopic(b, message.Topic(), 4) {
		return nil
	}

	channelId, err := b.checkTopic(message.Topic())
	if err != nil {
		return err
	}

	var request struct {
		X    int64 `xml:"x,omitempty"`
		Y    int64 `xml:"y,omitempty"`
		Zoom int64 `xml:"zoom,omitempty"`
	}

	if err := json.Unmarshal(message.Payload(), &request); err != nil {
		return err
	}

	if err := b.isapi.PTZRelative(ctx, channelId, request.X, request.Y, request.Zoom); err != nil {
		return err
	}

	return b.updateStatusByChannelId(ctx, channelId)
}

func (b *Bind) callbackMQTTPreset(ctx context.Context, client mqtt.Component, message mqtt.Message) error {
	if !boggart.CheckSerialNumberInMQTTTopic(b, message.Topic(), 4) {
		return nil
	}

	channelId, err := b.checkTopic(message.Topic())
	if err != nil {
		return err
	}

	presetId, err := strconv.ParseUint(string(message.Payload()), 10, 64)
	if err != nil {
		return err
	}

	if err := b.isapi.PTZPresetGoTo(ctx, channelId, presetId); err != nil {
		return err
	}

	return b.updateStatusByChannelId(ctx, channelId)
}

func (b *Bind) callbackMQTTMomentary(ctx context.Context, client mqtt.Component, message mqtt.Message) error {
	if !boggart.CheckSerialNumberInMQTTTopic(b, message.Topic(), 4) {
		return nil
	}

	channelId, err := b.checkTopic(message.Topic())
	if err != nil {
		return err
	}

	var request struct {
		Pan      int64         `json:"pan,omitempty"`
		Tilt     int64         `json:"tilt,omitempty"`
		Zoom     int64         `json:"zoom,omitempty"`
		Duration time.Duration `json:"duration,omitempty"`
	}

	if err := json.Unmarshal(message.Payload(), &request); err != nil {
		return err
	}

	duration := time.Duration(request.Duration) * time.Millisecond
	if err := b.isapi.PTZMomentary(ctx, channelId, request.Pan, request.Tilt, request.Zoom, duration); err != nil {
		return err
	}

	return b.updateStatusByChannelId(ctx, channelId)
}

/*
func (b *Bind) callbackMQTTMove(ctx context.Context, client mqtt.Component, message mqtt.Message) error {
	channelId, err := b.checkTopic(message.Topic())
	if err != nil {
		return err
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
		return errors.New("unknown operation " + string(message.Payload()))
	}

	if err := b.isapi.PTZContinuous(ctx, channelId, pan, tilt, zoom); err != nil {
		return err
	}

	return b.updateStatusByChannelId(ctx, channelId)
}
*/
