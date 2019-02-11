package hikvision

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/mqtt"
	er "github.com/pkg/errors"
	"go.uber.org/multierr"
)

const (
	MQTTSubscribeTopicPTZMove                 mqtt.Topic = boggart.ComponentName + "/cctv/+/ptz/+/move"
	MQTTSubscribeTopicPTZAbsolute             mqtt.Topic = boggart.ComponentName + "/cctv/+/ptz/+/absolute"
	MQTTSubscribeTopicPTZContinuous           mqtt.Topic = boggart.ComponentName + "/cctv/+/ptz/+/continuous"
	MQTTSubscribeTopicPTZRelative             mqtt.Topic = boggart.ComponentName + "/cctv/+/ptz/+/relative"
	MQTTSubscribeTopicPTZPreset               mqtt.Topic = boggart.ComponentName + "/cctv/+/ptz/+/preset"
	MQTTSubscribeTopicPTZMomentary            mqtt.Topic = boggart.ComponentName + "/cctv/+/ptz/+/momentary"
	MQTTPublishTopicEvent                     mqtt.Topic = boggart.ComponentName + "/cctv/+/+/+"
	MQTTPublishTopicPTZStatusElevation        mqtt.Topic = boggart.ComponentName + "/cctv/+/ptz/+/status/elevation"
	MQTTPublishTopicPTZStatusAzimuth          mqtt.Topic = boggart.ComponentName + "/cctv/+/ptz/+/status/azimuth"
	MQTTPublishTopicPTZStatusZoom             mqtt.Topic = boggart.ComponentName + "/cctv/+/ptz/+/status/zoom"
	MQTTPublishTopicStateModel                mqtt.Topic = boggart.ComponentName + "/cctv/+/state/model"
	MQTTPublishTopicStateFirmwareVersion      mqtt.Topic = boggart.ComponentName + "/cctv/+/state/firmware/version"
	MQTTPublishTopicStateFirmwareReleasedDate mqtt.Topic = boggart.ComponentName + "/cctv/+/state/firmware/release-date"
	MQTTPublishTopicStateUpTime               mqtt.Topic = boggart.ComponentName + "/cctv/+/state/uptime"
	MQTTPublishTopicStateMemoryUsage          mqtt.Topic = boggart.ComponentName + "/cctv/+/state/memory/usage"
	MQTTPublishTopicStateMemoryAvailable      mqtt.Topic = boggart.ComponentName + "/cctv/+/state/memory/available"
	MQTTPublishTopicStateHDDCapacity          mqtt.Topic = boggart.ComponentName + "/cctv/+/state/hdd/+/capacity"
	MQTTPublishTopicStateHDDFree              mqtt.Topic = boggart.ComponentName + "/cctv/+/state/hdd/+/free"
	MQTTPublishTopicStateHDDUsage             mqtt.Topic = boggart.ComponentName + "/cctv/+/state/hdd/+/usage"
)

func (b *Bind) MQTTPublishes() []mqtt.Topic {
	topics := []mqtt.Topic{
		MQTTPublishTopicStateModel,
		MQTTPublishTopicStateFirmwareVersion,
		MQTTPublishTopicStateFirmwareReleasedDate,
		MQTTPublishTopicStateUpTime,
		MQTTPublishTopicStateMemoryUsage,
		MQTTPublishTopicStateMemoryAvailable,
		MQTTPublishTopicStateHDDCapacity,
		MQTTPublishTopicStateHDDFree,
		MQTTPublishTopicStateHDDUsage,
	}

	if b.eventsEnabled {
		topics = append(topics, MQTTPublishTopicEvent)
	}

	if b.ptzEnabled {
		topics = append(topics,
			MQTTPublishTopicPTZStatusElevation,
			MQTTPublishTopicPTZStatusAzimuth,
			MQTTPublishTopicPTZStatusZoom,
		)
	}

	return topics
}

func (b *Bind) MQTTSubscribers() []mqtt.Subscriber {
	if !b.ptzEnabled {
		return nil
	}

	return []mqtt.Subscriber{
		mqtt.NewSubscriber(MQTTSubscribeTopicPTZMove.String(), 0, boggart.WrapMQTTSubscribeDeviceIsOnline(b.Status, b.callbackMQTTAbsolute)),
		mqtt.NewSubscriber(MQTTSubscribeTopicPTZAbsolute.String(), 0, boggart.WrapMQTTSubscribeDeviceIsOnline(b.Status, b.callbackMQTTAbsolute)),
		mqtt.NewSubscriber(MQTTSubscribeTopicPTZContinuous.String(), 0, boggart.WrapMQTTSubscribeDeviceIsOnline(b.Status, b.callbackMQTTContinuous)),
		mqtt.NewSubscriber(MQTTSubscribeTopicPTZRelative.String(), 0, boggart.WrapMQTTSubscribeDeviceIsOnline(b.Status, b.callbackMQTTRelative)),
		mqtt.NewSubscriber(MQTTSubscribeTopicPTZPreset.String(), 0, boggart.WrapMQTTSubscribeDeviceIsOnline(b.Status, b.callbackMQTTPreset)),
		mqtt.NewSubscriber(MQTTSubscribeTopicPTZMomentary.String(), 0, boggart.WrapMQTTSubscribeDeviceIsOnline(b.Status, b.callbackMQTTMomentary)),
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
	var result error

	if channel.Status == nil || channel.Status.AbsoluteHigh.Elevation != status.AbsoluteHigh.Elevation {
		if err := b.MQTTPublishAsync(ctx, MQTTPublishTopicPTZStatusElevation.Format(sn, channelId), status.AbsoluteHigh.Elevation); err != nil {
			result = multierr.Append(result, err)
		}
	}

	if channel.Status == nil || channel.Status.AbsoluteHigh.Azimuth != status.AbsoluteHigh.Azimuth {
		if err := b.MQTTPublishAsync(ctx, MQTTPublishTopicPTZStatusAzimuth.Format(sn, channelId), status.AbsoluteHigh.Azimuth); err != nil {
			result = multierr.Append(result, err)
		}
	}

	if channel.Status == nil || channel.Status.AbsoluteHigh.AbsoluteZoom != status.AbsoluteHigh.AbsoluteZoom {
		if err := b.MQTTPublishAsync(ctx, MQTTPublishTopicPTZStatusZoom.Format(sn, channelId), status.AbsoluteHigh.AbsoluteZoom); err != nil {
			result = multierr.Append(result, err)
		}
	}

	channel.Status = &status
	b.ptzChannels[channelId] = channel

	if result != nil {
		result = er.Wrap(result, "Failed send to MQTT")
	}

	return result
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

	if err := message.UnmarshalJSON(&request); err != nil {
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

	if err := message.UnmarshalJSON(&request); err != nil {
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

	if err := message.UnmarshalJSON(&request); err != nil {
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

	presetId, err := strconv.ParseUint(message.String(), 10, 64)
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

	if err := message.UnmarshalJSON(&request); err != nil {
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

	switch message.String() {
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
		return errors.New("unknown operation " + message.String())
	}

	if err := b.isapi.PTZContinuous(ctx, channelId, pan, tilt, zoom); err != nil {
		return err
	}

	return b.updateStatusByChannelId(ctx, channelId)
}
*/
