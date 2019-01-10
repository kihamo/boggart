package bind

import (
	"context"
	"encoding/json"
	"net"
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/providers/broadlink"
	"github.com/kihamo/boggart/components/mqtt"
	"github.com/kihamo/go-workers"
	"github.com/kihamo/go-workers/task"
	"github.com/kihamo/shadow/components/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

const (
	BroadlinkRMCaptureDuration = time.Second * 15

	BroadlinkRMMQTTTopicCommand         mqtt.Topic = boggart.ComponentName + "/remote-control/+/command"
	BroadlinkRMMQTTTopicRawCount        mqtt.Topic = boggart.ComponentName + "/remote-control/+/command/raw/count"
	BroadlinkRMMQTTTopicRaw             mqtt.Topic = boggart.ComponentName + "/remote-control/+/command/raw"
	BroadlinkRMMQTTTopicIRCount         mqtt.Topic = boggart.ComponentName + "/remote-control/+/command/ir/count"
	BroadlinkRMMQTTTopicIR              mqtt.Topic = boggart.ComponentName + "/remote-control/+/command/ir"
	BroadlinkRMMQTTTopicRF315mhz        mqtt.Topic = boggart.ComponentName + "/remote-control/+/command/rf315mhz"
	BroadlinkRMMQTTTopicRF433mhz        mqtt.Topic = boggart.ComponentName + "/remote-control/+/command/rf433mhz"
	BroadlinkRMMQTTTopicCapture         mqtt.Topic = boggart.ComponentName + "/remote-control/+/command/capture"
	BroadlinkRMMQTTTopicCaptureState    mqtt.Topic = boggart.ComponentName + "/remote-control/+/command/capture/state"
	BroadlinkRMMQTTTopicCaptureIR       mqtt.Topic = boggart.ComponentName + "/remote-control/+/command/capture/ir"
	BroadlinkRMMQTTTopicCaptureRF315mhz mqtt.Topic = boggart.ComponentName + "/remote-control/+/command/capture/rf315mhz"
	BroadlinkRMMQTTTopicCaptureRF433mhz mqtt.Topic = boggart.ComponentName + "/remote-control/+/command/capture/rf433mhz"
)

type BroadlinkRM struct {
	boggart.DeviceBindBase
	boggart.DeviceBindMQTT

	provider *broadlink.RMProPlus
}

type BroadlinkRMConfig struct {
	IP  string `valid:"ip,required"`
	MAC string `valid:"mac,required"`
}

func (d BroadlinkRM) Config() interface{} {
	return &BroadlinkRMConfig{}
}

func (d BroadlinkRM) CreateBind(c interface{}) (boggart.DeviceBind, error) {
	config := c.(*BroadlinkRMConfig)

	localAddr, err := broadlink.LocalAddr()
	if err != nil {
		return nil, err
	}

	mac, err := net.ParseMAC(config.MAC)
	if err != nil {
		return nil, err
	}

	ip := net.UDPAddr{
		IP:   net.ParseIP(config.IP),
		Port: broadlink.DevicePort,
	}

	device := &BroadlinkRM{
		provider: broadlink.NewRMProPlus(mac, ip, *localAddr),
	}
	device.Init()
	device.SetSerialNumber(mac.String())

	return device, nil
}

func (d *BroadlinkRM) SetMQTTClient(client mqtt.Component) {
	d.DeviceBindMQTT.SetMQTTClient(client)

	client.Publish(context.Background(), BroadlinkRMMQTTTopicCaptureState.Format(d.SerialNumberMQTTEscaped()), 2, true, "0")
}

func (d *BroadlinkRM) Tasks() []workers.Task {
	taskLiveness := task.NewFunctionTask(d.taskLiveness)
	taskLiveness.SetTimeout(time.Second * 5)
	taskLiveness.SetRepeats(-1)
	taskLiveness.SetRepeatInterval(time.Second * 30)
	taskLiveness.SetName("bind-broadlink-rm-liveness")

	return []workers.Task{
		taskLiveness,
	}
}

func (d *BroadlinkRM) taskLiveness(ctx context.Context) (interface{}, error) {
	// TODO:
	d.UpdateStatus(boggart.DeviceStatusOnline)

	return nil, nil
}

func (d *BroadlinkRM) MQTTTopics() []mqtt.Topic {
	sn := d.SerialNumberMQTTEscaped()

	return []mqtt.Topic{
		mqtt.Topic(BroadlinkRMMQTTTopicCommand.Format(sn)),
		mqtt.Topic(BroadlinkRMMQTTTopicRawCount.Format(sn)),
		mqtt.Topic(BroadlinkRMMQTTTopicRaw.Format(sn)),
		mqtt.Topic(BroadlinkRMMQTTTopicIRCount.Format(sn)),
		mqtt.Topic(BroadlinkRMMQTTTopicIR.Format(sn)),
		mqtt.Topic(BroadlinkRMMQTTTopicRF315mhz.Format(sn)),
		mqtt.Topic(BroadlinkRMMQTTTopicRF433mhz.Format(sn)),
		mqtt.Topic(BroadlinkRMMQTTTopicCapture.Format(sn)),
		mqtt.Topic(BroadlinkRMMQTTTopicCaptureState.Format(sn)),
		mqtt.Topic(BroadlinkRMMQTTTopicCaptureIR.Format(sn)),
		mqtt.Topic(BroadlinkRMMQTTTopicCaptureRF315mhz.Format(sn)),
		mqtt.Topic(BroadlinkRMMQTTTopicCaptureRF433mhz.Format(sn)),
	}
}

func (d *BroadlinkRM) MQTTSubscribers() []mqtt.Subscriber {
	type codeRequest struct {
		Code  string `json:"code"`
		Count int    `json:"count"`
	}

	// создаем таймер который сразу завершается, чтобы получить проиницилизированный таймер
	captureFlush := make(chan struct{}, 1)
	captureDone := make(chan struct{}, 1)

	captureTimer := time.NewTimer(0)
	<-captureTimer.C

	sn := d.SerialNumberMQTTEscaped()

	return []mqtt.Subscriber{
		mqtt.NewSubscriber(BroadlinkRMMQTTTopicRawCount.Format(sn), 0, d.wrapMQTTSubscriber("command_raw_count",
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
		mqtt.NewSubscriber(BroadlinkRMMQTTTopicRaw.Format(sn), 0, d.wrapMQTTSubscriber("command_raw",
			func(_ context.Context, _ mqtt.Component, message mqtt.Message) error {
				return d.provider.SendRemoteControlCodeRawAsString(string(message.Payload()), 0)
			})),
		mqtt.NewSubscriber(BroadlinkRMMQTTTopicIRCount.Format(sn), 0, d.wrapMQTTSubscriber("command_ir_count",
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
		mqtt.NewSubscriber(BroadlinkRMMQTTTopicIR.Format(sn), 0, d.wrapMQTTSubscriber("command_ir",
			func(_ context.Context, _ mqtt.Component, message mqtt.Message) error {
				return d.provider.SendIRRemoteControlCodeAsString(string(message.Payload()), 0)
			})),
		mqtt.NewSubscriber(BroadlinkRMMQTTTopicRF315mhz.Format(sn), 0, d.wrapMQTTSubscriber("command_rf315mhz",
			func(_ context.Context, _ mqtt.Component, message mqtt.Message) error {
				return d.provider.SendRF315MhzRemoteControlCodeAsString(string(message.Payload()), 0)
			})),
		mqtt.NewSubscriber(BroadlinkRMMQTTTopicRF433mhz.Format(sn), 0, d.wrapMQTTSubscriber("command_rf433mhz",
			func(_ context.Context, _ mqtt.Component, message mqtt.Message) error {
				return d.provider.SendRF433MhzRemoteControlCodeAsString(string(message.Payload()), 0)
			})),
		mqtt.NewSubscriber(BroadlinkRMMQTTTopicCapture.Format(sn), 0, d.wrapMQTTSubscriber("command_capture_start",
			func(ctx context.Context, client mqtt.Component, message mqtt.Message) error {
				if string(message.Payload()) != "1" {
					return nil
				}

				if err := d.provider.StartCaptureRemoteControlCode(); err != nil {
					return err
				}

				// завершаем предыдущие запуски
				if captureTimer.Reset(BroadlinkRMCaptureDuration) {
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
				err := d.MQTTPublish(ctx, BroadlinkRMMQTTTopicCaptureState.Format(sn), 2, true, "1")
				if err != nil {
					return err
				}

				select {
				case <-captureFlush:
				case <-captureTimer.C:

				case <-captureDone:
					return nil
				}

				d.MQTTPublishAsync(ctx, BroadlinkRMMQTTTopicCaptureState.Format(sn), 2, true, "0")

				remoteType, code, err := d.provider.ReadCapturedRemoteControlCodeAsString()
				if err != nil {
					if err != broadlink.ErrSignalNotCaptured {
						return err
					}

					return nil
				}

				var topicCaptureCode string

				switch remoteType {
				case broadlink.RemoteIR:
					topicCaptureCode = BroadlinkRMMQTTTopicCaptureIR.Format(sn)
				case broadlink.RemoteRF315Mhz:
					topicCaptureCode = BroadlinkRMMQTTTopicCaptureRF315mhz.Format(sn)
				case broadlink.RemoteRF433Mhz:
					topicCaptureCode = BroadlinkRMMQTTTopicCaptureRF433mhz.Format(sn)
				}

				if topicCaptureCode != "" {
					if err = d.MQTTPublish(ctx, topicCaptureCode, 0, false, code); err != nil {
						return err
					}
				}

				return nil
			})),
		mqtt.NewSubscriber(BroadlinkRMMQTTTopicCapture.Format(sn), 0, d.wrapMQTTSubscriber("command_capture_stop",
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

func (d *BroadlinkRM) wrapMQTTSubscriber(operationName string, fn func(context.Context, mqtt.Component, mqtt.Message) error) mqtt.MessageHandler {
	return func(ctx context.Context, client mqtt.Component, message mqtt.Message) {
		if d.Status() != boggart.DeviceStatusOnline {
			return
		}

		span, ctx := tracing.StartSpanFromContext(ctx, "remote-control", operationName)
		span.LogFields(
			log.String("mac", d.provider.MAC().String()),
			log.String("ip", d.provider.Addr().String()))
		defer span.Finish()

		if err := fn(ctx, client, message); err != nil {
			tracing.SpanError(span, err)
		}
	}
}
