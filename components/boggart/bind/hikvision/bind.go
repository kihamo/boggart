package hikvision

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"sync"
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/providers/hikvision"
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

	mutex    sync.RWMutex
	initOnce sync.Once

	isapi                 *hikvision.ISAPI
	address               url.URL
	alertStreamingHistory map[string]time.Time
	alertStreamingCancel  context.CancelFunc

	ptzChannels map[uint64]PTZChannel

	livenessInterval     time.Duration
	livenessTimeout      time.Duration
	updaterInterval      time.Duration
	updaterTimeout       time.Duration
	ptzInterval          time.Duration
	ptzTimeout           time.Duration
	ptzEnabled           bool
	eventsEnabled        bool
	eventsIgnoreInterval time.Duration
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

				cacheKey := fmt.Sprintf("%d-%s", event.DynChannelID, event.EventType)

				b.mutex.Lock()
				lastFire, ok := b.alertStreamingHistory[cacheKey]
				b.alertStreamingHistory[cacheKey] = event.DateTime
				b.mutex.Unlock()

				if !ok || event.DateTime.Sub(lastFire) > b.eventsIgnoreInterval {
					// TODO: log
					_ = b.MQTTPublishAsync(ctx, MQTTPublishTopicEvent.Format(sn, event.DynChannelID, event.EventType), event.EventDescription)
				}

			case _ = <-stream.NextError():
				// TODO: log errors

			case <-ctx.Done():
				b.alertStreamingCancel = nil
				return
			}
		}
	}()

	return nil
}

func (b *Bind) Snapshot(ctx context.Context, channel uint64, writer io.Writer) error {
	return b.isapi.StreamingPictureToWriter(ctx, channel, writer)
}

func (b *Bind) Close() error {
	if b.alertStreamingCancel != nil {
		b.alertStreamingCancel()
	}

	return nil
}
