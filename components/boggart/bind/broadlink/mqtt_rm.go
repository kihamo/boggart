package broadlink

import (
	"bytes"
	"context"
	"encoding/json"
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/providers/broadlink"
	"github.com/kihamo/boggart/components/mqtt"
)

const (
	RMMQTTSubscribeTopicCapture       mqtt.Topic = boggart.ComponentName + "/remote-control/+/capture"
	RMMQTTPublishTopicCaptureState    mqtt.Topic = boggart.ComponentName + "/remote-control/+/capture/state"
	RMMQTTSubscribeTopicIR            mqtt.Topic = boggart.ComponentName + "/remote-control/+/ir"
	RMMQTTSubscribeTopicIRCount       mqtt.Topic = boggart.ComponentName + "/remote-control/+/ir/count"
	RMMQTTPublishTopicIRCapture       mqtt.Topic = boggart.ComponentName + "/remote-control/+/capture/ir"
	RMMQTTSubscribeTopicRF315mhz      mqtt.Topic = boggart.ComponentName + "/remote-control/+/rf315mhz"
	RMMQTTSubscribeTopicRF315mhzCount mqtt.Topic = boggart.ComponentName + "/remote-control/+/rf315mhz/count"
	RMMQTTPublishTopicRF315mhzCapture mqtt.Topic = boggart.ComponentName + "/remote-control/+/capture/rf315mhz"
	RMMQTTSubscribeTopicRF433mhz      mqtt.Topic = boggart.ComponentName + "/remote-control/+/rf433mhz"
	RMMQTTSubscribeTopicRF433mhzCount mqtt.Topic = boggart.ComponentName + "/remote-control/+/rf433mhz/count"
	RMMQTTPublishTopicRF433mhzCapture mqtt.Topic = boggart.ComponentName + "/remote-control/+/capture/rf433mhz"
)

func (b *BindRM) SetMQTTClient(client mqtt.Component) {
	b.BindMQTT.SetMQTTClient(client)

	if client != nil {
		client.Publish(context.Background(), RMMQTTPublishTopicCaptureState.Format(mqtt.NameReplace(b.SerialNumber())), 2, true, false)
	}
}

func (b *BindRM) MQTTPublishes() []mqtt.Topic {
	sn := mqtt.NameReplace(b.SerialNumber())

	topics := make([]mqtt.Topic, 0)

	_, supportCapture := b.provider.(BindRMSupportCapture)
	if supportCapture {
		topics = append(topics, mqtt.Topic(RMMQTTPublishTopicCaptureState.Format(sn)))

		if _, ok := b.provider.(BindRMSupportIR); ok {
			topics = append(topics, mqtt.Topic(RMMQTTPublishTopicIRCapture.Format(sn)))
		}

		if _, ok := b.provider.(BindRMSupportRF315Mhz); ok {
			topics = append(topics, mqtt.Topic(RMMQTTPublishTopicRF315mhzCapture.Format(sn)))
		}

		if _, ok := b.provider.(BindRMSupportRF433Mhz); ok {
			topics = append(topics, mqtt.Topic(RMMQTTPublishTopicRF433mhzCapture.Format(sn)))
		}
	}

	return topics
}

