package webos

import (
	"context"
	"errors"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/kihamo/boggart/atomic"
	"github.com/kihamo/boggart/components/boggart/di"
	"github.com/snabb/webostv"
)

var defaultDialerLGWebOS = webostv.Dialer{
	DisableTLS: true,
	WebsocketDialer: &websocket.Dialer{
		Proxy: http.ProxyFromEnvironment,
		NetDial: (&net.Dialer{
			Timeout:   time.Second * 5,
			KeepAlive: time.Second * 30, // ensure we notice if the TV goes away
		}).Dial,
	},
}

type Bind struct {
	di.ConfigBind
	di.LoggerBind
	di.MetaBind
	di.MQTTBind
	di.ProbesBind
	di.WidgetBind
	di.WorkersBind

	mutex  sync.RWMutex
	client *webostv.Tv

	power        *atomic.BoolNull
	quitMonitors chan struct{}
}

func (b *Bind) config() *Config {
	return b.Config().Bind().(*Config)
}

func (b *Bind) Run() error {
	b.power.Nil()

	b.mutex.Lock()
	defer b.mutex.Unlock()

	b.client = nil
	b.quitMonitors = make(chan struct{})

	if cfg := b.config(); cfg.MAC != nil {
		b.Meta().SetMAC(cfg.MAC.HardwareAddr)

		return b.MQTT().PublishAsync(context.Background(), cfg.TopicStatePower.Format(b.Meta().MACAsString()), false)
	}

	return nil
}

func (b *Bind) initClient() error {
	cfg := b.config()

	client, err := defaultDialerLGWebOS.Dial(cfg.Host)
	if err != nil {
		return err
	}

	client.SetDebug(func(s string) {
		b.Logger().Debug(s)
	})

	go func() {
		if err := client.MessageHandler(); err != nil {
			b.mutex.Lock()
			b.client = nil
			b.mutex.Unlock()
		}
	}()

	newKey, err := client.Register(cfg.Key)
	if err != nil {
		return err
	}

	if cfg.Key != newKey {
		b.Logger().Warnf("Key changed before %s after %s", cfg.Key, newKey)
	}

	b.mutex.Lock()
	b.client = client
	b.mutex.Unlock()

	return nil
}

func (b *Bind) Client() *webostv.Tv {
	b.mutex.RLock()
	defer b.mutex.RUnlock()

	return b.client
}

func (b *Bind) Toast(message string) error {
	client := b.Client()
	if client == nil {
		return errors.New("client isn't init")
	}

	_, err := client.SystemNotificationsCreateToast(message)

	return err
}

func (b *Bind) Close() error {
	close(b.quitMonitors)

	if client := b.Client(); client != nil {
		return client.Close()
	}

	return nil
}
