package chromecast

import (
	"context"
	"errors"
	"fmt"
	"math"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/barnybug/go-cast"
	"github.com/barnybug/go-cast/controllers"
	"github.com/barnybug/go-cast/events"
	castnet "github.com/barnybug/go-cast/net"
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/atomic"
	"github.com/kihamo/boggart/components/mqtt"
	"go.uber.org/multierr"
)

const (
	PlayerStateIdle      = "IDLE"
	PlayerStatePlaying   = "PLAYING"
	PlayerStateBuffering = "BUFFERING"
	PlayerStatePaused    = "PAUSED"

	IdleReasonCancelled   = "CANCELLED"
	IdleReasonInterrupted = "INTERRUPTED"
	IdleReasonFinished    = "FINISHED"
	IdleReasonError       = "ERROR"
)

type Bind struct {
	boggart.BindBase
	boggart.BindMQTT

	host net.IP
	port int

	volume         *atomic.Uint32Null
	mute           *atomic.BoolNull
	status         *atomic.String
	mediaContentID *atomic.String

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

func (b *Bind) Run() error {
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
		b.Logger().Debug("Connect failed", "error", err.Error())
		return err
	}

	if err := b.heartbeat.Start(ctx); err != nil {
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

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	// close running aps
	if _, e := b.receiver.QuitApp(ctx); e != nil {
		err = multierr.Append(err, e)
	}

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
		if e := b.media.Close(); e != nil {
			err = multierr.Append(err, e)
		}

		b.media = nil
	}
	b.mutex.Unlock()

	return err
}

func (b *Bind) doEvents() {
	ctx := context.Background()

	for {
		select {
		case event := <-b.events:
			switch t := event.(type) {
			case events.Connected:
				b.Logger().Debug("Event Connected")

				b.UpdateStatus(boggart.BindStatusOnline)

			case events.Disconnected: // from ReceiverController
				b.Logger().Debug("Event Disconnected")

				if err := b.Close(); err != nil {
					b.Logger().Error("Close failed", "error", err.Error())
				}

			case events.AppStarted: // from ReceiverController
				b.Logger().Debug("Event AppStarted")

			case events.AppStopped: // from ReceiverController
				b.Logger().Debug("Event AppStopped")

				if t.AppID == cast.AppMedia {
					b.mutex.Lock()

					if b.media != nil {
						if err := b.media.Close(); err != nil {
							b.Logger().Error("Media close failed", "error", err.Error())
						}

						b.media = nil
					}

					b.mutex.Unlock()
				}

			case events.StatusUpdated: // from ReceiverController
				b.Logger().Debug("Event StatusUpdated",
					"level", t.Level,
					"muted", t.Muted,
				)

				sn := mqtt.NameReplace(b.SerialNumber())

				volume := uint32(math.Round(t.Level * 100))
				if ok := b.volume.Set(volume); ok {
					_ = b.MQTTPublishAsync(ctx, MQTTPublishTopicStateVolume.Format(sn), volume)
				}

				if ok := b.mute.Set(t.Muted); ok {
					_ = b.MQTTPublishAsync(ctx, MQTTPublishTopicStateMute.Format(sn), t.Muted)
				}

			case controllers.MediaStatus:
				b.Logger().Debug("Event MediaStatus",
					"state", t.PlayerState,
					"reason", t.IdleReason,
				)

				sn := mqtt.NameReplace(b.SerialNumber())

				if ok := b.status.Set(t.PlayerState); ok {
					_ = b.MQTTPublishAsync(ctx, MQTTPublishTopicStateStatus.Format(sn), strings.ToLower(t.PlayerState))
				}

				if t.PlayerState == PlayerStateIdle && t.IdleReason == IdleReasonFinished {
					if _, err := b.receiver.QuitApp(ctx); err != nil {
						b.Logger().Error("Quit app failed", "error", err.Error())
					}
				}

				if t.Media != nil {
					if ok := b.mediaContentID.Set(t.Media.ContentId); ok {
						_ = b.MQTTPublishAsync(ctx, MQTTPublishTopicStateContent.Format(sn), t.Media.ContentId)
					}
				}

			default:
				b.Logger().Error("Unknown event", "event", fmt.Sprintf("%#v", t))
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
