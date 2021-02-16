package rm

import (
	"context"
	"time"

	"github.com/kihamo/boggart/components/mqtt"
	"github.com/kihamo/boggart/providers/broadlink"
)

func (b *Bind) MQTTSubscribers() []mqtt.Subscriber {
	subscribers := make([]mqtt.Subscriber, 0)
	cfg := b.config()
	mac := cfg.MAC.String()

	if capture, ok := b.provider.(SupportCapture); ok {
		// создаем таймер который сразу завершается, чтобы получить проиницилизированный таймер
		captureFlush := make(chan struct{}, 1)
		captureDone := make(chan struct{}, 1)

		captureTimer := time.NewTimer(0)
		<-captureTimer.C

		subscribers = append(subscribers,
			mqtt.NewSubscriber(cfg.TopicCapture.Format(mac), 0, b.MQTT().WrapSubscribeDeviceIsOnline(
				func(ctx context.Context, _ mqtt.Component, message mqtt.Message) error {
					if message.IsTrue() { // start
						if err := capture.StartCaptureRemoteControlCode(); err != nil {
							return err
						}

						// завершаем предыдущие запуски
						if captureTimer.Reset(cfg.CaptureDuration) {
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
						if err := b.MQTT().Publish(ctx, cfg.TopicCaptureState.Format(mac), true); err != nil {
							return err
						}

						select {
						case <-captureFlush:
						case <-captureTimer.C:

						case <-captureDone:
							return nil
						}

						if err := b.MQTT().PublishAsync(ctx, cfg.TopicCaptureState.Format(mac), false); err != nil {
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
							topicCaptureCode = cfg.TopicIRCapture.Format(mac)
						case broadlink.RemoteRF315Mhz:
							topicCaptureCode = cfg.TopicRF315mhzCapture.Format(mac)
						case broadlink.RemoteRF433Mhz:
							topicCaptureCode = cfg.TopicRF433mhzCapture.Format(mac)
						}

						if topicCaptureCode != "" {
							if err = b.MQTT().Publish(ctx, topicCaptureCode, code); err != nil {
								return err
							}
						}
					} else if len(captureFlush) == 0 {
						captureFlush <- struct{}{}
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
			mqtt.NewSubscriber(cfg.TopicIR.Format(mac), 0, b.MQTT().WrapSubscribeDeviceIsOnline(
				func(_ context.Context, _ mqtt.Component, message mqtt.Message) error {
					return ir.SendIRRemoteControlCodeAsString(message.String(), 0)
				})),
			mqtt.NewSubscriber(cfg.TopicIRCount.Format(mac), 0, b.MQTT().WrapSubscribeDeviceIsOnline(
				func(ctx context.Context, _ mqtt.Component, message mqtt.Message) error {
					var request codeRequest

					if err := message.JSONUnmarshal(&request); err != nil {
						return err
					}

					return ir.SendIRRemoteControlCodeAsString(request.Code, request.Count)
				})),
		)
	}

	// RF315mhz support
	if rf315mhz, ok := b.provider.(SupportRF315Mhz); ok {
		subscribers = append(subscribers,
			mqtt.NewSubscriber(cfg.TopicRF315mhz.Format(mac), 0, b.MQTT().WrapSubscribeDeviceIsOnline(
				func(_ context.Context, _ mqtt.Component, message mqtt.Message) error {
					return rf315mhz.SendRF315MhzRemoteControlCodeAsString(message.String(), 0)
				})),
			mqtt.NewSubscriber(cfg.TopicRF315mhzCount.Format(mac), 0, b.MQTT().WrapSubscribeDeviceIsOnline(
				func(ctx context.Context, _ mqtt.Component, message mqtt.Message) error {
					var request codeRequest

					if err := message.JSONUnmarshal(&request); err != nil {
						return err
					}

					return rf315mhz.SendRF315MhzRemoteControlCodeAsString(request.Code, request.Count)
				})),
		)
	}

	// RF433mhz support
	if rf433mhz, ok := b.provider.(SupportRF433Mhz); ok {
		subscribers = append(subscribers,
			mqtt.NewSubscriber(cfg.TopicRF433mhz.Format(mac), 0, b.MQTT().WrapSubscribeDeviceIsOnline(
				func(_ context.Context, _ mqtt.Component, message mqtt.Message) error {
					return rf433mhz.SendRF433MhzRemoteControlCodeAsString(message.String(), 0)
				})),
			mqtt.NewSubscriber(cfg.TopicRF433mhzCount.Format(mac), 0, b.MQTT().WrapSubscribeDeviceIsOnline(
				func(ctx context.Context, _ mqtt.Component, message mqtt.Message) error {
					var request codeRequest

					if err := message.JSONUnmarshal(&request); err != nil {
						return err
					}

					return rf433mhz.SendRF433MhzRemoteControlCodeAsString(request.Code, request.Count)
				})),
		)
	}

	return subscribers
}
