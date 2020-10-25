package chromecast

import (
	"context"
	"errors"
	"fmt"
	"math"
	"strings"
	"sync"

	"github.com/barnybug/go-cast"
	"github.com/barnybug/go-cast/controllers"
	"github.com/barnybug/go-cast/events"
	castnet "github.com/barnybug/go-cast/net"
	"github.com/kihamo/boggart/atomic"
	"github.com/kihamo/boggart/components/boggart/di"
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

type eventClose struct{}

type Bind struct {
	di.LoggerBind
	di.MQTTBind
	di.ProbesBind
	di.WidgetBind

	config *Config

	disconnected   *atomic.BoolNull
	volume         *atomic.Uint32Null
	mute           *atomic.BoolNull
	status         *atomic.String
	mediaContentID *atomic.String

	events chan events.Event
	mutex  sync.RWMutex

	conn           *castnet.Connection
	ctrlConnection *controllers.ConnectionController
	ctrlHeartbeat  *controllers.HeartbeatController
	ctrlReceiver   *controllers.ReceiverController
	ctrlMedia      *controllers.MediaController
}

func (b *Bind) Run() error {
	b.disconnected.Nil()

	return nil
}

func (b *Bind) initConnect() error {
	ctx := context.Background()

	conn := castnet.NewConnection()
	if err := conn.Connect(ctx, b.config.Host.IP, b.config.Port); err != nil {
		return err
	}

	ctrlConnection := controllers.NewConnectionController(conn, b.events, cast.DefaultSender, cast.DefaultReceiver)
	if err := ctrlConnection.Start(ctx); err != nil {
		return err
	}

	ctrlHeartbeat := controllers.NewHeartbeatController(conn, b.events, cast.TransportSender, cast.TransportReceiver)
	if err := ctrlHeartbeat.Start(ctx); err != nil {
		return err
	}

	ctrlReceiver := controllers.NewReceiverController(conn, b.events, cast.DefaultSender, cast.DefaultReceiver)
	if err := ctrlReceiver.Start(ctx); err != nil {
		return err
	}

	go b.doEvents()

	b.mutex.Lock()
	b.conn = conn
	b.ctrlConnection = ctrlConnection
	b.ctrlHeartbeat = ctrlHeartbeat
	b.ctrlReceiver = ctrlReceiver
	b.mutex.Unlock()

	b.disconnected.False()

	return nil
}

func (b *Bind) Close() (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), b.config.LivenessTimeout)
	defer cancel()
	defer func() {
		b.events <- eventClose{}
	}()

	b.mutex.RLock()
	conn := b.conn
	ctrlConnection := b.ctrlConnection
	ctrlHeartbeat := b.ctrlHeartbeat
	ctrlReceiver := b.ctrlReceiver
	ctrlMedia := b.ctrlMedia
	b.mutex.RUnlock()

	if conn == nil {
		return nil
	}

	if _, e := ctrlReceiver.QuitApp(ctx); e != nil {
		err = multierr.Append(err, e)
	}

	ctrlHeartbeat.Stop()

	if e := ctrlConnection.Close(); e != nil {
		err = multierr.Append(err, e)
	}

	if e := conn.Close(); e != nil {
		err = multierr.Append(err, e)
	}

	if ctrlMedia != nil {
		if e := ctrlMedia.Close(); e != nil {
			err = multierr.Append(err, e)
		}
	}

	return err
}

func (b *Bind) doEvents() {
	ctx := context.Background()

	for {
		event := <-b.events
		switch t := event.(type) {
		case events.Disconnected: // from HeartbeatController
			b.Logger().Debug("Event Disconnected", "reason", t.Reason.Error())

			b.disconnected.True()

		case events.AppStarted: // from ReceiverController
			b.Logger().Debug("Event AppStarted")

		case events.AppStopped: // from ReceiverController
			b.Logger().Debug("Event AppStopped",
				"app-id", t.AppID,
				"display-name", t.DisplayName,
				"status-text", t.StatusText,
			)

			if t.AppID == cast.AppMedia {
				b.mutex.RLock()
				ctrlMedia := b.ctrlMedia
				b.mutex.RUnlock()

				if ctrlMedia != nil {
					if err := ctrlMedia.Close(); err != nil {
						b.Logger().Error("Media close failed", "error", err.Error())
					}

					b.mutex.Lock()
					b.ctrlMedia = nil
					b.mutex.Unlock()
				}
			}

		case events.StatusUpdated: // from ReceiverController
			b.Logger().Debug("Event StatusUpdated",
				"level", t.Level,
				"muted", t.Muted,
			)

			volume := uint32(math.Round(t.Level * 100))
			_ = b.MQTT().PublishAsync(ctx, b.config.TopicStateVolume, volume)

			_ = b.MQTT().PublishAsync(ctx, b.config.TopicStateMute, t.Muted)

		case controllers.MediaStatus:
			b.Logger().Debug("Event MediaStatus",
				"state", t.PlayerState,
				"reason", t.IdleReason,
			)

			_ = b.MQTT().PublishAsync(ctx, b.config.TopicStateStatus, strings.ToLower(t.PlayerState))

			if t.PlayerState == PlayerStateIdle && t.IdleReason == IdleReasonFinished {
				b.mutex.RLock()
				ctrlReceiver := b.ctrlReceiver
				b.mutex.RUnlock()

				if _, err := ctrlReceiver.QuitApp(ctx); err != nil {
					b.Logger().Error("Quit app failed", "error", err.Error())
				}
			}

			if t.Media != nil {
				_ = b.MQTT().PublishAsync(ctx, b.config.TopicStateContent, t.Media.ContentId)
			}

		case eventClose:
			return

		default:
			b.Logger().Error("Unknown event", "event", fmt.Sprintf("%#v", t))
		}
	}
}

func (b *Bind) Media(ctx context.Context) (*controllers.MediaController, error) {
	b.mutex.RLock()
	ctrlMedia := b.ctrlMedia
	conn := b.conn
	b.mutex.RUnlock()

	if ctrlMedia == nil {
		transportID, err := b.launchApp(ctx, cast.AppMedia)
		if err != nil {
			return nil, err
		}

		ctrlMedia = controllers.NewMediaController(conn, b.events, cast.DefaultSender, transportID)
		if err := ctrlMedia.Start(ctx); err != nil {
			return nil, err
		}

		b.mutex.Lock()
		b.ctrlMedia = ctrlMedia
		b.mutex.Unlock()
	}

	return ctrlMedia, nil
}

func (b *Bind) launchApp(ctx context.Context, appID string) (string, error) {
	b.mutex.RLock()
	ctrlReceiver := b.ctrlReceiver
	b.mutex.RUnlock()

	if ctrlReceiver == nil {
		return "", errors.New("receiver controller isn't init")
	}

	result, err := ctrlReceiver.GetAppAvailability(ctx, appID)
	if err != nil {
		return "", err
	}

	if !result {
		return "", errors.New("unsupported app with ID " + appID)
	}

	status, err := ctrlReceiver.GetStatus(ctx)
	if err != nil {
		return "", err
	}

	app := status.GetSessionByAppId(appID)
	if app == nil {
		status, err = ctrlReceiver.LaunchApp(ctx, appID)
		if err != nil {
			return "", err
		}

		app = status.GetSessionByAppId(appID)
	}

	if app == nil {
		return "", errors.New("failed to get transport")
	}

	return *app.TransportId, nil
}
