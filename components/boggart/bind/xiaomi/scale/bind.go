package scale

import (
	"context"
	"time"

	"github.com/kihamo/boggart/atomic"
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/providers/xiaomi/scale"
)

type Bind struct {
	boggart.BindBase
	boggart.BindMQTT
	config *Config

	provider           *scale.Client
	currentProfile     atomic.Value
	setProfileDatetime atomic.Time
}

func (b *Bind) Run() error {
	b.notifyCurrentProfile(context.Background())
	return nil
}

func (b *Bind) CurrentProfile() *Profile {
	profile := b.currentProfile.Load()
	if profile != nil {
		if p, ok := profile.(*Profile); ok {
			return p
		}
	}

	return nil
}

func (b *Bind) SetProfile(name string) *Profile {
	for _, profile := range b.config.Profiles {
		if profile.Name == name {
			b.currentProfile.Store(profile)
			b.setProfileDatetime.Set(time.Now())

			return profile
		}
	}

	return nil
}

func (b *Bind) Close() error {
	return b.provider.Close()
}
