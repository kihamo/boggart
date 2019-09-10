package hikvision

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/mqtt"
	"github.com/kihamo/boggart/providers/hikvision/client/ptz"
	"github.com/kihamo/boggart/providers/hikvision/models"
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
	MQTTPublishTopicPTZStatusElevation        mqtt.Topic = boggart.ComponentName + "/cctv/+/ptz/+/status/elevation"
	MQTTPublishTopicPTZStatusAzimuth          mqtt.Topic = boggart.ComponentName + "/cctv/+/ptz/+/status/azimuth"
	MQTTPublishTopicPTZStatusZoom             mqtt.Topic = boggart.ComponentName + "/cctv/+/ptz/+/status/zoom"
	MQTTPublishTopicEvent                     mqtt.Topic = boggart.ComponentName + "/cctv/+/event/+/+"
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

	if b.config.EventsEnabled {
		topics = append(topics, MQTTPublishTopicEvent)
	}

	if b.config.PTZEnabled {
		topics = append(topics,
			MQTTPublishTopicPTZStatusElevation,
			MQTTPublishTopicPTZStatusAzimuth,
			MQTTPublishTopicPTZStatusZoom,
		)
	}

	return topics
}

func (b *Bind) MQTTSubscribers() []mqtt.Subscriber {
	subscribers := make([]mqtt.Subscriber, 0)

	if b.config.PTZEnabled {
		subscribers = append(subscribers,
			mqtt.NewSubscriber(MQTTSubscribeTopicPTZMove.String(), 0, boggart.WrapMQTTSubscribeDeviceIsOnline(b.Status, b.callbackMQTTAbsolute)),
			mqtt.NewSubscriber(MQTTSubscribeTopicPTZAbsolute.String(), 0, boggart.WrapMQTTSubscribeDeviceIsOnline(b.Status, b.callbackMQTTAbsolute)),
			mqtt.NewSubscriber(MQTTSubscribeTopicPTZContinuous.String(), 0, boggart.WrapMQTTSubscribeDeviceIsOnline(b.Status, b.callbackMQTTContinuous)),
			mqtt.NewSubscriber(MQTTSubscribeTopicPTZRelative.String(), 0, boggart.WrapMQTTSubscribeDeviceIsOnline(b.Status, b.callbackMQTTRelative)),
			mqtt.NewSubscriber(MQTTSubscribeTopicPTZPreset.String(), 0, boggart.WrapMQTTSubscribeDeviceIsOnline(b.Status, b.callbackMQTTPreset)),
			mqtt.NewSubscriber(MQTTSubscribeTopicPTZMomentary.String(), 0, boggart.WrapMQTTSubscribeDeviceIsOnline(b.Status, b.callbackMQTTMomentary)),
		)
	}

	return subscribers
}

func (b *Bind) updateStatusByChannelId(ctx context.Context, channelId uint64) error {
	b.mutex.RLock()
	channel, ok := b.ptzChannels[channelId]
	b.mutex.RUnlock()

	if !ok {
		return fmt.Errorf("channel %d not found", channelId)
	}

	params := ptz.NewGetPtzStatusParamsWithContext(ctx).
		WithChannel(channelId)
	status, err := b.client.Ptz.GetPtzStatus(params, nil)
	if err != nil {
		return err
	}

	sn := mqtt.NameReplace(b.SerialNumber())
	var result error

	if channel.Status == nil || channel.Status.Elevation != status.Payload.AbsoluteHigh.Elevation {
		if err := b.MQTTPublishAsync(ctx, MQTTPublishTopicPTZStatusElevation.Format(sn, channelId), status.Payload.AbsoluteHigh.Elevation); err != nil {
			result = multierr.Append(result, err)
		}
	}

	if channel.Status == nil || channel.Status.Azimuth != status.Payload.AbsoluteHigh.Azimuth {
		if err := b.MQTTPublishAsync(ctx, MQTTPublishTopicPTZStatusAzimuth.Format(sn, channelId), status.Payload.AbsoluteHigh.Azimuth); err != nil {
			result = multierr.Append(result, err)
		}
	}

	if channel.Status == nil || channel.Status.Zoom != status.Payload.AbsoluteHigh.Zoom {
		if err := b.MQTTPublishAsync(ctx, MQTTPublishTopicPTZStatusZoom.Format(sn, channelId), status.Payload.AbsoluteHigh.Zoom); err != nil {
			result = multierr.Append(result, err)
		}
	}

	channel.Status = status.Payload.AbsoluteHigh

	b.mutex.Lock()
	b.ptzChannels[channelId] = channel
	b.mutex.Unlock()

	if result != nil {
		result = er.Wrap(result, "Failed send to MQTT")
	}

	return result
}

func (b *Bind) checkTopic(topic string) (uint64, error) {
	b.mutex.RLock()
	channels := b.ptzChannels
	b.mutex.RUnlock()

	if channels == nil || len(channels) == 0 {
		return 0, errors.New("channels is empty")
	}

	parts := mqtt.RouteSplit(topic)

	channelId, err := strconv.ParseUint(parts[4], 10, 64)
	if err != nil {
		return 0, err
	}

	_, ok := channels[channelId]
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

	var request models.PtzAbsoluteHigh

	if err := message.UnmarshalJSON(&request); err != nil {
		return err
	}

	params := ptz.NewSetPtzPositionAbsoluteParamsWithContext(ctx).
		WithChannel(channelId).
		WithPTZData(&models.PTZData{
			AbsoluteHigh: &request,
		})

	if _, err := b.client.Ptz.SetPtzPositionAbsolute(params, nil); err != nil {
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

	var request models.PTZData

	if err := message.UnmarshalJSON(&request); err != nil {
		return err
	}

	params := ptz.NewSetPtzContinuousParamsWithContext(ctx).
		WithChannel(channelId).
		WithPTZData(&request)

	if _, err = b.client.Ptz.SetPtzContinuous(params, nil); err != nil {
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

	var request models.PtzRelative

	if err := message.UnmarshalJSON(&request); err != nil {
		return err
	}

	params := ptz.NewSetPtzPositionRelativeParamsWithContext(ctx).
		WithChannel(channelId).
		WithPTZData(&models.PTZData{
			Relative: &request,
		})

	if _, err := b.client.Ptz.SetPtzPositionRelative(params, nil); err != nil {
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

	params := ptz.NewGotoPtzPresetParamsWithContext(ctx).
		WithChannel(channelId).
		WithPreset(presetId)
	if _, err := b.client.Ptz.GotoPtzPreset(params, nil); err != nil {
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

	var request models.PTZData

	if err := message.UnmarshalJSON(&request); err != nil {
		return err
	}

	request.Duration = strfmt.Duration(time.Duration(request.Duration) * time.Millisecond)

	params := ptz.NewSetPtzMomentaryParamsWithContext(ctx).
		WithChannel(channelId).
		WithPTZData(&request)

	if _, err := b.client.Ptz.SetPtzMomentary(params, nil); err != nil {
		return err
	}

	return b.updateStatusByChannelId(ctx, channelId)
}
