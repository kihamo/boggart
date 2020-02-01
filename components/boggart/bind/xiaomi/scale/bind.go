package scale

import (
	"context"
	"time"

	"github.com/kihamo/boggart/atomic"
	"github.com/kihamo/boggart/components/boggart/di"
	"github.com/kihamo/boggart/providers/xiaomi/scale"
)

type Bind struct {
	di.MetaBind
	di.MQTTBind
	di.LoggerBind
	di.WorkersBind

	config               *Config
	provider             *scale.Client
	currentProfile       atomic.Value
	measureStartDatetime *atomic.Time
}

func (b *Bind) Run() error {
	b.notifyCurrentProfile(context.Background())
	b.Meta().SetMAC(b.config.MAC.HardwareAddr)
	return nil
}

func (b *Bind) CurrentProfile() *Profile {
	if profile := b.currentProfile.Load(); profile != nil {
		if p, ok := profile.(*Profile); ok {
			return p
		}
	}

	return profileGuest
}

func (b *Bind) SetProfile(name string) {
	profile := b.Profile(name)

	if profile == nil {
		profile = profileGuest
	}

	b.currentProfile.Store(profile)
	b.measureStartDatetime.Set(time.Now())
}

func (b *Bind) Profile(name string) *Profile {
	if name == "" {
		return nil
	}

	for _, p := range b.Profiles() {
		if p.Name == name {
			return p
		}
	}

	return nil
}

func (b *Bind) Profiles() []*Profile {
	profiles := make([]*Profile, 0, len(b.config.Profiles)+1)
	guestExist := false

	for _, profile := range b.config.Profiles {
		profiles = append(profiles, profile)

		if profile.Name == profileGuest.Name {
			guestExist = true
		}
	}

	if !guestExist {
		profiles = append(profiles, profileGuest)
	}

	return profiles
}

func (b *Bind) Close() error {
	return b.provider.Close()
}
