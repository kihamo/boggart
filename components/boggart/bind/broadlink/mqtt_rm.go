package broadlink

import (
	"bytes"
	"context"
	"encoding/json"
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/providers/broadlink"
	"github.com/kihamo/boggart/components/mqtt"
	"github.com/kihamo/shadow/components/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

const (
	RMMQTTTopicCapture         mqtt.Topic = boggart.ComponentName + "/remote-control/+/capture"
	RMMQTTTopicCaptureState    mqtt.Topic = boggart.ComponentName + "/remote-control/+/capture/state"
	RMMQTTTopicIR              mqtt.Topic = boggart.ComponentName + "/remote-control/+/ir"
	RMMQTTTopicIRCount         mqtt.Topic = boggart.ComponentName + "/remote-control/+/ir/count"
	RMMQTTTopicIRCapture       mqtt.Topic = boggart.ComponentName + "/remote-control/+/capture/ir"
	RMMQTTTopicRF315mhz        mqtt.Topic = boggart.ComponentName + "/remote-control/+/rf315mhz"
	RMMQTTTopicRF315mhzCount   mqtt.Topic = boggart.ComponentName + "/remote-control/+/rf315mhz/count"
	RMMQTTTopicRF315mhzCapture mqtt.Topic = boggart.ComponentName + "/remote-control/+/capture/rf315mhz"
	RMMQTTTopicRF433mhz        mqtt.Topic = boggart.ComponentName + "/remote-control/+/rf433mhz"
	RMMQTTTopicRF433mhzCount   mqtt.Topic = boggart.ComponentName + "/remote-control/+/rf433mhz/count"
	RMMQTTTopicRF433mhzCapture mqtt.Topic = boggart.ComponentName + "/remote-control/+/capture/rf433mhz"
)

func (b *BindRM) SetMQTTClient(client mqtt.Component) {
	b.DeviceBindMQTT.SetMQTTClient(client)

	if client != nil {
		client.Publish(context.Background(), RMMQTTTopicCaptureState.Format(mqtt.NameReplace(b.SerialNumber())), 2, true, false)
	}
}

func (b *BindRM) MQTTPublishes() []mqtt.Topic {
	sn := mqtt.NameReplace(b.SerialNumber())

	topics := make([]mqtt.Topic, 0)

	_, supportCapture := b.provider.(BindRMSupportCapture)
	if supportCapture {
		topics = append(topics, mqtt.Topic(RMMQTTTopicCaptureState.Format(sn)))

		if _, ok := b.provider.(BindRMSupportIR); ok {
			topics = append(topics, mqtt.Topic(RMMQTTTopicIRCapture.Format(sn)))
		}

		if _, ok := b.provider.(BindRMSupportRF315Mhz); ok {
			topics = append(topics, mqtt.Topic(RMMQTTTopicRF315mhzCapture.Format(sn)))
		}

		if _, ok := b.provider.(BindRMSupportRF433Mhz); ok {
			topics = append(topics, mqtt.Topic(RMMQTTTopicRF433mhzCapture.Format(sn)))
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
			mqtt.NewSubscriber(RMMQTTTopicCapture.Format(sn), 0, b.wrapMQTTSubscriber("command_capture_start",
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
						err := b.MQTTPublish(ctx, RMMQTTTopicCaptureState.Format(sn), 2, true, true)
						if err != nil {
							return err
						}

						select {
						case <-captureFlush:
						case <-captureTimer.C:

						case <-captureDone:
							return nil
						}

						b.MQTTPublishAsync(ctx, RMMQTTTopicCaptureState.Format(sn), 2, true, false)

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
							topicCaptureCode = RMMQTTTopicIRCapture.Format(sn)
						case broadlink.RemoteRF315Mhz:
							topicCaptureCode = RMMQTTTopicRF315mhzCapture.Format(sn)
						case broadlink.RemoteRF433Mhz:
							topicCaptureCode = RMMQTTTopicRF433mhzCapture.Format(sn)
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
			mqtt.NewSubscriber(RMMQTTTopicIR.Format(sn), 0, b.wrapMQTTSubscriber("command_ir",
				func(_ context.Context, _ mqtt.Component, message mqtt.Message) error {
					return ir.SendIRRemoteControlCodeAsString(string(message.Payload()), 0)
				})),
			mqtt.NewSubscriber(RMMQTTTopicIRCount.Format(sn), 0, b.wrapMQTTSubscriber("command_ir_count",
				func(ctx context.Context, _ mqtt.Component, message mqtt.Message) error {
					var request codeRequest

					err := json.Unmarshal(message.Payload(), &request)
					if err == nil {
						if span := opentracing.SpanFromContext(ctx); span != nil {
							span.LogFields(
								log.String("code", request.Code),
								log.Int("count", request.Count),
							)
						}

						err = ir.SendIRRemoteControlCodeAsString(request.Code, request.Count)
					}

					return err
				})),
		)
	}

	// RF315mhz support
	if rf315mhz, ok := b.provider.(BindRMSupportRF315Mhz); ok {
		subscribers = append(subscribers,
			mqtt.NewSubscriber(RMMQTTTopicRF315mhz.Format(sn), 0, b.wrapMQTTSubscriber("command_rf315mhz",
				func(_ context.Context, _ mqtt.Component, message mqtt.Message) error {
					return rf315mhz.SendRF315MhzRemoteControlCodeAsString(string(message.Payload()), 0)
				})),
			mqtt.NewSubscriber(RMMQTTTopicRF315mhzCount.Format(sn), 0, b.wrapMQTTSubscriber("command_rf315mhz_count",
				func(ctx context.Context, _ mqtt.Component, message mqtt.Message) error {
					var request codeRequest

					err := json.Unmarshal(message.Payload(), &request)
					if err == nil {
						if span := opentracing.SpanFromContext(ctx); span != nil {
							span.LogFields(
								log.String("code", request.Code),
								log.Int("count", request.Count),
							)
						}

						err = rf315mhz.SendRF315MhzRemoteControlCodeAsString(request.Code, request.Count)
					}

					return err
				})),
		)
	}

	// RF433mhz support
	if rf433mhz, ok := b.provider.(BindRMSupportRF433Mhz); ok {
		subscribers = append(subscribers,
			mqtt.NewSubscriber(RMMQTTTopicRF433mhz.Format(sn), 0, b.wrapMQTTSubscriber("command_rf433mhz",
				func(_ context.Context, _ mqtt.Component, message mqtt.Message) error {
					return rf433mhz.SendRF433MhzRemoteControlCodeAsString(string(message.Payload()), 0)
				})),
			mqtt.NewSubscriber(RMMQTTTopicRF433mhzCount.Format(sn), 0, b.wrapMQTTSubscriber("command_rf433mhz_count",
				func(ctx context.Context, _ mqtt.Component, message mqtt.Message) error {
					var request codeRequest

					err := json.Unmarshal(message.Payload(), &request)
					if err == nil {
						if span := opentracing.SpanFromContext(ctx); span != nil {
							span.LogFields(
								log.String("code", request.Code),
								log.Int("count", request.Count),
							)
						}

						err = rf433mhz.SendRF433MhzRemoteControlCodeAsString(request.Code, request.Count)
					}

					return err
				})),
		)
	}

	return subscribers
}

func (b *BindRM) wrapMQTTSubscriber(operationName string, fn func(context.Context, mqtt.Component, mqtt.Message) error) mqtt.MessageHandler {
	return func(ctx context.Context, client mqtt.Component, message mqtt.Message) {
		if b.Status() != boggart.DeviceStatusOnline {
			return
		}

		span, ctx := tracing.StartSpanFromContext(ctx, "remote-control", operationName)
		span.LogFields(
			log.String("mac", b.mac.String()),
			log.String("ip", b.ip.String()))
		defer span.Finish()

		if err := fn(ctx, client, message); err != nil {
			tracing.SpanError(span, err)
		}
	}
}