func (b *BindRM) MQTTSubscribers() []mqtt.Subscriber {
	sn := mqtt.NameReplace(b.SerialNumber())
	subscribers := make([]mqtt.Subscriber, 0)

	if capture, ok := b.provider.(BindRMSupportCapture); ok {
		// создаем таймер который сразу завершается, чтобы получить проиницилизированный таймер
		captureFlush := make(chan struct{}, 1)
		captureDone := make(chan struct{}, 1)

		captureTimer := time.NewTimer(0)
		<-captureTimer.C

		subscribers = append(subscribers,
			mqtt.NewSubscriber(RMMQTTSubscribeTopicCapture.Format(sn), 0, boggart.WrapMQTTSubscribeDeviceIsOnline(b.Status,
				func(ctx context.Context, client mqtt.Component, message mqtt.Message) error {
					if bytes.Equal(message.Payload(), []byte(`1`)) { // start
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
						err := b.MQTTPublish(ctx, RMMQTTPublishTopicCaptureState.Format(sn), 2, true, true)
						if err != nil {
							return err
						}

						select {
						case <-captureFlush:
						case <-captureTimer.C:

						case <-captureDone:
							return nil
						}

						if err := b.MQTTPublishAsync(ctx, RMMQTTPublishTopicCaptureState.Format(sn), 2, true, false); err != nil {
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
							topicCaptureCode = RMMQTTPublishTopicIRCapture.Format(sn)
						case broadlink.RemoteRF315Mhz:
							topicCaptureCode = RMMQTTPublishTopicRF315mhzCapture.Format(sn)
						case broadlink.RemoteRF433Mhz:
							topicCaptureCode = RMMQTTPublishTopicRF433mhzCapture.Format(sn)
						}

						if topicCaptureCode != "" {
							if err = b.MQTTPublish(ctx, topicCaptureCode, 0, false, code); err != nil {
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
	if ir, ok := b.provider.(BindRMSupportIR); ok {
		subscribers = append(subscribers,
			mqtt.NewSubscriber(RMMQTTSubscribeTopicIR.Format(sn), 0, boggart.WrapMQTTSubscribeDeviceIsOnline(b.Status,
				func(_ context.Context, _ mqtt.Component, message mqtt.Message) error {
					return ir.SendIRRemoteControlCodeAsString(string(message.Payload()), 0)
				})),
			mqtt.NewSubscriber(RMMQTTSubscribeTopicIRCount.Format(sn), 0, boggart.WrapMQTTSubscribeDeviceIsOnline(b.Status,
				func(ctx context.Context, _ mqtt.Component, message mqtt.Message) error {
					var request codeRequest

					if err := json.Unmarshal(message.Payload(), &request); err != nil {
						return err
					}

					return ir.SendIRRemoteControlCodeAsString(request.Code, request.Count)
				})),
		)
	}

	// RF315mhz support
	if rf315mhz, ok := b.provider.(BindRMSupportRF315Mhz); ok {
		subscribers = append(subscribers,
			mqtt.NewSubscriber(RMMQTTSubscribeTopicRF315mhz.Format(sn), 0, boggart.WrapMQTTSubscribeDeviceIsOnline(b.Status,
				func(_ context.Context, _ mqtt.Component, message mqtt.Message) error {
					return rf315mhz.SendRF315MhzRemoteControlCodeAsString(string(message.Payload()), 0)
				})),
			mqtt.NewSubscriber(RMMQTTSubscribeTopicRF315mhzCount.Format(sn), 0, boggart.WrapMQTTSubscribeDeviceIsOnline(b.Status,
				func(ctx context.Context, _ mqtt.Component, message mqtt.Message) error {
					var request codeRequest

					if err := json.Unmarshal(message.Payload(), &request); err != nil {
						return err
					}

					return rf315mhz.SendRF315MhzRemoteControlCodeAsString(request.Code, request.Count)
				})),
		)
	}

	// RF433mhz support
	if rf433mhz, ok := b.provider.(BindRMSupportRF433Mhz); ok {
		subscribers = append(subscribers,
			mqtt.NewSubscriber(RMMQTTSubscribeTopicRF433mhz.Format(sn), 0, boggart.WrapMQTTSubscribeDeviceIsOnline(b.Status,
				func(_ context.Context, _ mqtt.Component, message mqtt.Message) error {
					return rf433mhz.SendRF433MhzRemoteControlCodeAsString(string(message.Payload()), 0)
				})),
			mqtt.NewSubscriber(RMMQTTSubscribeTopicRF433mhzCount.Format(sn), 0, boggart.WrapMQTTSubscribeDeviceIsOnline(b.Status,
				func(ctx context.Context, _ mqtt.Component, message mqtt.Message) error {
					var request codeRequest

					if err := json.Unmarshal(message.Payload(), &request); err != nil {
						return err
					}

					return rf433mhz.SendRF433MhzRemoteControlCodeAsString(request.Code, request.Count)
				})),
		)
	}

	return subscribers
}
