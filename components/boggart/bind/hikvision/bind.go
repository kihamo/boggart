package hikvision

import (
	"context"
	"io"
	"strconv"
	"sync"
	"time"

	"github.com/go-openapi/runtime"
	"github.com/kihamo/boggart/components/boggart/di"
	"github.com/kihamo/boggart/protocols/swagger"
	"github.com/kihamo/boggart/providers/hikvision"
	"github.com/kihamo/boggart/providers/hikvision/client/system"
	"github.com/kihamo/boggart/providers/hikvision/models"
)

const (
	MB int64 = 1024 * 1024
)

type PTZChannel struct {
	Channel *models.PtzChannel
	Status  *models.PtzAbsoluteHigh
}

type Bind struct {
	di.ConfigBind
	di.LoggerBind
	di.MetaBind
	di.MetricsBind
	di.MQTTBind
	di.ProbesBind
	di.WidgetBind
	di.WorkersBind

	mutex sync.RWMutex

	client                *hikvision.Client
	alertStreamingHistory map[string]time.Time
	alertStreaming        *hikvision.AlertStreaming
}

func (b *Bind) config() *Config {
	return b.Config().Bind().(*Config)
}

func (b *Bind) Run() error {
	cfg := b.config()

	password, _ := cfg.Address.User.Password()

	b.client = hikvision.New(cfg.Address.Host, cfg.Address.User.Username(), password, cfg.Debug, swagger.NewLogger(
		func(message string) {
			b.Logger().Info(message)
		},
		func(message string) {
			b.Logger().Debug(message)
		}))
	b.alertStreamingHistory = make(map[string]time.Time)

	b.Meta().SetLink(&cfg.Address.URL)

	return nil
}

func (b *Bind) registerEvent(event *models.EventNotificationAlert) {
	if event.EventState != models.EventNotificationAlertEventStateActive {
		return
	}

	ch := event.ChannelID
	if ch == 0 && event.DynChannelID > 0 {
		ch = event.DynChannelID
	}

	cacheKey := strconv.FormatUint(ch, 10) + "-" + event.EventType
	dt := time.Time(event.DateTime)

	b.mutex.Lock()
	lastFire, ok := b.alertStreamingHistory[cacheKey]
	b.alertStreamingHistory[cacheKey] = dt
	b.mutex.Unlock()

	cfg := b.config()

	if !ok || dt.Sub(lastFire) > cfg.EventsIgnoreInterval {
		if err := b.MQTT().PublishAsync(context.Background(), cfg.TopicEvent.Format(b.Meta().SerialNumber(), ch, event.EventType), event.ActivePostCount); err != nil {
			b.Logger().Error("Send event to MQTT failed", "error", err.Error())
		}
	}
}

func (b *Bind) startAlertStreaming() {
	ctx := context.Background()
	b.alertStreaming = b.client.EventNotificationAlertStream(ctx)

	go func() {
		for {
			select {
			case event := <-b.alertStreaming.NextAlert():
				b.registerEvent(event)

			case err := <-b.alertStreaming.NextError():
				b.Logger().Error("Stream error", "error", err.Error())
			}
		}
	}()
}

func (b *Bind) FirmwareUpdate(firmware io.Reader) {
	go func() {
		ctx := context.Background()
		params := system.NewUpdateSystemFirmwareParamsWithContext(ctx).
			WithFile(runtime.NamedReader("digicap.dav", firmware))

		response, err := b.client.System.UpdateSystemFirmware(params, nil)
		if err != nil {
			b.Logger().Error("Firmware update failed", "error", err.Error())
		} else if response.Payload.SubCode == models.StatusSubCodeRebootRequired {
			if _, err := b.client.System.Reboot(system.NewRebootParamsWithContext(ctx), nil); err != nil {
				b.Logger().Error("Reboot after firmware update failed", "error", err.Error())
			}
		}
	}()
}

func (b *Bind) Close() error {
	if b.alertStreaming != nil {
		b.alertStreaming.Close()
	}

	return nil
}
