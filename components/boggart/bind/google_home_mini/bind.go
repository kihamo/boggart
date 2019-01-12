package google_home_mini

import (
	"sync"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/providers/google/home"
	"github.com/kihamo/boggart/components/boggart/providers/google/home/client"
	"github.com/kihamo/boggart/components/voice/players/chromecast"
)

type Bind struct {
	boggart.DeviceBindBase
	boggart.DeviceBindMQTT

	mutex    sync.RWMutex
	initOnce sync.Once

	clientGoogleHome *client.GoogleHome
	clientChromecast *chromecast.Player
	host             string
}

func (b *Bind) ClientGoogleHome() *client.GoogleHome {
	b.mutex.RLock()
	c := b.clientGoogleHome
	b.mutex.RUnlock()

	if c != nil {
		return c
	}

	ctrl := home.NewClient(b.host)

	b.mutex.Lock()
	b.clientGoogleHome = ctrl
	b.mutex.Unlock()

	return ctrl
}

func (b *Bind) ClientChromecast() *chromecast.Player {
	b.mutex.RLock()
	c := b.clientChromecast
	b.mutex.RUnlock()

	if c != nil {
		return c
	}

	ctrl := chromecast.New(b.host, chromecast.DefaultPort)

	b.mutex.Lock()
	b.clientChromecast = ctrl
	b.mutex.Unlock()

	return ctrl
}

func (b *Bind) UpdateStatus(status boggart.DeviceStatus) {
	if status == boggart.DeviceStatusOffline && status != b.Status() {
		b.mutex.Lock()
		b.clientGoogleHome = nil

		if b.clientChromecast != nil {
			go b.clientChromecast.Close()
		}

		b.clientChromecast = nil
		b.mutex.Unlock()
	}

	b.DeviceBindBase.UpdateStatus(status)
}
