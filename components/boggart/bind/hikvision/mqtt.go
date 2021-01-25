package hikvision

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/kihamo/boggart/components/mqtt"
	"github.com/kihamo/boggart/providers/hikvision/client/ptz"
	"github.com/kihamo/boggart/providers/hikvision/models"
	"go.uber.org/multierr"
)

func (b *Bind) updateStatusByChannelID(ctx context.Context, channelID uint64) error {
	params := ptz.NewGetPtzStatusParamsWithContext(ctx).
		WithChannel(channelID)
	status, err := b.client.Ptz.GetPtzStatus(params, nil)

	if err != nil {
		return err
	}

	sn := b.Meta().SerialNumber()

	if e := b.MQTT().PublishAsync(ctx, b.config.TopicPTZStatusElevation.Format(sn, channelID), status.Payload.AbsoluteHigh.Elevation); e != nil {
		err = multierr.Append(err, e)
	}

	if e := b.MQTT().PublishAsync(ctx, b.config.TopicPTZStatusAzimuth.Format(sn, channelID), status.Payload.AbsoluteHigh.Azimuth); e != nil {
		err = multierr.Append(err, e)
	}

	if e := b.MQTT().PublishAsync(ctx, b.config.TopicPTZStatusZoom.Format(sn, channelID), status.Payload.AbsoluteHigh.Zoom); e != nil {
		err = multierr.Append(err, e)
	}

	return err
}

func (b *Bind) checkTopic(topic mqtt.Topic) (uint64, error) {
	parts := topic.Split()
	if len(parts) < 5 {
		return 0, errors.New("bad topic")
	}

	return strconv.ParseUint(parts[4], 10, 64)
}

func (b *Bind) callbackMQTTAbsolute(ctx context.Context, client mqtt.Component, message mqtt.Message) error {
	if !b.MQTT().CheckSerialNumberInTopic(message.Topic(), -4) {
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
	if !b.MQTT().CheckSerialNumberInTopic(message.Topic(), -4) {
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
	if !b.MQTT().CheckSerialNumberInTopic(message.Topic(), -4) {
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
	if !b.MQTT().CheckSerialNumberInTopic(message.Topic(), -4) {
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
	if !b.MQTT().CheckSerialNumberInTopic(message.Topic(), -4) {
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
