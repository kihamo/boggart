package hikvision

import (
	"context"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/providers/hikvision"
	"github.com/kihamo/boggart/components/mqtt"
)

const (
	MB             uint64 = 1024 * 1024
	IgnoreInterval        = time.Second * 5
)

type HikVisionPTZChannel struct {
	Channel hikvision.PTZChannel
	Status  *hikvision.PTZStatus
}

type Bind struct {
	boggart.DeviceBindBase
	boggart.DeviceBindMQTT

	mutex    sync.RWMutex
	initOnce sync.Once

	isapi                 *hikvision.ISAPI
	alertStreamingHistory map[string]time.Time

	ptzChannels map[uint64]HikVisionPTZChannel
}

func (d *Bind) startAlertStreaming() error {
	ctx := context.Background()

	stream, err := d.isapi.EventNotificationAlertStream(ctx)
	if err != nil {
		return err
	}

	go func() {
		sn := mqtt.NameReplace(d.SerialNumber())

		for {
			select {
			case event := <-stream.NextAlert():
				if event.EventState != hikvision.EventEventStateActive {
					continue
				}

				cacheKey := fmt.Sprintf("%d-%s", event.DynChannelID, event.EventType)

				d.mutex.Lock()
				lastFire, ok := d.alertStreamingHistory[cacheKey]
				d.alertStreamingHistory[cacheKey] = event.DateTime
				d.mutex.Unlock()

				if !ok || event.DateTime.Sub(lastFire) > IgnoreInterval {
					d.MQTTPublishAsync(ctx, MQTTTopicEvent.Format(sn, event.DynChannelID, event.EventType), 0, false, event.EventDescription)
				}

			case _ = <-stream.NextError():
				// TODO: log errors

			case <-ctx.Done():
				return
			}
		}
	}()

	return nil
}

func (d *Bind) Snapshot(ctx context.Context, channel uint64, writer io.Writer) error {
	return d.isapi.StreamingPictureToWriter(ctx, channel, writer)
}
