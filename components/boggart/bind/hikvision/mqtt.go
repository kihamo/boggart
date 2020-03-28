package hikvision

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/kihamo/boggart/components/mqtt"
	"github.com/kihamo/boggart/providers/hikvision/client/ptz"
	"github.com/kihamo/boggart/providers/hikvision/models"
	er "github.com/pkg/errors"
	"go.uber.org/multierr"
)

func (b *Bind) MQTTSubscribers() []mqtt.Subscriber {
	subscribers := make([]mqtt.Subscriber, 0)

	if b.config.PTZEnabled {
		subscribers = append(subscribers,
			mqtt.NewSubscriber(b.config.TopicPTZMove, 0, b.MQTT().WrapSubscribeDeviceIsOnline(b.callbackMQTTAbsolute)),
			mqtt.NewSubscriber(b.config.TopicPTZAbsolute, 0, b.MQTT().WrapSubscribeDeviceIsOnline(b.callbackMQTTAbsolute)),
			mqtt.NewSubscriber(b.config.TopicPTZContinuous, 0, b.MQTT().WrapSubscribeDeviceIsOnline(b.callbackMQTTContinuous)),
			mqtt.NewSubscriber(b.config.TopicPTZRelative, 0, b.MQTT().WrapSubscribeDeviceIsOnline(b.callbackMQTTRelative)),
			mqtt.NewSubscriber(b.config.TopicPTZPreset, 0, b.MQTT().WrapSubscribeDeviceIsOnline(b.callbackMQTTPreset)),
			mqtt.NewSubscriber(b.config.TopicPTZMomentary, 0, b.MQTT().WrapSubscribeDeviceIsOnline(b.callbackMQTTMomentary)),
		)
	}

	return subscribers
}

func (b *Bind) updateStatusByChannelID(ctx context.Context, channelID uint64) error {
	b.mutex.RLock()
	channel, ok := b.ptzChannels[channelID]
	b.mutex.RUnlock()

	if !ok {
		return fmt.Errorf("channel %d not found", channelID)
	}

	params := ptz.NewGetPtzStatusParamsWithContext(ctx).
		WithChannel(channelID)
	status, err := b.client.Ptz.GetPtzStatus(params, nil)

	if err != nil {
		return err
	}

	sn := b.Meta().SerialNumber()

	var result error

	if channel.Status == nil || channel.Status.Elevation != status.Payload.AbsoluteHigh.Elevation {
		if err := b.MQTT().PublishAsync(ctx, b.config.TopicPTZStatusElevation.Format(sn, channelID), status.Payload.AbsoluteHigh.Elevation); err != nil {
			result = multierr.Append(result, err)
		}
	}

	if channel.Status == nil || channel.Status.Azimuth != status.Payload.AbsoluteHigh.Azimuth {
		if err := b.MQTT().PublishAsync(ctx, b.config.TopicPTZStatusAzimuth.Format(sn, channelID), status.Payload.AbsoluteHigh.Azimuth); err != nil {
			result = multierr.Append(result, err)
		}
	}

	if channel.Status == nil || channel.Status.Zoom != status.Payload.AbsoluteHigh.Zoom {
		if err := b.MQTT().PublishAsync(ctx, b.config.TopicPTZStatusZoom.Format(sn, channelID), status.Payload.AbsoluteHigh.Zoom); err != nil {
			result = multierr.Append(result, err)
		}
	}

	channel.Status = status.Payload.AbsoluteHigh

	b.mutex.Lock()
	b.ptzChannels[channelID] = channel
	b.mutex.Unlock()

	if result != nil {
		result = er.Wrap(result, "Failed send to MQTT")
	}

	return result
}

func (b *Bind) checkTopic(topic mqtt.Topic) (uint64, error) {
	b.mutex.RLock()
	channels := b.ptzChannels
	b.mutex.RUnlock()

	if len(channels) == 0 {
		return 0, errors.New("channels is empty")
	}

	parts := topic.Split()

	channelID, err := strconv.ParseUint(parts[4], 10, 64)
	if err != nil {
		return 0, err
	}

	_, ok := channels[channelID]
	if !ok {
		return 0, fmt.Errorf("channel %d not found", channelID)
	}

	return channelID, nil
}

