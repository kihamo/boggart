package chromecast

import (
	"context"
	"errors"
	"math"
	"net"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"go.uber.org/multierr"

	"github.com/barnybug/go-cast"
	"github.com/barnybug/go-cast/controllers"
	"github.com/barnybug/go-cast/events"
	castnet "github.com/barnybug/go-cast/net"
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/mqtt"
)

type Bind struct {
	volume      int64
	mute        int64
	status      atomic.Value
	mediaStatus atomic.Value

	boggart.BindBase
	boggart.BindMQTT

	host net.IP
	port int

	events chan events.Event
	mutex  sync.RWMutex

	conn       *castnet.Connection
	connection *controllers.ConnectionController
	heartbeat  *controllers.HeartbeatController
	receiver   *controllers.ReceiverController
	media      *controllers.MediaController

	livenessInterval time.Duration
	livenessTimeout  time.Duration
}

func (b *Bind) initCast() error {
	b.events = make(chan events.Event, 16)
	b.conn = castnet.NewConnection()

	b.connection = controllers.NewConnectionController(b.conn, b.events, cast.DefaultSender, cast.DefaultReceiver)
	b.heartbeat = controllers.NewHeartbeatController(b.conn, b.events, cast.TransportSender, cast.TransportReceiver)
	b.receiver = controllers.NewReceiverController(b.conn, b.events, cast.DefaultSender, cast.DefaultReceiver)

	go b.doEvents()

	return nil
}

func (b *Bind) Connect(_ context.Context) error {
	status := b.Status()
	if status != boggart.BindStatusInitializing && status != boggart.BindStatusOffline {
		return nil
	}

	ctx := context.Background()

	// open TCP connection
	err := b.conn.Connect(ctx, b.host, b.port)
	if err != nil {
		return err
	}

	if err := b.heartbeat.Start(context.Background()); err != nil {
		return err
	}

	// start main connection controller
	if err := b.connection.Start(ctx); err != nil {
		return err
	}

	// start receiver
	if err := b.receiver.Start(ctx); err != nil {
		return err
	}

	b.events <- events.Connected{}
	return nil
}

func (b *Bind) Close() (err error) {
	b.UpdateStatus(boggart.BindStatusOffline)

	// close running aps
	b.receiver.QuitApp(context.Background())

	// close main connection controller
	if e := b.connection.Close(); e != nil {
		err = multierr.Append(err, e)
	}

	// close TCP connection
	if e := b.conn.Close(); e != nil {
		err = multierr.Append(err, e)
	}

	b.mutex.Lock()
	if b.media != nil {
		b.media.Close()
		b.media = nil
	}
	b.mutex.Unlock()

	return err
}

func (b *Bind) doEvents() {
	for {
		select {
		case event := <-b.events:
			switch t := event.(type) {
			case events.Connected:
				//fmt.Println("events.Connected", t)
				b.UpdateStatus(boggart.BindStatusOnline)

			case events.Disconnected: // from ReceiverController
				//fmt.Println("events.Disconnected", t)

				if err := b.Close(); err != nil {
					// TODO: log
				}

			case events.AppStarted: // from ReceiverController
				//fmt.Println("events.AppStarted", t)

			case events.AppStopped: // from ReceiverController
				//fmt.Println("events.AppStopped", t)

				if t.AppID == cast.AppMedia {
					b.mutex.Lock()

					if b.media != nil {
						if err := b.media.Close(); err != nil {
							// TODO: log
						}

						b.media = nil
					}

					b.mutex.Unlock()
				}

			case events.StatusUpdated: // from ReceiverController
				// fmt.Println("events.StatusUpdated", t)

				ctx := context.Background()
				sn := mqtt.NameReplace(b.SerialNumber())

				currentVolume := int64(math.Round(t.Level * 100))
				prevVolume := atomic.LoadInt64(&b.volume)
				if currentVolume != prevVolume {
					atomic.StoreInt64(&b.volume, currentVolume)

					_ = b.MQTTPublishAsync(ctx, MQTTPublishTopicStateVolume.Format(sn), 0, true, currentVolume)
				}

				prevMute := atomic.LoadInt64(&b.mute) == 1
				if t.Muted != prevMute {
					if t.Muted {
						atomic.StoreInt64(&b.mute, 1)
					} else {
						atomic.StoreInt64(&b.mute, 0)
					}

					_ = b.MQTTPublishAsync(ctx, MQTTPublishTopicStateMute.Format(sn), 0, true, t.Muted)
				}

			case controllers.MediaStatus:
				ctx := context.Background()
				sn := mqtt.NameReplace(b.SerialNumber())

				prev := b.status.Load()
				if prev == nil || prev.(string) != t.PlayerState {
					b.status.Store(t.PlayerState)

					_ = b.MQTTPublishAsync(ctx, MQTTPublishTopicStateStatus.Format(sn), 0, true, strings.ToLower(t.PlayerState))
				}

				if t.PlayerState == "IDLE" && t.IdleReason == "FINISHED" {
					// TODO: error log
					_, _ = b.receiver.QuitApp(ctx)
				}

				if t.Media != nil {
					prev := b.mediaStatus.Load()
					if prev != nil {
						prevMedia := prev.(*controllers.MediaStatusMedia)
						if t.Media.ContentId != prevMedia.ContentId {
							b.mediaStatus.Store(t.Media)

							_ = b.MQTTPublishAsync(ctx, MQTTPublishTopicStateContent.Format(sn), 0, true, t.Media.ContentId)
						}
					} else {
						b.mediaStatus.Store(t.Media)

						_ = b.MQTTPublishAsync(ctx, MQTTPublishTopicStateContent.Format(sn), 0, true, t.Media.ContentId)
					}
				}

			default:
				// TODO: log
				// fmt.Printf("Unknown event: %#v\n", t)
			}
		}
	}
}

func (b *Bind) Media(ctx context.Context) (*controllers.MediaController, error) {
	b.mutex.RLock()
	controller := b.media
	b.mutex.RUnlock()

	if controller == nil {
		transportId, err := b.launchApp(ctx, cast.AppMedia)
		if err != nil {
			return nil, err
		}

		controller = controllers.NewMediaController(b.conn, b.events, cast.DefaultSender, transportId)
		if err := controller.Start(ctx); err != nil {
			b.UpdateStatus(boggart.BindStatusOffline)
			return nil, err
		}

		b.mutex.Lock()
		b.media = controller
		b.mutex.Unlock()
	}

	return controller, nil
}

func (b *Bind) launchApp(ctx context.Context, appID string) (string, error) {
	result, err := b.receiver.GetAppAvailability(ctx, appID)
	if err != nil {
		b.UpdateStatus(boggart.BindStatusOffline)
		return "", err
	}

	if !result {
		return "", errors.New("unsupported app with ID " + appID)
	}

	status, err := b.receiver.GetStatus(ctx)
	if err != nil {
		b.UpdateStatus(boggart.BindStatusOffline)
		return "", err
	}

	app := status.GetSessionByAppId(appID)
	if app == nil {
		status, err = b.receiver.LaunchApp(ctx, appID)
		if err != nil {
			b.UpdateStatus(boggart.BindStatusOffline)
			return "", err
		}

		app = status.GetSessionByAppId(appID)
	}

	if app == nil {
		return "", errors.New("failed to get transport")
	}

	return *app.TransportId, nil
}
