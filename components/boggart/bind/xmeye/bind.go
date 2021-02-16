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
	di.ConfigBind
	di.LoggerBind
	di.MetaBind
	di.MetricsBind
	di.MQTTBind
	di.ProbesBind
	di.WidgetBind
	di.WorkersBind

	alarmStreaming *xmeye.AlertStreaming
}

func (b *Bind) config() *Config {
	return b.Config().Bind().(*Config)
}

func (b *Bind) client(ctx context.Context) (*xmeye.Client, error) {
	cfg := b.config()

	password, _ := cfg.Address.User.Password()

	provider, err := xmeye.New(cfg.Address.Host, cfg.Address.User.Username(), password)
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

	cfg := b.config()

	b.alarmStreaming = client.AlarmStreaming(ctx, cfg.AlarmStreamingInterval)
	defer client.Close()

	go func() {
		for {
			select {
			case alarm := <-b.alarmStreaming.NextAlarm():
				// close channel
				if alarm == nil {
					return
				}

				if err := b.MQTT().PublishAsync(ctx, cfg.TopicEvent.Format(b.Meta().SerialNumber(), alarm.Channel, alarm.Event), alarm.Status); err != nil {
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