func (b *Bind) callbackMQTTAbsolute(ctx context.Context, client mqtt.Component, message mqtt.Message) error {
	if !b.MQTT().CheckSerialNumberInTopic(message.Topic(), 4) {
		return nil
	}

	channelID, err := b.checkTopic(message.Topic())
	if err != nil {
		return err
	}

	var request models.PtzAbsoluteHigh

	if err := message.JSONUnmarshal(&request); err != nil {
		return err
	}

	params := ptz.NewSetPtzPositionAbsoluteParamsWithContext(ctx).
		WithChannel(channelID).
		WithPTZData(&models.PTZData{
			AbsoluteHigh: &request,
		})

	if _, err := b.client.Ptz.SetPtzPositionAbsolute(params, nil); err != nil {
		return err
	}

	return b.updateStatusByChannelID(ctx, channelID)
}

func (b *Bind) callbackMQTTContinuous(ctx context.Context, client mqtt.Component, message mqtt.Message) error {
	if !b.MQTT().CheckSerialNumberInTopic(message.Topic(), 4) {
		return nil
	}

	channelID, err := b.checkTopic(message.Topic())
	if err != nil {
		return err
	}

	var request models.PTZData

	if err := message.JSONUnmarshal(&request); err != nil {
		return err
	}

	params := ptz.NewSetPtzContinuousParamsWithContext(ctx).
		WithChannel(channelID).
		WithPTZData(&request)

	if _, err = b.client.Ptz.SetPtzContinuous(params, nil); err != nil {
		return err
	}

	return b.updateStatusByChannelID(ctx, channelID)
}

func (b *Bind) callbackMQTTRelative(ctx context.Context, client mqtt.Component, message mqtt.Message) error {
	if !b.MQTT().CheckSerialNumberInTopic(message.Topic(), 4) {
		return nil
	}

	channelID, err := b.checkTopic(message.Topic())
	if err != nil {
		return err
	}

	var request models.PtzRelative

	if err := message.JSONUnmarshal(&request); err != nil {
		return err
	}

	params := ptz.NewSetPtzPositionRelativeParamsWithContext(ctx).
		WithChannel(channelID).
		WithPTZData(&models.PTZData{
			Relative: &request,
		})

	if _, err := b.client.Ptz.SetPtzPositionRelative(params, nil); err != nil {
		return err
	}

	return b.updateStatusByChannelID(ctx, channelID)
}

func (b *Bind) callbackMQTTPreset(ctx context.Context, client mqtt.Component, message mqtt.Message) error {
	if !b.MQTT().CheckSerialNumberInTopic(message.Topic(), 4) {
		return nil
	}

	channelID, err := b.checkTopic(message.Topic())
	if err != nil {
		return err
	}

	presetID, err := strconv.ParseUint(message.String(), 10, 64)
	if err != nil {
		return err
	}

	params := ptz.NewGotoPtzPresetParamsWithContext(ctx).
		WithChannel(channelID).
		WithPreset(presetID)
	if _, err := b.client.Ptz.GotoPtzPreset(params, nil); err != nil {
		return err
	}

	return b.updateStatusByChannelID(ctx, channelID)
}

func (b *Bind) callbackMQTTMomentary(ctx context.Context, client mqtt.Component, message mqtt.Message) error {
	if !b.MQTT().CheckSerialNumberInTopic(message.Topic(), 4) {
		return nil
	}

	channelID, err := b.checkTopic(message.Topic())
	if err != nil {
		return err
	}

	var request models.PTZData

	if err := message.JSONUnmarshal(&request); err != nil {
		return err
	}

	request.Duration = strfmt.Duration(time.Duration(request.Duration) * time.Millisecond)

	params := ptz.NewSetPtzMomentaryParamsWithContext(ctx).
		WithChannel(channelID).
		WithPTZData(&request)

	if _, err := b.client.Ptz.SetPtzMomentary(params, nil); err != nil {
		return err
	}

	return b.updateStatusByChannelID(ctx, channelID)
}
