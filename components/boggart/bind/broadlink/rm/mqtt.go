package rm

import (
	"context"
	"encoding/json"
	"time"

	"github.com/kihamo/boggart/components/mqtt"
	"github.com/kihamo/boggart/providers/broadlink"
)

type captureResult struct {
	Type string `json:"type"`
	Code string `json:"code"`
}

func (r captureResult) MarshalBinary() (data []byte, err error) {
	return json.Marshal(r)
}

type codeRequest struct {
	Code  string `json:"code"`
	Count int    `json:"count,omitempty"`
}

func (r *codeRequest) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, r)
}

func (b *Bind) MQTTSubscribers() []mqtt.Subscriber {
	subscribers := make([]mqtt.Subscriber, 0)
	cfg := b.config()
	id := b.Meta().ID()

	if capture, ok := b.provider.(SupportCapture); ok {
		// создаем таймер который сразу завершается, чтобы получить проиницилизированный таймер
		captureFlush := make(chan struct{}, 1)
		captureDone := make(chan struct{}, 1)

		captureTimer := time.NewTimer(0)
		<-captureTimer.C

		subscribers = append(subscribers,
			mqtt.NewSubscriber(cfg.TopicCaptureSwitch.Format(id), 0, b.MQTT().WrapSubscribeDeviceIsOnline(
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
						if err := b.MQTT().Publish(ctx, cfg.TopicCaptureState.Format(id), true); err != nil {
							return err
						}

						select {
						case <-captureFlush:
						case <-captureTimer.C:

						case <-captureDone:
							return nil
						}

						if err := b.MQTT().PublishAsync(ctx, cfg.TopicCaptureState.Format(id), false); err != nil {
							return err
						}

						remoteType, code, err := capture.ReadCapturedRemoteControlCodeAsString()
						if err != nil {
							if err != broadlink.ErrSignalNotCaptured {
								return err
							}

							return nil
						}

						result := captureResult{
							Code: code,
						}

						switch remoteType {
						case broadlink.RemoteIR:
							result.Type = "ir"
						case broadlink.RemoteRF315Mhz:
							result.Type = "rf315"
						case broadlink.RemoteRF433Mhz:
							result.Type = "rf433"
						}

						if err = b.MQTT().Publish(ctx, cfg.TopicCaptureResult.Format(id), result); err != nil {
							return err
						}
					} else if len(captureFlush) == 0 {
						captureFlush <- struct{}{}
					}

					return nil
				})),
		)
	}

	// IR support
	if ir, ok := b.provider.(SupportIR); ok {
		subscribers = append(subscribers,
			mqtt.NewSubscriber(cfg.TopicIR.Format(id), 0, b.MQTT().WrapSubscribeDeviceIsOnline(
				func(ctx context.Context, _ mqtt.Component, message mqtt.Message) error {
					var request codeRequest

					if err := message.Unmarshal(&request); err != nil {
						return err
					}

					return ir.SendIRRemoteControlCodeAsString(request.Code, request.Count)
				})),
		)
	}

	// RF315mhz support
	if rf315mhz, ok := b.provider.(SupportRF315Mhz); ok {
		subscribers = append(subscribers,
			mqtt.NewSubscriber(cfg.TopicRF315.Format(id), 0, b.MQTT().WrapSubscribeDeviceIsOnline(
				func(ctx context.Context, _ mqtt.Component, message mqtt.Message) error {
					var request codeRequest

					if err := message.Unmarshal(&request); err != nil {
						return err
					}

					return rf315mhz.SendRF315MhzRemoteControlCodeAsString(request.Code, request.Count)
				})),
		)
	}

	// RF433mhz support
	if rf433mhz, ok := b.provider.(SupportRF433Mhz); ok {
		subscribers = append(subscribers,
			mqtt.NewSubscriber(cfg.TopicRF433.Format(id), 0, b.MQTT().WrapSubscribeDeviceIsOnline(
				func(ctx context.Context, _ mqtt.Component, message mqtt.Message) error {
					var request codeRequest

					if err := message.Unmarshal(&request); err != nil {
						return err
					}

					return rf433mhz.SendRF433MhzRemoteControlCodeAsString(request.Code, request.Count)
				})),
		)
	}

	return subscribers
}
