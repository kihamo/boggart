package chromecast

import (
	"bytes"
	"context"
	"errors"
	"strconv"
	"strings"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/mqtt"
)

const (
	MQTTPrefix = boggart.ComponentName + "/chromecast/+/"

	MQTTSubscribeTopicVolume     mqtt.Topic = MQTTPrefix + "volume"
	MQTTSubscribeTopicMute       mqtt.Topic = MQTTPrefix + "mute"
	MQTTSubscribeTopicPause      mqtt.Topic = MQTTPrefix + "pause"
	MQTTSubscribeTopicStop       mqtt.Topic = MQTTPrefix + "stop"
	MQTTSubscribeTopicPlay       mqtt.Topic = MQTTPrefix + "play"
	MQTTSubscribeTopicResume     mqtt.Topic = MQTTPrefix + "resume"
	MQTTSubscribeTopicSeek       mqtt.Topic = MQTTPrefix + "seek"
	MQTTSubscribeTopicAction     mqtt.Topic = MQTTPrefix + "action"
	MQTTPublishTopicStateStatus  mqtt.Topic = MQTTPrefix + "state/status"
	MQTTPublishTopicStateVolume  mqtt.Topic = MQTTPrefix + "state/volume"
	MQTTPublishTopicStateMute    mqtt.Topic = MQTTPrefix + "state/mute"
	MQTTPublishTopicStateContent mqtt.Topic = MQTTPrefix + "state/content"
)

func (b *Bind) MQTTPublishes() []mqtt.Topic {
	return []mqtt.Topic{
		MQTTPublishTopicStateStatus,
		MQTTPublishTopicStateVolume,
		MQTTPublishTopicStateMute,
		MQTTPublishTopicStateContent,
	}
}

func (b *Bind) MQTTSubscribers() []mqtt.Subscriber {
	return []mqtt.Subscriber{
		mqtt.NewSubscriber(MQTTSubscribeTopicVolume.String(), 0, boggart.WrapMQTTSubscribeDeviceIsOnline(b, func(ctx context.Context, _ mqtt.Component, message mqtt.Message) error {
			if !boggart.CheckSerialNumberInMQTTTopic(b, message.Topic(), 2) {
				return nil
			}

			volume, err := strconv.ParseInt(string(message.Payload()), 10, 64)
			if err != nil {
				return err
			}

			return b.SetVolume(ctx, volume)
		})),
		mqtt.NewSubscriber(MQTTSubscribeTopicMute.String(), 0, boggart.WrapMQTTSubscribeDeviceIsOnline(b, func(ctx context.Context, _ mqtt.Component, message mqtt.Message) error {
			if !boggart.CheckSerialNumberInMQTTTopic(b, message.Topic(), 2) {
				return nil
			}

			return b.SetMute(ctx, bytes.Equal(message.Payload(), []byte(`1`)))
		})),
		mqtt.NewSubscriber(MQTTSubscribeTopicPause.String(), 0, boggart.WrapMQTTSubscribeDeviceIsOnline(b, func(ctx context.Context, _ mqtt.Component, message mqtt.Message) error {
			if !boggart.CheckSerialNumberInMQTTTopic(b, message.Topic(), 2) {
				return nil
			}

			return b.Pause(ctx)
		})),
		mqtt.NewSubscriber(MQTTSubscribeTopicStop.String(), 0, boggart.WrapMQTTSubscribeDeviceIsOnline(b, func(ctx context.Context, _ mqtt.Component, message mqtt.Message) error {
			if !boggart.CheckSerialNumberInMQTTTopic(b, message.Topic(), 2) {
				return nil
			}

			return b.Stop(ctx)
		})),
		mqtt.NewSubscriber(MQTTSubscribeTopicPlay.String(), 0, boggart.WrapMQTTSubscribeDeviceIsOnline(b, func(ctx context.Context, _ mqtt.Component, message mqtt.Message) error {
			if !boggart.CheckSerialNumberInMQTTTopic(b, message.Topic(), 2) {
				return nil
			}

			if u := string(message.Payload()); u != "" {
				return b.PlayFromURL(ctx, u)
			}

			return b.Resume(ctx)
		})),
		mqtt.NewSubscriber(MQTTSubscribeTopicSeek.String(), 0, boggart.WrapMQTTSubscribeDeviceIsOnline(b, func(ctx context.Context, _ mqtt.Component, message mqtt.Message) error {
			if !boggart.CheckSerialNumberInMQTTTopic(b, message.Topic(), 2) {
				return nil
			}

			second, err := strconv.ParseUint(string(message.Payload()), 10, 64)
			if err != nil {
				return err
			}

			return b.Seek(ctx, second)
		})),
		mqtt.NewSubscriber(MQTTSubscribeTopicResume.String(), 0, boggart.WrapMQTTSubscribeDeviceIsOnline(b, func(ctx context.Context, _ mqtt.Component, message mqtt.Message) error {
			if !boggart.CheckSerialNumberInMQTTTopic(b, message.Topic(), 2) {
				return nil
			}

			return b.Resume(ctx)
		})),
		mqtt.NewSubscriber(MQTTSubscribeTopicAction.String(), 0, boggart.WrapMQTTSubscribeDeviceIsOnline(b, func(ctx context.Context, _ mqtt.Component, message mqtt.Message) error {
			if !boggart.CheckSerialNumberInMQTTTopic(b, message.Topic(), 2) {
				return nil
			}

			action := string(message.Payload())

			switch strings.ToLower(action) {
			case "stop":
				return b.Stop(ctx)

			case "pause":
				return b.Pause(ctx)

			case "play", "resume":
				return b.Resume(ctx)
			}

			return errors.New("unknown action " + action)
		})),
	}
}
