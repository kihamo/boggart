package rm

import (
	"context"
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/providers/broadlink"
	"github.com/kihamo/boggart/components/mqtt"
)

const (
	MQTTPrefix mqtt.Topic = boggart.ComponentName + "/remote-control/+/"

	MQTTSubscribeTopicCapture       = MQTTPrefix + "capture"
	MQTTPublishTopicCaptureState    = MQTTPrefix + "capture/state"
	MQTTSubscribeTopicIR            = MQTTPrefix + "ir"
	MQTTSubscribeTopicIRCount       = MQTTPrefix + "ir/count"
	MQTTPublishTopicIRCapture       = MQTTPrefix + "capture/ir"
	MQTTSubscribeTopicRF315mhz      = MQTTPrefix + "rf315mhz"
	MQTTSubscribeTopicRF315mhzCount = MQTTPrefix + "rf315mhz/count"
	MQTTPublishTopicRF315mhzCapture = MQTTPrefix + "capture/rf315mhz"
	MQTTSubscribeTopicRF433mhz      = MQTTPrefix + "rf433mhz"
	MQTTSubscribeTopicRF433mhzCount = MQTTPrefix + "rf433mhz/count"
	MQTTPublishTopicRF433mhzCapture = MQTTPrefix + "capture/rf433mhz"
)

func (b *Bind) SetMQTTClient(client mqtt.Component) {
	b.BindMQTT.SetMQTTClient(client)

	if client != nil {
		client.Publish(context.Background(), MQTTPublishTopicCaptureState.Format(mqtt.NameReplace(b.SerialNumber())), 2, true, false)
	}
}

func (b *Bind) MQTTPublishes() []mqtt.Topic {
	sn := mqtt.NameReplace(b.SerialNumber())

	topics := make([]mqtt.Topic, 0)

	_, supportCapture := b.provider.(SupportCapture)
	if supportCapture {
		topics = append(topics, mqtt.Topic(MQTTPublishTopicCaptureState.Format(sn)))

		if _, ok := b.provider.(SupportIR); ok {
			topics = append(topics, mqtt.Topic(MQTTPublishTopicIRCapture.Format(sn)))
		}

		if _, ok := b.provider.(SupportRF315Mhz); ok {
			topics = append(topics, mqtt.Topic(MQTTPublishTopicRF315mhzCapture.Format(sn)))
		}

		if _, ok := b.provider.(SupportRF433Mhz); ok {
			topics = append(topics, mqtt.Topic(MQTTPublishTopicRF433mhzCapture.Format(sn)))
		}
	}

	return topics
}

