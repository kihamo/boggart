package broadlink

import (
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
	RMMQTTTopicCommand         mqtt.Topic = boggart.ComponentName + "/remote-control/+/command"
	RMMQTTTopicRawCount        mqtt.Topic = boggart.ComponentName + "/remote-control/+/command/raw/count"
	RMMQTTTopicRaw             mqtt.Topic = boggart.ComponentName + "/remote-control/+/command/raw"
	RMMQTTTopicIRCount         mqtt.Topic = boggart.ComponentName + "/remote-control/+/command/ir/count"
	RMMQTTTopicIR              mqtt.Topic = boggart.ComponentName + "/remote-control/+/command/ir"
	RMMQTTTopicRF315mhz        mqtt.Topic = boggart.ComponentName + "/remote-control/+/command/rf315mhz"
	RMMQTTTopicRF433mhz        mqtt.Topic = boggart.ComponentName + "/remote-control/+/command/rf433mhz"
	RMMQTTTopicCapture         mqtt.Topic = boggart.ComponentName + "/remote-control/+/command/capture"
	RMMQTTTopicCaptureState    mqtt.Topic = boggart.ComponentName + "/remote-control/+/command/capture/state"
	RMMQTTTopicCaptureIR       mqtt.Topic = boggart.ComponentName + "/remote-control/+/command/capture/ir"
	RMMQTTTopicCaptureRF315mhz mqtt.Topic = boggart.ComponentName + "/remote-control/+/command/capture/rf315mhz"
	RMMQTTTopicCaptureRF433mhz mqtt.Topic = boggart.ComponentName + "/remote-control/+/command/capture/rf433mhz"
)

func (b *BindRM) SetMQTTClient(client mqtt.Component) {
	b.DeviceBindMQTT.SetMQTTClient(client)

	client.Publish(context.Background(), RMMQTTTopicCaptureState.Format(mqtt.NameReplace(b.SerialNumber())), 2, true, "0")
}

func (b *BindRM) MQTTTopics() []mqtt.Topic {
	sn := mqtt.NameReplace(b.SerialNumber())

	return []mqtt.Topic{
		mqtt.Topic(RMMQTTTopicCommand.Format(sn)),
		mqtt.Topic(RMMQTTTopicRawCount.Format(sn)),
		mqtt.Topic(RMMQTTTopicRaw.Format(sn)),
		mqtt.Topic(RMMQTTTopicIRCount.Format(sn)),
		mqtt.Topic(RMMQTTTopicIR.Format(sn)),
		mqtt.Topic(RMMQTTTopicRF315mhz.Format(sn)),
		mqtt.Topic(RMMQTTTopicRF433mhz.Format(sn)),
		mqtt.Topic(RMMQTTTopicCapture.Format(sn)),
		mqtt.Topic(RMMQTTTopicCaptureState.Format(sn)),
		mqtt.Topic(RMMQTTTopicCaptureIR.Format(sn)),
		mqtt.Topic(RMMQTTTopicCaptureRF315mhz.Format(sn)),
		mqtt.Topic(RMMQTTTopicCaptureRF433mhz.Format(sn)),
	}
}

