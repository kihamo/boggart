package google_home

import (
	"strconv"
	"sync"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/providers/google/home"
	"github.com/kihamo/boggart/providers/google/home/client"
)

type Bind struct {
	boggart.BindBase
	config *Config
	mutex  sync.RWMutex
	client *client.GoogleHome
}

func (b *Bind) ClientGoogleHome() *client.GoogleHome {
	b.mutex.RLock()
	c := b.client
	b.mutex.RUnlock()

	if c != nil {
		return c
	}

	ctrl := home.NewClient(b.config.Host.String() + ":" + strconv.Itoa(b.config.Port))

	b.mutex.Lock()
	b.client = ctrl
	b.mutex.Unlock()

	return ctrl
}

func (b *Bind) UpdateStatus(status boggart.BindStatus) {
	if status == boggart.BindStatusOffline && !b.IsStatus(status) {
		b.mutex.Lock()
		b.client = nil
		b.mutex.Unlock()
	}

	b.BindBase.UpdateStatus(status)
}
