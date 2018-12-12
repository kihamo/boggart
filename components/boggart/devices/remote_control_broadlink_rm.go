package devices

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/providers/broadlink"
	"github.com/kihamo/boggart/components/mqtt"
	"github.com/kihamo/shadow/components/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

const (
	RemoteControlBroadlinkRMCaptureDuration = time.Second * 15
	RemoteControlBroadlinkRMMQTTTopicPrefix = boggart.ComponentName + "/remote-control/"
)

type BroadlinkRMRemoteControl struct {
	boggart.DeviceBase
	boggart.DeviceSerialNumber

	provider *broadlink.RM3Mini
}

func NewBroadlinkRMRemoteControl(provider *broadlink.RM3Mini) *BroadlinkRMRemoteControl {
	device := &BroadlinkRMRemoteControl{
		provider: provider,
	}
	device.Init()
	device.SetSerialNumber(provider.MAC().String())
	device.SetDescription("Socket of Broadlink with IP " + provider.Addr().String() + " and MAC " + provider.MAC().String())

	return device
}

func (d *BroadlinkRMRemoteControl) Types() []boggart.DeviceType {
	return []boggart.DeviceType{
		boggart.DeviceTypeRemoteControll,
	}
}

func (d *BroadlinkRMRemoteControl) Ping(_ context.Context) bool {
	return true
}

func (d *BroadlinkRMRemoteControl) MQTTSubscribers() []mqtt.Subscriber {
	mac := strings.Replace(d.provider.MAC().String(), ":", "-", -1)
	topicCode := RemoteControlBroadlinkRMMQTTTopicPrefix + mac + "/command/"

	type codeRequest struct {
		Code  string `json:"code"`
		Count int    `json:"count"`
	}

	// создаем таймер который сразу завершается, чтобы получить проиницилизированный таймер
	captureFlush := make(chan struct{}, 1)
	captureDone := make(chan struct{}, 1)

	captureTimer := time.NewTimer(0)
	<-captureTimer.C

	return []mqtt.Subscriber{
		mqtt.NewSubscriber(topicCode+"raw/count", 0, d.wrapMQTTSubscriber("command_raw_count",
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

					err = d.provider.SendRemoteControlCodeRawAsString(request.Code, request.Count)
				}

				return err
			})),
		mqtt.NewSubscriber(topicCode+"raw", 0, d.wrapMQTTSubscriber("command_raw",
			func(_ context.Context, _ mqtt.Component, message mqtt.Message) error {
				return d.provider.SendRemoteControlCodeRawAsString(string(message.Payload()), 0)
			})),
		mqtt.NewSubscriber(topicCode+"ir/count", 0, d.wrapMQTTSubscriber("command_ir_count",
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

					err = d.provider.SendIRRemoteControlCodeAsString(request.Code, request.Count)
				}

				return err
			})),
		mqtt.NewSubscriber(topicCode+"ir", 0, d.wrapMQTTSubscriber("command_ir",
			func(_ context.Context, _ mqtt.Component, message mqtt.Message) error {
				return d.provider.SendIRRemoteControlCodeAsString(string(message.Payload()), 0)
			})),
		mqtt.NewSubscriber(topicCode+"capture/start", 0, d.wrapMQTTSubscriber("command_capture_start",
			func(ctx context.Context, client mqtt.Component, message mqtt.Message) error {
				if err := d.provider.StartCaptureRemoteControlCode(); err != nil {
					return err
				}

				// завершаем предыдущие запуски
				if captureTimer.Reset(RemoteControlBroadlinkRMCaptureDuration) {
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
				select {
				case <-captureFlush:
				case <-captureTimer.C:

				case <-captureDone:
					return nil
				}

				remoteType, code, err := d.provider.ReadCapturedRemoteControlCodeAsString()
				if err != nil {
					if err != broadlink.ErrSignalNotCaptured {
						return err
					}

					return nil
				}

				if remoteType == broadlink.RemoteIR {
					if err = client.Publish(ctx, topicCode+"capture/ir", 0, false, code); err != nil {
						return err
					}
				}

				return nil

			})),
		mqtt.NewSubscriber(topicCode+"capture/stop", 0, d.wrapMQTTSubscriber("command_capture_stop",
			func(ctx context.Context, client mqtt.Component, message mqtt.Message) error {
				if len(captureFlush) == 0 {
					captureFlush <- struct{}{}
				}

				return nil
			})),
	}
}

func (d *BroadlinkRMRemoteControl) wrapMQTTSubscriber(operationName string, fn func(context.Context, mqtt.Component, mqtt.Message) error) mqtt.MessageHandler {
	return func(ctx context.Context, client mqtt.Component, message mqtt.Message) {
		if !d.IsEnabled() {
			return
		}

		span, ctx := tracing.StartSpanFromContext(ctx, boggart.DeviceTypeRemoteControll.String(), operationName)
		span.LogFields(
			log.String("mac", d.provider.MAC().String()),
			log.String("ip", d.provider.Addr().String()))
		defer span.Finish()

		if err := fn(ctx, client, message); err != nil {
			tracing.SpanError(span, err)
		}
	}
}