func (b *BindRM) MQTTSubscribers() []mqtt.Subscriber {
	type codeRequest struct {
		Code  string `json:"code"`
		Count int    `json:"count"`
	}

	// создаем таймер который сразу завершается, чтобы получить проиницилизированный таймер
	captureFlush := make(chan struct{}, 1)
	captureDone := make(chan struct{}, 1)

	captureTimer := time.NewTimer(0)
	<-captureTimer.C

	sn := mqtt.NameReplace(b.SerialNumber())

	return []mqtt.Subscriber{
		mqtt.NewSubscriber(RMMQTTTopicRawCount.Format(sn), 0, b.wrapMQTTSubscriber("command_raw_count",
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

					err = b.provider.SendRemoteControlCodeRawAsString(request.Code, request.Count)
				}

				return err
			})),
		mqtt.NewSubscriber(RMMQTTTopicRaw.Format(sn), 0, b.wrapMQTTSubscriber("command_raw",
			func(_ context.Context, _ mqtt.Component, message mqtt.Message) error {
				return b.provider.SendRemoteControlCodeRawAsString(string(message.Payload()), 0)
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

					err = b.provider.SendIRRemoteControlCodeAsString(request.Code, request.Count)
				}

				return err
			})),
		mqtt.NewSubscriber(RMMQTTTopicIR.Format(sn), 0, b.wrapMQTTSubscriber("command_ir",
			func(_ context.Context, _ mqtt.Component, message mqtt.Message) error {
				return b.provider.SendIRRemoteControlCodeAsString(string(message.Payload()), 0)
			})),
		mqtt.NewSubscriber(RMMQTTTopicRF315mhz.Format(sn), 0, b.wrapMQTTSubscriber("command_rf315mhz",
			func(_ context.Context, _ mqtt.Component, message mqtt.Message) error {
				return b.provider.SendRF315MhzRemoteControlCodeAsString(string(message.Payload()), 0)
			})),
		mqtt.NewSubscriber(RMMQTTTopicRF433mhz.Format(sn), 0, b.wrapMQTTSubscriber("command_rf433mhz",
			func(_ context.Context, _ mqtt.Component, message mqtt.Message) error {
				return b.provider.SendRF433MhzRemoteControlCodeAsString(string(message.Payload()), 0)
			})),
		mqtt.NewSubscriber(RMMQTTTopicCapture.Format(sn), 0, b.wrapMQTTSubscriber("command_capture_start",
			func(ctx context.Context, client mqtt.Component, message mqtt.Message) error {
				if string(message.Payload()) != "1" {
					return nil
				}

				if err := b.provider.StartCaptureRemoteControlCode(); err != nil {
					return err
				}

				// завершаем предыдущие запуски
				if captureTimer.Reset(RMCaptureDuration) {
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
				err := b.MQTTPublish(ctx, RMMQTTTopicCaptureState.Format(sn), 2, true, "1")
				if err != nil {
					return err
				}

				select {
				case <-captureFlush:
				case <-captureTimer.C:

				case <-captureDone:
					return nil
				}

				b.MQTTPublishAsync(ctx, RMMQTTTopicCaptureState.Format(sn), 2, true, "0")

				remoteType, code, err := b.provider.ReadCapturedRemoteControlCodeAsString()
				if err != nil {
					if err != broadlink.ErrSignalNotCaptured {
						return err
					}

					return nil
				}

				var topicCaptureCode string

				switch remoteType {
				case broadlink.RemoteIR:
					topicCaptureCode = RMMQTTTopicCaptureIR.Format(sn)
				case broadlink.RemoteRF315Mhz:
					topicCaptureCode = RMMQTTTopicCaptureRF315mhz.Format(sn)
				case broadlink.RemoteRF433Mhz:
					topicCaptureCode = RMMQTTTopicCaptureRF433mhz.Format(sn)
				}

				if topicCaptureCode != "" {
					if err = b.MQTTPublish(ctx, topicCaptureCode, 0, false, code); err != nil {
						return err
					}
				}

				return nil
			})),
		mqtt.NewSubscriber(RMMQTTTopicCapture.Format(sn), 0, b.wrapMQTTSubscriber("command_capture_stop",
			func(ctx context.Context, client mqtt.Component, message mqtt.Message) error {
				if string(message.Payload()) != "0" {
					return nil
				}

				if len(captureFlush) == 0 {
					captureFlush <- struct{}{}
				}

				return nil
			})),
	}
}

func (b *BindRM) wrapMQTTSubscriber(operationName string, fn func(context.Context, mqtt.Component, mqtt.Message) error) mqtt.MessageHandler {
	return func(ctx context.Context, client mqtt.Component, message mqtt.Message) {
		if b.Status() != boggart.DeviceStatusOnline {
			return
		}

		span, ctx := tracing.StartSpanFromContext(ctx, "remote-control", operationName)
		span.LogFields(
			log.String("mac", b.provider.MAC().String()),
			log.String("ip", b.provider.Addr().String()))
		defer span.Finish()

		if err := fn(ctx, client, message); err != nil {
			tracing.SpanError(span, err)
		}
	}
}
