package xmeye

import (
	"context"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/providers/xmeye"
)

const (
	MB uint64 = 1024 * 1024
)

type Bind struct {
	boggart.BindBase
	boggart.BindMQTT

	config         *Config
	client         *xmeye.Client
	alarmStreaming *xmeye.AlertStreaming
}

func (b *Bind) Close() error {
	return b.client.Close()
}

func (b *Bind) startAlarmStreaming() {
	ctx := context.Background()
	b.alarmStreaming = b.client.AlarmStreaming(ctx, b.config.AlarmStreamingInterval)

	go func() {
		for {
			select {
			case alarm := <-b.alarmStreaming.NextAlarm():
				// close channel
				if alarm == nil {
					return
				}

				if err := b.MQTTPublishAsync(ctx, b.config.TopicEvent.Format(b.SerialNumber(), alarm.Channel, alarm.Event), alarm.Status); err != nil {
					b.Logger().Error("Send alarm to MQTT failed", "error", err.Error())
				}

			case err := <-b.alarmStreaming.NextError():
				// close channel
				if err == nil {
					return
				}

				b.Logger().Error("Stream error", "error", err.Error())
			}
		}
	}()
}
