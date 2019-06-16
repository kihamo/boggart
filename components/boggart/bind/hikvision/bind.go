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
	"github.com/kihamo/boggart/components/boggart/providers/hikvision"
	"github.com/kihamo/boggart/components/boggart/providers/hikvision2"
	"github.com/kihamo/boggart/components/boggart/providers/hikvision2/client/system"
	"github.com/kihamo/boggart/components/boggart/providers/hikvision2/models"
	"github.com/kihamo/boggart/components/mqtt"
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

	client *hikvision2.Client

	address               url.URL
	alertStreamingHistory map[string]time.Time
	alertStreamingCancel  context.CancelFunc

	ptzChannels map[uint64]PTZChannel

	config *Config
}

func (b *Bind) startAlertStreaming() {
	ctx, cancel := context.WithCancel(context.Background())
	b.alertStreamingCancel = cancel

	stream := b.client.EventNotificationAlertStream(ctx)

	go func() {
		sn := mqtt.NameReplace(b.SerialNumber())

		for {
			select {
			case event := <-stream.NextAlert():
				if event.EventState != models.EventNotificationAlertEventStateActive {
					continue
				}

				cacheKey := strconv.FormatUint(event.DynChannelID, 10) + "-" + event.EventType
				dt := time.Time(event.DateTime)

				b.mutex.Lock()
				lastFire, ok := b.alertStreamingHistory[cacheKey]
				b.alertStreamingHistory[cacheKey] = dt
				b.mutex.Unlock()

				if !ok || dt.Sub(lastFire) > b.config.EventsIgnoreInterval {
					if err := b.MQTTPublishAsync(ctx, MQTTPublishTopicEvent.Format(sn, event.DynChannelID, event.EventType), event.EventDescription); err != nil {
						b.Logger().Error("Send event to MQTT failed", "error", err.Error())
					}
				}

			case err := <-stream.NextError():
				b.Logger().Error("Stream error", "error", err.Error())

			case <-ctx.Done():
				b.alertStreamingCancel = nil
				return
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
		} else if response.Payload.SubCode == hikvision.SubStatusCodeRebootRequired {
			if _, err := b.client.System.Reboot(system.NewRebootParamsWithContext(ctx), nil); err != nil {
				b.Logger().Error("Reboot after firmware update failed", "error", err.Error())
			}
		}
	}()
}

func (b *Bind) Close() error {
	if b.alertStreamingCancel != nil {
		b.alertStreamingCancel()
	}

	return nil
}
