package google_home

import (
	"net"
	"strconv"
	"sync"
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/providers/google/home"
	"github.com/kihamo/boggart/components/boggart/providers/google/home/client"
)

type Bind struct {
	boggart.BindBase

	mutex  sync.RWMutex
	client *client.GoogleHome

	host net.IP
	port int

	livenessInterval time.Duration
	livenessTimeout  time.Duration
}

func (b *Bind) ClientGoogleHome() *client.GoogleHome {
	b.mutex.RLock()
	c := b.client
	b.mutex.RUnlock()

	if c != nil {
		return c
	}

	ctrl := home.NewClient(b.host.String() + ":" + strconv.Itoa(b.port))

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
