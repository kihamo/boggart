package scale

import (
	"context"
	"time"

	"github.com/kihamo/boggart/atomic"
	"github.com/kihamo/boggart/components/boggart/di"
	"github.com/kihamo/boggart/providers/xiaomi/scale"
)

type Bind struct {
	di.MQTTBind
	di.WorkersBind

	config               *Config
	provider             *scale.Client
	currentProfile       atomic.Value
	measureStartDatetime atomic.Time
}

func (b *Bind) Run() error {
	b.notifyCurrentProfile(context.Background())
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

func (b *Bind) SetProfile(name string) *Profile {
	var profile *Profile

	for _, p := range b.config.Profiles {
		if p.Name == name {
			profile = p
			break
		}
	}

	if profile == nil {
		profile = profileGuest
	}

	b.currentProfile.Store(profile)
	b.measureStartDatetime.Set(time.Now())

	return nil
}

func (b *Bind) Close() error {
	return b.provider.Close()
}