func (b *Bind) MQTTSubscribers() []mqtt.Subscriber {
	sn := mqtt.NameReplace(b.SerialNumber())
	subscribers := make([]mqtt.Subscriber, 0)

	if capture, ok := b.provider.(SupportCapture); ok {
		// создаем таймер который сразу завершается, чтобы получить проиницилизированный таймер
		captureFlush := make(chan struct{}, 1)
		captureDone := make(chan struct{}, 1)

		captureTimer := time.NewTimer(0)
		<-captureTimer.C

		subscribers = append(subscribers,
			mqtt.NewSubscriber(MQTTSubscribeTopicCapture.Format(sn), 0, boggart.WrapMQTTSubscribeDeviceIsOnline(b.Status,
				func(ctx context.Context, _ mqtt.Component, message mqtt.Message) error {
					if message.IsTrue() { // start
						if err := capture.StartCaptureRemoteControlCode(); err != nil {
							return err
						}

						// завершаем предыдущие запуски
						if captureTimer.Reset(b.captureDuration) {
							captureDone <- struct{}{}
						} else {
							for len(captureDone) > 0 {
								<-captureDone
							}
						}

						// вычитываем каналы, если успели забиться
						for len(captureFlush) > 0 {
							<-captureFlush
						}

						for len(captureTimer.C) > 0 {
							<-captureTimer.C
						}

						// стартуем новую запись
						if err := b.MQTTPublish(ctx, MQTTPublishTopicCaptureState.Format(sn), true); err != nil {
							return err
						}

						select {
						case <-captureFlush:
						case <-captureTimer.C:

						case <-captureDone:
							return nil
						}

						if err := b.MQTTPublishAsync(ctx, MQTTPublishTopicCaptureState.Format(sn), false); err != nil {
							return err
						}

						remoteType, code, err := capture.ReadCapturedRemoteControlCodeAsString()
						if err != nil {
							if err != broadlink.ErrSignalNotCaptured {
								return err
							}

							return nil
						}

						var topicCaptureCode string

						switch remoteType {
						case broadlink.RemoteIR:
							topicCaptureCode = MQTTPublishTopicIRCapture.Format(sn)
						case broadlink.RemoteRF315Mhz:
							topicCaptureCode = MQTTPublishTopicRF315mhzCapture.Format(sn)
						case broadlink.RemoteRF433Mhz:
							topicCaptureCode = MQTTPublishTopicRF433mhzCapture.Format(sn)
						}

						if topicCaptureCode != "" {
							if err = b.MQTTPublish(ctx, topicCaptureCode, code); err != nil {
								return err
							}
						}
					} else { // stop
						if len(captureFlush) == 0 {
							captureFlush <- struct{}{}
						}
					}

					return nil
				})),
		)
	}

	type codeRequest struct {
		Code  string `json:"code"`
		Count int    `json:"count"`
	}

	// IR support
	if ir, ok := b.provider.(SupportIR); ok {
		subscribers = append(subscribers,
			mqtt.NewSubscriber(MQTTSubscribeTopicIR.Format(sn), 0, boggart.WrapMQTTSubscribeDeviceIsOnline(b.Status,
				func(_ context.Context, _ mqtt.Component, message mqtt.Message) error {
					return ir.SendIRRemoteControlCodeAsString(message.String(), 0)
				})),
			mqtt.NewSubscriber(MQTTSubscribeTopicIRCount.Format(sn), 0, boggart.WrapMQTTSubscribeDeviceIsOnline(b.Status,
				func(ctx context.Context, _ mqtt.Component, message mqtt.Message) error {
					var request codeRequest

					if err := message.UnmarshalJSON(&request); err != nil {
						return err
					}

					return ir.SendIRRemoteControlCodeAsString(request.Code, request.Count)
				})),
		)
	}

	// RF315mhz support
	if rf315mhz, ok := b.provider.(SupportRF315Mhz); ok {
		subscribers = append(subscribers,
			mqtt.NewSubscriber(MQTTSubscribeTopicRF315mhz.Format(sn), 0, boggart.WrapMQTTSubscribeDeviceIsOnline(b.Status,
				func(_ context.Context, _ mqtt.Component, message mqtt.Message) error {
					return rf315mhz.SendRF315MhzRemoteControlCodeAsString(message.String(), 0)
				})),
			mqtt.NewSubscriber(MQTTSubscribeTopicRF315mhzCount.Format(sn), 0, boggart.WrapMQTTSubscribeDeviceIsOnline(b.Status,
				func(ctx context.Context, _ mqtt.Component, message mqtt.Message) error {
					var request codeRequest

					if err := message.UnmarshalJSON(&request); err != nil {
						return err
					}

					return rf315mhz.SendRF315MhzRemoteControlCodeAsString(request.Code, request.Count)
				})),
		)
	}

	// RF433mhz support
	if rf433mhz, ok := b.provider.(SupportRF433Mhz); ok {
		subscribers = append(subscribers,
			mqtt.NewSubscriber(MQTTSubscribeTopicRF433mhz.Format(sn), 0, boggart.WrapMQTTSubscribeDeviceIsOnline(b.Status,
				func(_ context.Context, _ mqtt.Component, message mqtt.Message) error {
					return rf433mhz.SendRF433MhzRemoteControlCodeAsString(message.String(), 0)
				})),
			mqtt.NewSubscriber(MQTTSubscribeTopicRF433mhzCount.Format(sn), 0, boggart.WrapMQTTSubscribeDeviceIsOnline(b.Status,
				func(ctx context.Context, _ mqtt.Component, message mqtt.Message) error {
					var request codeRequest

					if err := message.UnmarshalJSON(&request); err != nil {
						return err
					}

					return rf433mhz.SendRF433MhzRemoteControlCodeAsString(request.Code, request.Count)
				})),
		)
	}

	return subscribers
}
