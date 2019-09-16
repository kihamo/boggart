package hikvision

import (
	"context"
	"io"
	"net/url"
	"strconv"
	"sync"
	"time"

	"github.com/go-openapi/runtime"
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/mqtt"
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
	boggart.BindBase
	boggart.BindMQTT

	mutex sync.RWMutex

	client *hikvision.Client

	address               url.URL
	alertStreamingHistory map[string]time.Time
	alertStreaming        *hikvision.AlertStreaming

	ptzChannels map[uint64]PTZChannel

	config *Config
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

	if !ok || dt.Sub(lastFire) > b.config.EventsIgnoreInterval {
		sn := mqtt.NameReplace(b.SerialNumber())

		if err := b.MQTTPublishAsync(context.Background(), b.config.TopicEvent.Format(sn, ch, event.EventType), event.ActivePostCount); err != nil {
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
