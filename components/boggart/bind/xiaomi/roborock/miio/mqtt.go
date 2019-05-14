package miio

import (
	"context"
	"strconv"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/mqtt"
)

const (
	MQTTPrefix mqtt.Topic = boggart.ComponentName + "/xiaomi/roborock/+/"

	MQTTPublishTopicBattery             = MQTTPrefix + "battery"
	MQTTPublishTopicCleanArea           = MQTTPrefix + "clean/area"
	MQTTPublishTopicCleanTime           = MQTTPrefix + "clean/time"
	MQTTPublishTopicFanPower            = MQTTPrefix + "fan-power"
	MQTTPublishTopicVolume              = MQTTPrefix + "volume"
	MQTTPublishTopicConsumableFilter    = MQTTPrefix + "consumable/filter"
	MQTTPublishTopicConsumableBrushMain = MQTTPrefix + "consumable/brush-main"
	MQTTPublishTopicConsumableBrushSide = MQTTPrefix + "consumable/brush-side"
	MQTTPublishTopicConsumableSensor    = MQTTPrefix + "consumable/sensor"
	MQTTPublishTopicState               = MQTTPrefix + "state"
	MQTTPublishTopicError               = MQTTPrefix + "error"

	MQTTSubscribeTopicSetFanPower = MQTTPrefix + "fan-power/set"
	MQTTSubscribeTopicSetVolume   = MQTTPrefix + "volume/set"
	MQTTSubscribeTopicTestVolume  = MQTTPrefix + "volume/test"
	MQTTSubscribeTopicFind        = MQTTPrefix + "find"
	MQTTSubscribeTopicAction      = MQTTPrefix + "action"
)

func (b *Bind) MQTTPublishes() []mqtt.Topic {
	return []mqtt.Topic{
		MQTTPublishTopicBattery,
		MQTTPublishTopicCleanArea,
		MQTTPublishTopicCleanTime,
		MQTTPublishTopicFanPower,
		MQTTPublishTopicConsumableFilter,
		MQTTPublishTopicConsumableBrushMain,
		MQTTPublishTopicConsumableBrushSide,
		MQTTPublishTopicConsumableSensor,
		MQTTPublishTopicState,
		MQTTPublishTopicError,
	}
}

func (b *Bind) MQTTSubscribers() []mqtt.Subscriber {
	return []mqtt.Subscriber{
		mqtt.NewSubscriber(MQTTSubscribeTopicSetFanPower.String(), 0, boggart.WrapMQTTSubscribeDeviceIsOnline(b.Status, b.callbackMQTTSetFanPower)),
		mqtt.NewSubscriber(MQTTSubscribeTopicSetVolume.String(), 0, boggart.WrapMQTTSubscribeDeviceIsOnline(b.Status, b.callbackMQTTSetVolume)),
		mqtt.NewSubscriber(MQTTSubscribeTopicTestVolume.String(), 0, boggart.WrapMQTTSubscribeDeviceIsOnline(b.Status, b.callbackMQTTTestVolume)),
		mqtt.NewSubscriber(MQTTSubscribeTopicFind.String(), 0, boggart.WrapMQTTSubscribeDeviceIsOnline(b.Status, b.callbackMQTTFind)),
		mqtt.NewSubscriber(MQTTSubscribeTopicAction.String(), 0, boggart.WrapMQTTSubscribeDeviceIsOnline(b.Status, b.callbackMQTTAction)),
	}
}

func (b *Bind) callbackMQTTSetFanPower(ctx context.Context, client mqtt.Component, message mqtt.Message) error {
	if !boggart.CheckSerialNumberInMQTTTopic(b, message.Topic(), 3) {
		return nil
	}

	power, err := strconv.ParseUint(message.String(), 10, 64)
	if err != nil {
		return err
	}

	err = b.device.SetFanPower(ctx, power)
	if err == nil {
		err = b.updateFanPower(ctx)
	}

	return err
}

func (b *Bind) callbackMQTTSetVolume(ctx context.Context, client mqtt.Component, message mqtt.Message) error {
	if !boggart.CheckSerialNumberInMQTTTopic(b, message.Topic(), 3) {
		return nil
	}

	power, err := strconv.ParseUint(message.String(), 10, 64)
	if err != nil {
		return err
	}

	err = b.device.SetSoundVolume(ctx, uint32(power))
	if err == nil {
		err = b.updateVolume(ctx)
	}

	return err
}

func (b *Bind) callbackMQTTTestVolume(ctx context.Context, client mqtt.Component, message mqtt.Message) error {
	if !boggart.CheckSerialNumberInMQTTTopic(b, message.Topic(), 3) {
		return nil
	}

	return b.device.SoundVolumeTest(ctx)
}

func (b *Bind) callbackMQTTFind(ctx context.Context, client mqtt.Component, message mqtt.Message) error {
	if !boggart.CheckSerialNumberInMQTTTopic(b, message.Topic(), 2) {
		return nil
	}

	return b.device.Find(ctx)
}

func (b *Bind) callbackMQTTAction(ctx context.Context, client mqtt.Component, message mqtt.Message) error {
	if !boggart.CheckSerialNumberInMQTTTopic(b, message.Topic(), 2) {
		return nil
	}

	var err error

	switch message.String() {
	case "start":
		err = b.device.Start(ctx)
	case "spot":
		err = b.device.Spot(ctx)
	case "stop":
		err = b.device.Stop(ctx)
	case "pause":
		err = b.device.Pause(ctx)
	case "home":
		err = b.device.Home(ctx)
	}

	if err == nil {
		err = b.updateStatus(ctx)
	}

	return err
}

func (b *Bind) updateStatus(ctx context.Context) error {
	sn := b.SerialNumber()
	if sn == "" {
		return nil
	}

	status, err := b.device.Status(ctx)
	if err == nil {
		snMQTT := mqtt.NameReplace(sn)

		state := uint32(status.State)
		if ok := b.state.Set(state); ok {
			err = b.MQTTPublishAsync(ctx, MQTTPublishTopicState.Format(snMQTT), state)
		}

		e := uint32(status.Error)
		if ok := b.state.Set(e); ok {
			err = b.MQTTPublishAsync(ctx, MQTTPublishTopicError.Format(snMQTT), e)
		}
	}

	return err
}

func (b *Bind) updateFanPower(ctx context.Context) error {
	sn := b.SerialNumber()
	if sn == "" {
		return nil
	}

	fan, err := b.device.FanPower(ctx)
	if err == nil {
		if ok := b.fanPower.Set(fan); ok {
			err = b.MQTTPublishAsync(ctx, MQTTPublishTopicFanPower.Format(mqtt.NameReplace(sn)), fan)
		}
	}

	return err
}

func (b *Bind) updateVolume(ctx context.Context) error {
	sn := b.SerialNumber()
	if sn == "" {
		return nil
	}

	volume, err := b.device.SoundVolume(ctx)
	if err == nil {
		if ok := b.volume.Set(volume); ok {
			err = b.MQTTPublishAsync(ctx, MQTTPublishTopicVolume.Format(mqtt.NameReplace(sn)), volume)
		}
	}

	return err
}
