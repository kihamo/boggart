package google_home

import (
	"sync"
	"time"

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
	clientChromeCast *chromecast.Player
	host             string

	livenessInterval time.Duration
	livenessTimeout  time.Duration
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

func (b *Bind) ClientChromeCast() *chromecast.Player {
	b.mutex.RLock()
	c := b.clientChromeCast
	b.mutex.RUnlock()

	if c != nil {
		return c
	}

	ctrl := chromecast.New(b.host, chromecast.DefaultPort)

	b.mutex.Lock()
	b.clientChromeCast = ctrl
	b.mutex.Unlock()

	return ctrl
}

func (b *Bind) UpdateStatus(status boggart.DeviceStatus) {
	if status == boggart.DeviceStatusOffline && status != b.Status() {
		b.mutex.Lock()
		b.clientGoogleHome = nil

		if b.clientChromeCast != nil {
			go b.clientChromeCast.Close()
		}

		b.clientChromeCast = nil
		b.mutex.Unlock()
	}

	b.DeviceBindBase.UpdateStatus(status)
}
