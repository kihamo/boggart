package devices

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/providers/broadlink"
	"github.com/kihamo/boggart/components/mqtt"
	"github.com/kihamo/shadow/components/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

const (
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
	topicCode := RemoteControlBroadlinkRMMQTTTopicPrefix + mac + "/code/"

	type codeRequest struct {
		Code  string `json:"code"`
		Count int    `json:"count"`
	}

	return []mqtt.Subscriber{
		mqtt.NewSubscriber(topicCode+"raw/count", 0, d.wrapMQTTSubscriber("code_raw_count", func(message mqtt.Message, span opentracing.Span) error {
			var request codeRequest

			err := json.Unmarshal(message.Payload(), &request)
			if err == nil {
				span.LogFields(
					log.String("code", request.Code),
					log.Int("count", request.Count),
				)

				err = d.provider.SendRemoteControlCodeRawAsString(request.Code, request.Count)
			}

			return err
		})),
		mqtt.NewSubscriber(topicCode+"raw", 0, d.wrapMQTTSubscriber("code_raw", func(message mqtt.Message, span opentracing.Span) error {
			return d.provider.SendRemoteControlCodeRawAsString(string(message.Payload()), 0)
		})),
		mqtt.NewSubscriber(topicCode+"ir/count", 0, d.wrapMQTTSubscriber("code_ir_count", func(message mqtt.Message, span opentracing.Span) error {
			var request codeRequest

			err := json.Unmarshal(message.Payload(), &request)
			if err == nil {
				span.LogFields(
					log.String("code", request.Code),
					log.Int("count", request.Count),
				)

				err = d.provider.SendIRRemoteControlCodeAsString(request.Code, request.Count)
			}

			return err
		})),
		mqtt.NewSubscriber(topicCode+"ir", 0, d.wrapMQTTSubscriber("code_raw", func(message mqtt.Message, span opentracing.Span) error {
			return d.provider.SendIRRemoteControlCodeAsString(string(message.Payload()), 0)
		})),
	}
}

func (d *BroadlinkRMRemoteControl) wrapMQTTSubscriber(operationName string, fn func(mqtt.Message, opentracing.Span) error) mqtt.MessageHandler {
	return func(ctx context.Context, client mqtt.Component, message mqtt.Message) {
		if !d.IsEnabled() {
			return
		}

		span, ctx := tracing.StartSpanFromContext(ctx, boggart.DeviceTypeRemoteControll.String(), operationName)
		span.LogFields(
			log.String("mac", d.provider.MAC().String()),
			log.String("ip", d.provider.Addr().String()))
		defer span.Finish()

		if err := fn(message, span); err != nil {
			tracing.SpanError(span, err)
		}
	}
}
