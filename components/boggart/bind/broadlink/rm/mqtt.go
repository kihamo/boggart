package rm

import (
	"context"
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/mqtt"
	"github.com/kihamo/boggart/providers/broadlink"
)

func (b *Bind) MQTTPublishes() []mqtt.Topic {
	topics := make([]mqtt.Topic, 0)

	_, supportCapture := b.provider.(SupportCapture)
	if supportCapture {
		topics = append(topics, b.config.TopicCaptureState)

		if _, ok := b.provider.(SupportIR); ok {
			topics = append(topics, b.config.TopicIRCapture)
		}

		if _, ok := b.provider.(SupportRF315Mhz); ok {
			topics = append(topics, b.config.TopicRF315mhzCapture)
		}

		if _, ok := b.provider.(SupportRF433Mhz); ok {
			topics = append(topics, b.config.TopicRF433mhzCapture)
		}
	}

	return topics
}

func (b *Bind) MQTTSubscribers() []mqtt.Subscriber {
	subscribers := make([]mqtt.Subscriber, 0)

	if capture, ok := b.provider.(SupportCapture); ok {
		// создаем таймер который сразу завершается, чтобы получить проиницилизированный таймер
		captureFlush := make(chan struct{}, 1)
		captureDone := make(chan struct{}, 1)

		captureTimer := time.NewTimer(0)
		<-captureTimer.C

		subscribers = append(subscribers,
			mqtt.NewSubscriber(b.config.TopicCapture, 0, boggart.WrapMQTTSubscribeDeviceIsOnline(b.Status,
				func(ctx context.Context, _ mqtt.Component, message mqtt.Message) error {
					if message.IsTrue() { // start
						if err := capture.StartCaptureRemoteControlCode(); err != nil {
							return err
						}

						// завершаем предыдущие запуски
						if captureTimer.Reset(b.config.CaptureDuration) {
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
						if err := b.MQTTPublish(ctx, b.config.TopicCaptureState, true); err != nil {
							return err
						}

						select {
						case <-captureFlush:
						case <-captureTimer.C:

						case <-captureDone:
							return nil
						}

						if err := b.MQTTPublishAsync(ctx, b.config.TopicCaptureState, false); err != nil {
							return err
						}

						remoteType, code, err := capture.ReadCapturedRemoteControlCodeAsString()
						if err != nil {
							if err != broadlink.ErrSignalNotCaptured {
								return err
							}

							return nil
						}

						var topicCaptureCode mqtt.Topic

						switch remoteType {
						case broadlink.RemoteIR:
							topicCaptureCode = b.config.TopicIRCapture
						case broadlink.RemoteRF315Mhz:
							topicCaptureCode = b.config.TopicRF315mhzCapture
						case broadlink.RemoteRF433Mhz:
							topicCaptureCode = b.config.TopicRF433mhzCapture
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
			mqtt.NewSubscriber(b.config.TopicIR, 0, boggart.WrapMQTTSubscribeDeviceIsOnline(b.Status,
				func(_ context.Context, _ mqtt.Component, message mqtt.Message) error {
					return ir.SendIRRemoteControlCodeAsString(message.String(), 0)
				})),
			mqtt.NewSubscriber(b.config.TopicIRCount, 0, boggart.WrapMQTTSubscribeDeviceIsOnline(b.Status,
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
			mqtt.NewSubscriber(b.config.TopicRF315mhz, 0, boggart.WrapMQTTSubscribeDeviceIsOnline(b.Status,
				func(_ context.Context, _ mqtt.Component, message mqtt.Message) error {
					return rf315mhz.SendRF315MhzRemoteControlCodeAsString(message.String(), 0)
				})),
			mqtt.NewSubscriber(b.config.TopicRF315mhzCount, 0, boggart.WrapMQTTSubscribeDeviceIsOnline(b.Status,
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
			mqtt.NewSubscriber(b.config.TopicRF433mhz, 0, boggart.WrapMQTTSubscribeDeviceIsOnline(b.Status,
				func(_ context.Context, _ mqtt.Component, message mqtt.Message) error {
					return rf433mhz.SendRF433MhzRemoteControlCodeAsString(message.String(), 0)
				})),
			mqtt.NewSubscriber(b.config.TopicRF433mhzCount, 0, boggart.WrapMQTTSubscribeDeviceIsOnline(b.Status,
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
