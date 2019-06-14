package hikvision

import (
	"context"
	"io"
	"net/url"
	"strconv"
	"sync"
	"time"

	"github.com/kihamo/boggart/components/boggart/providers/hikvision2/client/operations"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/providers/hikvision"
	apiclient "github.com/kihamo/boggart/components/boggart/providers/hikvision2/client"
	"github.com/kihamo/boggart/components/mqtt"
)

const (
	MB uint64 = 1024 * 1024
)

type PTZChannel struct {
	Channel hikvision.PTZChannel
	Status  *hikvision.PTZStatus
}

type Bind struct {
	boggart.BindBase
	boggart.BindMQTT

	mutex sync.RWMutex

	client *apiclient.HikVision

	isapi                 *hikvision.ISAPI
	address               url.URL
	alertStreamingHistory map[string]time.Time
	alertStreamingCancel  context.CancelFunc

	ptzChannels map[uint64]PTZChannel

	config *Config
}

func (b *Bind) startAlertStreaming() error {
	ctx, cancel := context.WithCancel(context.Background())
	b.alertStreamingCancel = cancel

	stream, err := b.isapi.EventNotificationAlertStream(ctx)
	if err != nil {
		return err
	}

	go func() {
		sn := mqtt.NameReplace(b.SerialNumber())

		for {
			select {
			case event := <-stream.NextAlert():
				if event.EventState != hikvision.EventEventStateActive {
					continue
				}

				cacheKey := strconv.FormatUint(event.DynChannelID, 10) + "-" + event.EventType

				b.mutex.Lock()
				lastFire, ok := b.alertStreamingHistory[cacheKey]
				b.alertStreamingHistory[cacheKey] = event.DateTime
				b.mutex.Unlock()

				if !ok || event.DateTime.Sub(lastFire) > b.config.EventsIgnoreInterval {
					if err = b.MQTTPublishAsync(ctx, MQTTPublishTopicEvent.Format(sn, event.DynChannelID, event.EventType), event.EventDescription); err != nil {
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

	return nil
}

func (b *Bind) FirmwareUpdate(firmware io.Reader) {
	go func() {
		ctx := context.Background()

		code, _ := b.isapi.SystemUpdateFirmware(ctx, firmware)
		if code.SubStatusCode == hikvision.SubStatusCodeRebootRequired {
			if _, err := b.client.Operations.PutSystemReboot(operations.NewPutSystemRebootParamsWithContext(ctx), nil); err != nil {
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
