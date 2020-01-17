package xmeye

import (
	"context"

	"github.com/kihamo/boggart/components/boggart/di"
	"github.com/kihamo/boggart/providers/xmeye"
)

const (
	MB uint64 = 1024 * 1024
)

type Bind struct {
	di.MetaBind
	di.MQTTBind
	di.WorkersBind
	di.LoggerBind
	di.ProbesBind

	config         *Config
	alarmStreaming *xmeye.AlertStreaming
}

func (b *Bind) client(ctx context.Context) (*xmeye.Client, error) {
	password, _ := b.config.Address.User.Password()

	provider, err := xmeye.New(b.config.Address.Host, b.config.Address.User.Username(), password)
	if err != nil {
		return nil, err
	}

	if err := provider.Login(ctx); err != nil {
		return nil, err
	}

	return provider, nil
}

func (b *Bind) startAlarmStreaming() error {
	ctx := context.Background()
	client, err := b.client(ctx)
	if err != nil {
		return err
	}

	b.alarmStreaming = client.AlarmStreaming(ctx, b.config.AlarmStreamingInterval)
	defer client.Close()

	go func() {
		for {
			select {
			case alarm := <-b.alarmStreaming.NextAlarm():
				// close channel
				if alarm == nil {
					return
				}

				if err := b.MQTT().PublishAsync(ctx, b.config.TopicEvent.Format(b.Meta().SerialNumber(), alarm.Channel, alarm.Event), alarm.Status); err != nil {
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

	return nil
}
